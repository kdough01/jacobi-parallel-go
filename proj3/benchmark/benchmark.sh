#!/bin/bash
#SBATCH --mail-user=kdough01@uchicago.edu
#SBATCH --mail-type=ALL
#SBATCH --job-name=proj3
#SBATCH --output=./slurm/out/%j.%N.stdout
#SBATCH --error=./slurm/out/%j.%N.stderr
#SBATCH --time=00:30:00
#SBATCH --partition=debug
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --cpus-per-task=16
#SBATCH --mem-per-cpu=900
#SBATCH --exclusive

module load python/3.10
module load golang/1.19

python benchmark.py