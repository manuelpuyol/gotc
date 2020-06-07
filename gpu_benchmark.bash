#!/bin/bash

run_hashrate() {
  ./gotc -g -b > benchmarks/hashrate/gpu-1.txt
  ./gotc -g -b > benchmarks/hashrate/gpu-2.txt
  ./gotc -g -b > benchmarks/hashrate/gpu-3.txt
  ./gotc -g -b > benchmarks/hashrate/gpu-4.txt
  ./gotc -g -b > benchmarks/hashrate/gpu-5.txt
}

run_timings() {
  echo "difficulty 4"
  (time ./gotc -s -f data/benchmark4.txt -d 4 -g) > out.txt 2> benchmarks/timings/d4/parallel_gpu-$n.txt

  echo "difficulty 5"
  (time ./gotc -s -f data/benchmark3.txt -d 5 -g) > out.txt 2> benchmarks/timings/d5/parallel_gpu-$n.txt

  echo "difficulty 6"
  (time ./gotc -s -f data/benchmark2.txt -d 6 -g) > out.txt 2> benchmarks/timings/d6/parallel_gpu-$n.txt

  echo "difficulty 7"
  (time ./gotc -s -f data/benchmark1.txt -d 7 -g) > out.txt 2> benchmarks/timings/d7/parallel_gpu-$n.txt
}

make

# n=1
# run_timings
# n=2
# run_timings
# n=3
# run_timings

run_hashrate
