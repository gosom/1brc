# 1brc

This is an attempt to solve a 1 Billion Rows Challenge in Go.


## My Rules

- Allowed to use google, but not searching for the solution
- Not allowed ChatGPT and other similar tools
- Copilot is enabled but I only use it as autocompletion sometimes anyways
- Basically, no cheating

## My setup

the code is tested on a Lenovo laptop running Fedora.

CPU: AMD Ryzen 7 PRO 4750U with Radeon Graphics

while testing the laptop was plugged in to power.


## Instructions

Create the dataset from the official challenge:

(1)
Checkout the official challenge

```
git@github.com:gunnarmorling/1brc.git
cd 1brc
./mvnw clean verify
./create_measurements.sh 1000000000 # this creates a measurements.txt file
/usr/bin/time -o baseline.txt ./calculate_average_baseline.sh >baseline.results
```

The baseline solution coded in JAVA in my machine takes:

```
real 271.27 ~= 4.52 minutes
user 228.36
sys 19.98
```

(2)

copy the `measurements.txt` into the projects folder (from previous step)

Run my solution

```
./calculate_average.sh measurements.txt
```

My timings:

```
real 9.27
user 97.72
sys 10.76
```

The top from the profiler:

```
File: main
Type: cpu
Time: Oct 6, 2024 at 7:54pm (EEST)
Duration: 9.40s, Total samples = 107.86s (1147.96%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 100.13s, 92.83% of 107.86s total
Dropped 36 nodes (cum <= 0.54s)
Showing top 10 nodes out of 27
      flat  flat%   sum%        cum   cum%
    22.26s 20.64% 20.64%     34.34s 31.84%  main.parseLine
    16.97s 15.73% 36.37%     33.22s 30.80%  main.(*actor).processChunk-range1
    11.86s 11.00% 47.37%     11.91s 11.04%  main.atoi
    11.30s 10.48% 57.84%     11.36s 10.53%  main.djb2Hash (inline)
    10.03s  9.30% 67.14%     10.03s  9.30%  internal/runtime/syscall.Syscall6
     9.22s  8.55% 75.69%      9.22s  8.55%  indexbytebody
     5.96s  5.53% 81.22%     18.51s 17.16%  main.scanLines
     4.88s  4.52% 85.74%    107.79s 99.94%  main.(*actor).processChunk.(*MeasurementsReader).All.func2
     4.04s  3.75% 89.49%     33.82s 31.36%  bufio.(*Scanner).Scan
     3.61s  3.35% 92.83%      3.61s  3.35%  memeqbody
```
