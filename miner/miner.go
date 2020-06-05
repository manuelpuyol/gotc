package miner

import (
	"fmt"
	"gotc/block"
	"gotc/blockchain"
	"gotc/hash"
	"gotc/header"
	"gotc/merkle"
	"gotc/transaction"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/gitchander/permutation"
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
		bc.LastHash(),
		bc.Difficulty,
		threads,
		NotFound,
		0,
		&mutex,
		&group,
	}
}

func (m *CPUMiner) Mine() *block.Block {
	p := permutation.New(transaction.Slice(m.transactions))

	for m.found == NotFound && p.Next() {
		m.checkPermutation()
	}

	if m.found == Found {
		mt := merkle.NewTree(m.transactions)
		h := header.NewHeader(m.nonce, m.prev, mt.GetRoot())
		return block.NewBlock(h, m.transactions)
	}

	return nil
}

func (m *CPUMiner) checkPermutation() {
	m.group.Add(m.threads)

	mt := merkle.NewTree(m.transactions)
	root := mt.GetRoot()
	prefix := m.prev + root

	var id uint64
	for id = 0; id < uint64(m.threads); id++ {
		go findNonce(id, m, prefix)
	}

	m.group.Wait()
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
			if atomic.CompareAndSwapInt32(&m.found, NotFound, Found) {
				fmt.Println("\nGoroutine ", id, " found a block")
				fmt.Println("Nonce = ", nonce)
				fmt.Println("Hash = ", hash.BTCHash(test))
				m.nonce = nonce
			}
		}

		nonce++
	}

	m.group.Done()
}
