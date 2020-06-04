package miner

import (
	"crypto/sha256"
	"fmt"
	"gotc/block"
	"gotc/blockchain"
	"gotc/transaction"
	"sync"
)

const Uint64Size = 64 << (^uint(0) >> 64 & 1)
const MaxUint64 uint64 = 1<<Uint64Size - 1

type Miner interface {
	Mine() *block.Block
}

type CPUMiner struct {
	transactions []*transaction.Transaction
	prev         [sha256.Size]byte
	difficulty   uint
	threads      int
	found        bool
	mutex        *sync.Mutex
}

func NewCPUMiner(transactions []*transaction.Transaction, bc *blockchain.Blockchain, threads int) Miner {
	var mutex sync.Mutex

	return &CPUMiner{transactions, bc.LastHash(), bc.Difficulty, threads, false, &mutex}
}

func (m *CPUMiner) Mine() *block.Block {
	var group sync.WaitGroup
	group.Add(m.threads)

	var id uint64
	for id = 0; id < uint64(m.threads); id++ {
		go findNonce(id, m, &group)
	}
	group.Wait()
	return nil
}

func findNonce(id uint64, m *CPUMiner, group *sync.WaitGroup) {
	bucket := MaxUint64 / uint64(m.threads)
	nonce := id * bucket
	var max uint64

	if id == uint64(m.threads-1) {
		max = MaxUint64
	} else {
		max = (id + 1) * bucket
	}

	fmt.Println("Thread = ", id, " - Nonce = ", nonce, " - Max = ", max)
	group.Done()
}
