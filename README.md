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
real 10.97
user 121.80
sys 11.21
```

The top from the profiler:

```
File: main
Type: cpu
Time: Sep 29, 2024 at 8:15pm (EEST)
Duration: 10.78s, Total samples = 126.79s (1176.31%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 111.67s, 88.07% of 126.79s total
Dropped 37 nodes (cum <= 0.63s)
Showing top 10 nodes out of 28
      flat  flat%   sum%        cum   cum%
    27.73s 21.87% 21.87%     39.88s 31.45%  main.parseLine
    15.82s 12.48% 34.35%     34.62s 27.30%  runtime.mapaccess2_faststr
    11.80s  9.31% 43.65%     11.89s  9.38%  main.atoi
    11.27s  8.89% 52.54%     11.27s  8.89%  internal/runtime/syscall.Syscall6
    11.22s  8.85% 61.39%     45.88s 36.19%  main.(*actor).processChunk-range1
    10.29s  8.12% 69.51%     10.29s  8.12%  aeshashbody
     8.63s  6.81% 76.32%      8.63s  6.81%  indexbytebody
     5.72s  4.51% 80.83%     17.50s 13.80%  main.scanLines
     4.75s  3.75% 84.57%    126.71s 99.94%  main.(*actor).processChunk.(*MeasurementsReader).All.func1
     4.44s  3.50% 88.07%     34.63s 27.31%  bufio.(*Scanner).Scan
```
