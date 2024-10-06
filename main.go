package main

import (
	"bufio"
	"bytes"
	"cmp"
	"io"
	"iter"
	"maps"
	"os"
	"runtime"
	"runtime/pprof"
	"slices"
	"strconv"
	"sync"
)

func main() {
	if len(os.Args) > 2 {
		f, err := os.Create(os.Args[2])
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if err := pprof.StartCPUProfile(f); err != nil {
			panic(err)
		}
		defer pprof.StopCPUProfile()
	}

	if err := run(os.Args[1]); err != nil {
		panic(err)
	}
}

func run(fname string) error {
	fd, err := os.Open(fname)
	if err != nil {
		return err
	}

	defer fd.Close()

	w := bufio.NewWriter(os.Stdout)

	defer w.Flush()

	fileInfo, err := fd.Stat()
	if err != nil {
		return err
	}

	fileSize := fileInfo.Size()

	numberOfCPUs := runtime.NumCPU()
	chunks := make(chan chunk, numberOfCPUs)

	results := make(chan subResult, numberOfCPUs)

	wg := sync.WaitGroup{}
	wg.Add(numberOfCPUs)

	for i := 0; i < numberOfCPUs; i++ {
		actor := newActor(fname, i+1, chunks, results)
		go func() {
			defer wg.Done()
			actor.run()
		}()
	}

	done := make(chan struct{})

	merger := merger{inbox: results}
	go func() {
		defer close(done)

		merger.run(w)
	}()

	chunkSize := max(1, fileSize) / int64(numberOfCPUs)

	start := int64(0)
	end := chunkSize

	cnt := 0

	for start < end {
		chunk := chunk{}

		b := make([]byte, 1)

		for end < fileSize && b[0] != '\n' {
			_, err := fd.Seek(end, 0)
			if err != nil {
				panic(err)
			}

			_, err = fd.Read(b)
			if err != nil {
				panic(err)
			}

			end++
		}

		chunk.start = start
		chunk.end = end

		start = end
		end = min(fileSize, start+chunkSize)

		cnt++

		chunks <- chunk
	}

	close(chunks)

	wg.Wait()

	close(results)

	<-done

	return nil
}

type actor struct {
	id     int
	inbox  <-chan chunk
	outbox chan<- subResult
	fd     *os.File
}

func newActor(fname string, id int, inbox <-chan chunk, outbox chan<- subResult) *actor {
	fd, err := os.Open(fname)
	if err != nil {
		panic(err)
	}

	return &actor{
		id:     id,
		inbox:  inbox,
		outbox: outbox,
		fd:     fd,
	}
}

func (a *actor) run() {
	for chunk := range a.inbox {
		result := a.processChunk(chunk)

		a.outbox <- result
	}
}

func (a *actor) processChunk(chunk chunk) subResult {
	const size uint32 = 100000

	resultSlice := make([]result, size)
	positions := make([]int, 0, size/2)

	reader := MeasurementsReader{chunk: chunk, fd: a.fd}

	for measurement := range reader.All() {
		hv := djb2Hash(measurement.station, size)

		if resultSlice[hv].station == "" {
			resultSlice[hv].station = string(measurement.station)
			resultSlice[hv].min = measurement.temp
			resultSlice[hv].max = measurement.temp
			resultSlice[hv].total = 1
			resultSlice[hv].sum = measurement.temp
			positions = append(positions, int(hv))
		} else if resultSlice[hv].station == string(measurement.station) {
			resultSlice[hv].min = min(resultSlice[hv].min, measurement.temp)
			resultSlice[hv].max = max(resultSlice[hv].max, measurement.temp)
			resultSlice[hv].total++
			resultSlice[hv].sum += measurement.temp
		} else {
			for {
				hv = (hv + 1) % size

				if resultSlice[hv].station == "" {
					resultSlice[hv].station = string(measurement.station)
					resultSlice[hv].min = measurement.temp
					resultSlice[hv].max = measurement.temp
					resultSlice[hv].total = 1
					resultSlice[hv].sum = measurement.temp
					positions = append(positions, int(hv))

					break
				} else if resultSlice[hv].station == string(measurement.station) {
					resultSlice[hv].min = min(resultSlice[hv].min, measurement.temp)
					resultSlice[hv].max = max(resultSlice[hv].max, measurement.temp)
					resultSlice[hv].total++
					resultSlice[hv].sum += measurement.temp

					break
				}
			}
		}
	}

	it := func(yield func(*result) bool) {
		for _, pos := range positions {
			if !yield(&resultSlice[pos]) {
				break
			}
		}
	}

	ans := subResult{
		items: it,
	}

	return ans
}

// http://www.cse.yorku.ca/~oz/hash.html (almost)
func djb2Hash(b []byte, sliceSize uint32) uint32 {
	var hash uint32 = 5381

	for _, c := range b {
		hash = ((hash << 5) + hash) ^ uint32(c)
	}

	return uint32(hash % sliceSize)
}

type MeasurementsReader struct {
	chunk chunk
	fd    *os.File
}

func (m *MeasurementsReader) All() iter.Seq[measurement] {
	return func(yield func(measurement) bool) {

		sectionReader := io.NewSectionReader(m.fd, m.chunk.start, m.chunk.end-m.chunk.start)

		scanner := bufio.NewScanner(sectionReader)

		const maxbuf = 1024 * 256

		buffer := make([]byte, maxbuf)

		scanner.Buffer(buffer, maxbuf)
		scanner.Split(scanLines)

		var item measurement

		for scanner.Scan() {
			parseLine(scanner.Bytes(), &item)

			yield(item)
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}
	}
}

type chunk struct {
	start int64
	end   int64
}

type result struct {
	station string
	min     int64
	max     int64
	total   int64
	sum     int64
}

type subResult struct {
	items iter.Seq[*result]
}

type measurement struct {
	line    []byte
	station []byte
	temp    int64
}

func scanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		return i + 1, data[0:i], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}

func parseLine(b []byte, ans *measurement) {
	for i := 0; i < len(b); i++ {
		if b[i] == ';' {
			ans.station = b[:i]
			ans.temp = atoi(b[i+1:])

			return
		}
	}

	return
}

func atoi(b []byte) int64 {
	sign := int64(1)
	start := 0

	if b[0] == '-' {
		sign = -1
		start = 1
	}

	result := int64(b[start] - '0')

	if b[start+1] == '.' {
		return sign * (result*10 + int64(b[start+2]-'0'))
	}

	return sign * (result*100 + int64(b[start+1]-'0')*10 + int64(b[start+3]-'0'))
}

type merger struct {
	inbox <-chan subResult
}

func (m *merger) run(w io.Writer) {
	merged := m.mergeAndSortResults()

	// results look like this:
	// {station1=min/max/mean, station2=min/max/mean, ...}
	// so we need 2 + 30 (stationName) +1(=)+4(min)+1(/)+4(max)+1(/)+4(mean) = 46 bytes per station
	// and 2 bytes for the brackets and 2 for the space and comma
	estimatedSize := 2 + len(merged)*46 + 4

	sb := bytes.NewBuffer(make([]byte, 0, estimatedSize))

	sb.WriteString("{")

	for i := range merged {
		m.writerResult(sb, merged[i])

		if i < len(merged)-1 {
			sb.WriteString(", ")
		}
	}

	sb.WriteString("} ")

	w.Write(sb.Bytes())
}

func (m *merger) mergeAndSortResults() []*result {
	stations := make(map[string]*result, 10000)

	sortFunc := func(a, b *result) int {
		return cmp.Compare(a.station, b.station)
	}

	for sub := range m.inbox {
		for result := range sub.items {
			if existing, ok := stations[result.station]; ok {
				existing.min = min(existing.min, result.min)
				existing.max = max(existing.max, result.max)
				existing.total += result.total
				existing.sum += result.sum
			} else {
				stations[result.station] = result
			}
		}
	}

	merged := slices.SortedFunc(maps.Values(stations), sortFunc)

	return merged
}

func (m *merger) writerResult(sb *bytes.Buffer, res *result) {
	sb.WriteString(res.station)
	sb.WriteByte('=')
	writeTemperature(sb, res.min)
	sb.WriteByte('/')

	mean := roundToDecimal(res.sum * 10 / int64(res.total))

	writeTemperature(sb, mean)

	sb.WriteByte('/')

	writeTemperature(sb, res.max)
}

func writeTemperature(sb *bytes.Buffer, temp int64) {
	if temp < 0 {
		sb.WriteByte('-')

		temp = -temp
	}

	sb.WriteString(strconv.FormatInt(temp/10, 10))

	sb.WriteByte('.')
	sb.WriteByte('0' + byte(temp%10))
}

func roundToDecimal(temp int64) int64 {
	if temp >= 0 {
		return (temp + 5) / 10
	}

	return (temp - 5) / 10
}
