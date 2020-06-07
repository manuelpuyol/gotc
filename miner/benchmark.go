package miner

// Couldn't make test work with C

import (
	"fmt"
	"gotc/blockchain"
	"gotc/constants"
	"gotc/sync"
	"time"
)

func newMiner(threads int, gpu bool) *Miner {
	bc := blockchain.NewBlockchain(60)

	return &Miner{
		bc:      bc,
		prev:    "",
		found:   constants.NotFound,
		nonce:   0,
		gpu:     gpu,
		id:      0,
		barrier: sync.NewBarrier(threads),
	}
}

func benchmarkMiner(threads int, gpu bool) {
	// hashes to run
	constants.MaxUint32 = 2048 * 2048

	m := newMiner(threads, gpu)

	start := time.Now()
	m.Check("prefix")
	end := time.Now()
	seconds := end.Sub(start).Seconds()

	fmt.Println("Took", seconds, "seconds")
	fmt.Println(constants.MaxUint32/uint32(seconds), "Hashes per second")
}

func BenchmarkSerial()   { benchmarkMiner(0, false) }
func Benchmark1Thread()  { benchmarkMiner(1, false) }
func Benchmark2Threads() { benchmarkMiner(2, false) }
func Benchmark4Threads() { benchmarkMiner(4, false) }
func Benchmark6Threads() { benchmarkMiner(6, false) }
func Benchmark8Threads() { benchmarkMiner(8, false) }
func BenchmarkGPU()      { benchmarkMiner(0, true) }
