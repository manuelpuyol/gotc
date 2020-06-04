package miner

import (
	"crypto/sha256"
	"fmt"
	"gotc/block"
	"gotc/blockchain"
	"gotc/hash"
	"gotc/merkle"
	"gotc/transaction"
	"strconv"
	"sync"
	"sync/atomic"
)

const Uint64Size = 64 << (^uint(0) >> 64 & 1)
const MaxUint64 uint64 = 1<<Uint64Size - 1
const NotFound = -1
const Found = 1

type Miner interface {
	Mine() *block.Block
}

type CPUMiner struct {
	transactions []*transaction.Transaction
	prev         string
	difficulty   int
	threads      int
	found        int32
	nonce        uint64
	mutex        *sync.Mutex
	group        *sync.WaitGroup
}

func NewCPUMiner(transactions []*transaction.Transaction, bc *blockchain.Blockchain, threads int) Miner {
	var mutex sync.Mutex
	var group sync.WaitGroup

	return &CPUMiner{
		transactions,
		fmt.Sprintf("%x", bc.LastHash()),
		bc.Difficulty,
		threads,
		NotFound,
		0,
		&mutex,
		&group,
	}
}

func (m *CPUMiner) Mine() *block.Block {
	m.group.Add(m.threads)

	mt := merkle.NewTree(m.transactions)
	root := fmt.Sprintf("%x", mt.GetRoot())
	prefix := m.prev + root

	var id uint64
	for id = 0; id < uint64(m.threads); id++ {
		go findNonce(id, m, prefix)
	}
	m.group.Wait()
	return nil
}

func findNonce(id uint64, m *CPUMiner, prefix string) {
	bucket := MaxUint64 / uint64(m.threads)
	nonce := id * bucket
	var max uint64

	if id == uint64(m.threads-1) {
		max = MaxUint64
	} else {
		max = (id + 1) * bucket
	}

	h := hash.NewHash(m.difficulty)

	for nonce < max && m.found == NotFound {
		test := prefix + strconv.FormatUint(nonce, 10)

		if h.IsValid(test) {
			fmt.Println("Thread ", id, " found a nonce = ", nonce)
			fmt.Printf("Hash = %x\n", sha256.Sum256([]byte(test)))

			if atomic.CompareAndSwapInt32(&m.found, NotFound, Found) {
				m.nonce = nonce
			}
		}

		nonce++
	}

	m.group.Done()
}
