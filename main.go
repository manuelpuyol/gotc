package main

import (
	"flag"
	"fmt"
	"gotc/blockchain"
	"gotc/constants"
	"gotc/miner"
	"gotc/utils"
)

type Flags struct {
	difficulty int
	inPath     string
	outPath    string
	miners     int
	threads    int
	silent     bool
	gpu        bool
}

func main() {
	benchmark := flag.Bool("b", false, "Run hashing benchmark")
	difficulty := flag.Int("d", 5, "The number of trailing 0s needed for a block to be valid")
	inPath := flag.String("f", "data/transactions.txt", "Path to the file which contains the transactions to be read")
	outPath := flag.String("o", "data/blockchain.json", "Path to output the resulting blockchain")
	miners := flag.Int("m", 1, "The number of miners to spawn")
	threads := flag.Int("p", 0, "The number of threads for each miner to run, defaults to 0 (serial implementation).")
	silent := flag.Bool("s", false, "Enable silent mode (disable output)")
	gpu := flag.Bool("g", false, "Enable GPU")

	flag.Parse()

	f := &Flags{
		*difficulty,
		*inPath,
		*outPath,
		*miners,
		*threads,
		*silent,
		*gpu,
	}

	if *benchmark {
		runBenchmark()
	} else {
		run(f)
	}
}

func run(f *Flags) {
	constants.Silent = f.silent

	if f.miners == 0 {
		fmt.Println("Need at least one miner")
		return
	}

	if f.threads == 0 && f.miners > 1 {
		fmt.Println("Can't run miner pool in sequential mode")
		return
	}

	bc := blockchain.NewBlockchain(f.difficulty)
	pool := miner.NewPool(f.miners, f.threads, f.inPath, f.outPath, f.gpu, bc)
	pool.Prepare()

	if !constants.Silent {
		go utils.Spinner("Mining...")
	}

	minedAll := pool.Process()

	if !constants.Silent {
		if !minedAll {
			fmt.Println("\nCould not find blocks for the following transactions:")

			for pool.Queue.Size > 0 {
				pool.Queue.Dequeue().Print()
			}
		}

		pool.Finish()
	}
}

func runBenchmark() {
	miner.BenchmarkSerial()
	miner.Benchmark1Thread()
	miner.Benchmark2Threads()
	miner.Benchmark4Threads()
	miner.Benchmark6Threads()
	miner.Benchmark8Threads()
	miner.BenchmarkGPU()
}
