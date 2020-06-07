#!/bin/bash

run_hashrate() {
  ./gotc -b > benchmarks/hashrate/1.txt
  ./gotc -b > benchmarks/hashrate/2.txt
  ./gotc -b > benchmarks/hashrate/3.txt
  ./gotc -b > benchmarks/hashrate/4.txt
  ./gotc -b > benchmarks/hashrate/5.txt
}

run_timings() {
  echo "difficulty 4"
  (time ./gotc -s -f data/benchmark4.txt -d 4) > out.txt 2> benchmarks/timings/d4/serial-$n.txt
  echo "Running with 1 cores"
  (time ./gotc -s -f data/benchmark4.txt -d 4 -p 1) > out.txt 2> benchmarks/timings/d4/parallel_1-$n.txt
  echo "Running with 4 cores"
  (time ./gotc -s -f data/benchmark4.txt -d 4 -p 4) > out.txt 2> benchmarks/timings/d4/parallel_4-$n.txt
  echo "Running with 8 cores"
  (time ./gotc -s -f data/benchmark4.txt -d 4 -p 8) > out.txt 2> benchmarks/timings/d4/parallel_8-$n.txt
  echo "Running with 12 cores"
  (time ./gotc -s -f data/benchmark4.txt -d 4 -p 12) > out.txt 2> benchmarks/timings/d4/parallel_12-$n.txt
  echo "Running with 16 cores"
  (time ./gotc -s -f data/benchmark4.txt -d 4 -p 16) > out.txt 2> benchmarks/timings/d4/parallel_16-$n.txt

  echo "difficulty 5"
  (time ./gotc -s -f data/benchmark3.txt -d 5) > out.txt 2> benchmarks/timings/d5/serial-$n.txt
  echo "Running with 1 cores"
  (time ./gotc -s -f data/benchmark3.txt -d 5 -p 1) > out.txt 2> benchmarks/timings/d5/parallel_1-$n.txt
  echo "Running with 4 cores"
  (time ./gotc -s -f data/benchmark3.txt -d 5 -p 4) > out.txt 2> benchmarks/timings/d5/parallel_4-$n.txt
  echo "Running with 8 cores"
  (time ./gotc -s -f data/benchmark3.txt -d 5 -p 8) > out.txt 2> benchmarks/timings/d5/parallel_8-$n.txt
  echo "Running with 12 cores"
  (time ./gotc -s -f data/benchmark3.txt -d 5 -p 12) > out.txt 2> benchmarks/timings/d5/parallel_12-$n.txt
  echo "Running with 16 cores"
  (time ./gotc -s -f data/benchmark3.txt -d 5 -p 16) > out.txt 2> benchmarks/timings/d5/parallel_16-$n.txt

  echo "difficulty 6"
  (time ./gotc -s -f data/benchmark2.txt -d 6) > out.txt 2> benchmarks/timings/d6/serial-$n.txt
  echo "Running with 1 cores"
  (time ./gotc -s -f data/benchmark2.txt -d 6 -p 1) > out.txt 2> benchmarks/timings/d6/parallel_1-$n.txt
  echo "Running with 4 cores"
  (time ./gotc -s -f data/benchmark2.txt -d 6 -p 4) > out.txt 2> benchmarks/timings/d6/parallel_4-$n.txt
  echo "Running with 8 cores"
  (time ./gotc -s -f data/benchmark2.txt -d 6 -p 8) > out.txt 2> benchmarks/timings/d6/parallel_8-$n.txt
  echo "Running with 12 cores"
  (time ./gotc -s -f data/benchmark2.txt -d 6 -p 12) > out.txt 2> benchmarks/timings/d6/parallel_12-$n.txt
  echo "Running with 16 cores"
  (time ./gotc -s -f data/benchmark2.txt -d 6 -p 16) > out.txt 2> benchmarks/timings/d6/parallel_16-$n.txt

  echo "difficulty 7"
  (time ./gotc -s -f data/benchmark1.txt -d 7) > out.txt 2> benchmarks/timings/d7/serial-$n.txt
  echo "Running with 1 cores"
  (time ./gotc -s -f data/benchmark1.txt -d 7 -p 1) > out.txt 2> benchmarks/timings/d7/parallel_1-$n.txt
  echo "Running with 4 cores"
  (time ./gotc -s -f data/benchmark1.txt -d 7 -p 4) > out.txt 2> benchmarks/timings/d7/parallel_4-$n.txt
  echo "Running with 8 cores"
  (time ./gotc -s -f data/benchmark1.txt -d 7 -p 8) > out.txt 2> benchmarks/timings/d7/parallel_8-$n.txt
  echo "Running with 12 cores"
  (time ./gotc -s -f data/benchmark1.txt -d 7 -p 12) > out.txt 2> benchmarks/timings/d7/parallel_12-$n.txt
  echo "Running with 16 cores"
  (time ./gotc -s -f data/benchmark1.txt -d 7 -p 16) > out.txt 2> benchmarks/timings/d7/parallel_16-$n.txt
}

make

rm -rf benchmarks/
mkdir benchmarks
mkdir benchmarks/hashrate/
mkdir benchmarks/timings/
mkdir benchmarks/timings/d4
mkdir benchmarks/timings/d5
mkdir benchmarks/timings/d6
mkdir benchmarks/timings/d7

n=1
run_timings
n=2
run_timings
n=3
run_timings

run_hashrate
