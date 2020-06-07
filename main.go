package main

import (
	"flag"
	"fmt"
	"gotc/blockchain"
	"gotc/constants"
	"gotc/miner"
	"gotc/utils"
)

func main() {
	difficulty := flag.Int("d", 5, "The number of trailing 0s needed for a block to be valid")
	inPath := flag.String("f", "data/transactions.txt", "Path to the file which contains the transactions to be read")
	outPath := flag.String("o", "data/blockchain.json", "Path to output the resulting blockchain")
	miners := flag.Int("m", 1, "The number of miners to spawn")
	threads := flag.Int("p", 0, "The number of threads for each miner to run, defaults to 0 (serial implementation).")
	benchmark := flag.Bool("b", false, "Enable benchmark mode (disable output)")
	gpu := flag.Bool("g", false, "Enable GPU")

	flag.Parse()

	constants.Benchmark = *benchmark

	if *miners == 0 {
		fmt.Println("Need at least one miner")
		return
	}

	if *threads == 0 && *miners > 1 {
		fmt.Println("Can't run miner pool in sequential mode")
		return
	}

	bc := blockchain.NewBlockchain(*difficulty)
	pool := miner.NewPool(*miners, *threads, *inPath, *outPath, *gpu, bc)
	pool.Prepare()

	if !constants.Benchmark {
		go utils.Spinner("Mining...")
	}

	minedAll := pool.Process()

	if !constants.Benchmark {
		if !minedAll {
			fmt.Println("\nCould not find blocks for the following transactions:")

			for pool.Queue.Size > 0 {
				pool.Queue.Dequeue().Print()
			}
		}

		pool.Finish()
	}
}
