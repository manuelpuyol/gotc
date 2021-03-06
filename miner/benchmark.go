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

	hashrate := int(float64(constants.MaxUint32) / seconds)

	if gpu {
		fmt.Println("GPU -", hashrate, "Hashes per second")
	} else if threads == 0 {
		fmt.Println("Serial -", hashrate, "Hashes per second")
	} else {
		fmt.Println(threads, "Threads -", hashrate, "Hashes per second")
	}
}

func benchmarkSerial()    { benchmarkMiner(0, false) }
func benchmark4Threads()  { benchmarkMiner(4, false) }
func benchmark8Threads()  { benchmarkMiner(8, false) }
func benchmark12Threads() { benchmarkMiner(12, false) }
func benchmarkGPU()       { benchmarkMiner(0, true) }

func BenchmarkAll(gpu bool) {
	if gpu {
		benchmarkGPU()
	} else {
		benchmarkSerial()
		benchmark4Threads()
		benchmark8Threads()
		benchmark12Threads()
	}
}
