#!/bin/bash
#
#SBATCH --mail-user=manuelpuyol@cs.uchicago.edu
#SBATCH --mail-type=ALL
#SBATCH --output=/home/manuelpuyol/slurm/out/%j.%N.stdout
#SBATCH --error=/home/manuelpuyol/slurm/out/%j.%N.stderr
#SBATCH --workdir=/home/manuelpuyol/slurm/gotcbenchmarks
#SBATCH --partition=titan
#SBATCH --job-name=benchmark_gotc
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --mem-per-cpu=500
#SBATCH --gres=gpu:1

git pull
module load golang/1.14.1
./gpu_benchmark.bash
