# 1brc

This is an attempt to solve a 1 Billion Rows Challenge in Go.


## My Rules

- Allowed to use google, but not searching for the solution
- Not allowed ChatGPT and other similar tools
- Copilot is enabled but I only use it as autocompletion sometimes anyways
- Basically, no cheating

## My setup

the code is tested on a Lenovo laptop running Ubuntu.

CPU: Intel(R) Core(TM) Ultra 7 155H


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

The baseline solution coded in JAVA in my machine takes around `2mins`

(2)

copy the `measurements.txt` into the projects folder (from previous step)

Run my solution

```
./calculate_average.sh measurements.txt
```

My timings:

```
real 2.98
user 57.13
sys 1.89
```

The top from the profiler:

```
File: main
Type: cpu
Time: Oct 29, 2024 at 7:28pm (EET)
Duration: 3.20s, Total samples = 61.97s (1935.33%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 59.32s, 95.72% of 61.97s total
Dropped 17 nodes (cum <= 0.31s)
Showing top 10 nodes out of 26
      flat  flat%   sum%        cum   cum%
    14.09s 22.74% 22.74%     25.98s 41.92%  main.(*actor).processChunk-range1
    13.47s 21.74% 44.47%     22.93s 37.00%  main.parseLine
     9.41s 15.18% 59.66%      9.42s 15.20%  main.atoi
     8.88s 14.33% 73.99%      8.89s 14.35%  main.djb2Hash (inline)
     3.56s  5.74% 79.73%     10.97s 17.70%  bufio.(*Scanner).Scan
     2.56s  4.13% 83.86%      2.56s  4.13%  memeqbody
     2.29s  3.70% 87.56%      2.29s  3.70%  indexbytebody
     1.87s  3.02% 90.58%      1.87s  3.02%  internal/runtime/syscall.Syscall6
     1.68s  2.71% 93.29%      4.71s  7.60%  main.scanLines
     1.51s  2.44% 95.72%     61.94s   100%  main.(*actor).processChunk.(*MeasurementsReader).All.func2
```
