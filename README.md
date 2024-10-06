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
real 9.67
user 96.25
sys 12.09
```

The top from the profiler:

```
File: main
Type: cpu
Time: Oct 6, 2024 at 7:24pm (EEST)
Duration: 9.87s, Total samples = 108.67s (1100.63%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 102.30s, 94.14% of 108.67s total
Dropped 31 nodes (cum <= 0.54s)
Showing top 10 nodes out of 26
      flat  flat%   sum%        cum   cum%
    24.57s 22.61% 22.61%     37.01s 34.06%  main.parseLine
    17.07s 15.71% 38.32%     31.64s 29.12%  main.(*actor).processChunk-range1
    12.48s 11.48% 49.80%     12.48s 11.48%  internal/runtime/syscall.Syscall6
    12.06s 11.10% 60.90%     12.17s 11.20%  main.atoi
    10.74s  9.88% 70.78%     10.86s  9.99%  main.djb2Hash (inline)
     8.49s  7.81% 78.60%      8.49s  7.81%  indexbytebody
     5.31s  4.89% 83.48%     16.79s 15.45%  main.scanLines
     4.45s  4.09% 87.58%    108.65s   100%  main.(*actor).processChunk.(*MeasurementsReader).All.func2
     3.80s  3.50% 91.07%     34.13s 31.41%  bufio.(*Scanner).Scan
     3.33s  3.06% 94.14%      3.33s  3.06%  memeqbody
(pprof)
```
