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

	if gpu {
		fmt.Println("GPU -", constants.MaxUint32/uint32(seconds), "Hashes per second")
	} else if threads == 0 {
		fmt.Println("Serial -", constants.MaxUint32/uint32(seconds), "Hashes per second")
	} else {
		fmt.Println(threads, "Threads -", constants.MaxUint32/uint32(seconds), "Hashes per second")
	}
}

func benchmarkSerial()   { benchmarkMiner(0, false) }
func benchmark2Threads() { benchmarkMiner(2, false) }
func benchmark4Threads() { benchmarkMiner(4, false) }
func benchmark6Threads() { benchmarkMiner(6, false) }
func benchmark8Threads() { benchmarkMiner(8, false) }
func benchmarkGPU()      { benchmarkMiner(0, true) }

func BenchmarkAll() {
	benchmarkSerial()
	benchmark2Threads()
	benchmark4Threads()
	benchmark6Threads()
	benchmark8Threads()
	benchmarkGPU()
}
