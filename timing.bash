#!/bin/bash

run_timings() {
  echo "difficulty 4"
  (time ./main -d 4) > out.txt 2> timings/d4/serial-$n.txt
  (time ./main -d 4 -p 1) > out.txt 2> timings/d4/parallel_1-$n.txt
  (time ./main -d 4 -p 2) > out.txt 2> timings/d4/parallel_2-$n.txt
  (time ./main -d 4 -p 4) > out.txt 2> timings/d4/parallel_4-$n.txt
  (time ./main -d 4 -p 6) > out.txt 2> timings/d4/parallel_6-$n.txt
  (time ./main -d 4 -p 8) > out.txt 2> timings/d4/parallel_8-$n.txt

  echo "difficulty 5"
  (time ./main -d 5) > out.txt 2> timings/d5/serial-$n.txt
  (time ./main -d 5 -p 1) > out.txt 2> timings/d5/parallel_1-$n.txt
  (time ./main -d 5 -p 2) > out.txt 2> timings/d5/parallel_2-$n.txt
  (time ./main -d 5 -p 4) > out.txt 2> timings/d5/parallel_4-$n.txt
  (time ./main -d 5 -p 6) > out.txt 2> timings/d5/parallel_6-$n.txt
  (time ./main -d 5 -p 8) > out.txt 2> timings/d5/parallel_8-$n.txt

  echo "difficulty 6"
  (time ./main -d 6) > out.txt 2> timings/d6/serial-$n.txt
  (time ./main -d 6 -p 1) > out.txt 2> timings/d6/parallel_1-$n.txt
  (time ./main -d 6 -p 2) > out.txt 2> timings/d6/parallel_2-$n.txt
  (time ./main -d 6 -p 4) > out.txt 2> timings/d6/parallel_4-$n.txt
  (time ./main -d 6 -p 6) > out.txt 2> timings/d6/parallel_6-$n.txt
  (time ./main -d 6 -p 8) > out.txt 2> timings/d6/parallel_8-$n.txt

  echo "difficulty 7"
  (time ./main -d 7) > out.txt 2> timings/d7/serial-$n.txt
  (time ./main -d 7 -p 1) > out.txt 2> timings/d7/parallel_1-$n.txt
  (time ./main -d 7 -p 2) > out.txt 2> timings/d7/parallel_2-$n.txt
  (time ./main -d 7 -p 4) > out.txt 2> timings/d7/parallel_4-$n.txt
  (time ./main -d 7 -p 6) > out.txt 2> timings/d7/parallel_6-$n.txt
  (time ./main -d 7 -p 8) > out.txt 2> timings/d7/parallel_8-$n.txt
}

rm -rf timings/
mkdir timings
mkdir timings/d4
mkdir timings/d5
mkdir timings/d6
mkdir timings/d7

n=1
run_timings
n=2
run_timings
n=3
run_timings
