package miner

import (
	"fmt"
	"gotc/block"
	"gotc/blockchain"
	"gotc/constants"
	"gotc/hash"
	"gotc/header"
	"gotc/merkle"
	"gotc/transaction"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/gitchander/permutation"
)

type Miner interface {
	Mine() bool
	Reset(t []*transaction.Transaction)
}

type CPUMiner struct {
	transactions []*transaction.Transaction
	bc           *blockchain.Blockchain
	prev         string
	threads      int
	found        int32
	nonce        uint64
	mutex        *sync.Mutex
	group        *sync.WaitGroup
}

func NewCPUMiner(bc *blockchain.Blockchain, threads int) Miner {
	var mutex sync.Mutex
	var group sync.WaitGroup

	return &CPUMiner{
		bc:      bc,
		prev:    bc.LastHash(),
		threads: threads,
		found:   constants.NotFound,
		nonce:   0,
		mutex:   &mutex,
		group:   &group,
	}
}

func (m *CPUMiner) Reset(t []*transaction.Transaction) {
	m.transactions = t
	m.prev = m.bc.LastHash()
	m.found = constants.NotFound
	m.nonce = 0
}

func (m *CPUMiner) Mine() bool {
	p := permutation.New(transaction.Slice(m.transactions))

	for m.found == constants.NotFound && p.Next() {
		m.checkPermutation()
	}

	if m.found == constants.Found {
		if m.sendBlock() {
			return true
		}

		// another miner beat me to it, so I have to mine again
		m.Reset(m.transactions)
		m.Mine()
	}

	return false
}

func (m *CPUMiner) sendBlock() bool {
	mt := merkle.NewTree(m.transactions)
	h := header.NewHeader(m.nonce, m.prev, mt.GetRoot())
	b := block.NewBlock(h, m.transactions)

	return m.bc.AddBlock(b)
}

func (m *CPUMiner) checkPermutation() {
	mt := merkle.NewTree(m.transactions)
	root := mt.GetRoot()
	prefix := m.prev + root

	if m.threads > 0 {
		m.group.Add(m.threads)
		var id uint64
		for id = 0; id < uint64(m.threads); id++ {
			go findNonce(id, m, prefix)
		}
		m.group.Wait()
	} else {
		findNonce(0, m, prefix)
	}
}

func findNonce(id uint64, m *CPUMiner, prefix string) {
	bucket := constants.MaxUint64

	if m.threads > 0 {
		bucket /= uint64(m.threads)
	}

	nonce := id * bucket

	var max uint64
	if m.threads == 0 || id == uint64(m.threads-1) {
		max = constants.MaxUint64
	} else {
		max = (id + 1) * bucket
	}

	h := hash.NewHash(m.bc.Difficulty)

	for nonce < max && m.found == constants.NotFound {
		test := prefix + strconv.FormatUint(nonce, 10)

		if h.IsValid(test) {
			if atomic.CompareAndSwapInt32(&m.found, constants.NotFound, constants.Found) {
				fmt.Println("\nGoroutine ", id, " found a block")
				fmt.Println("Nonce = ", nonce)
				fmt.Println("Hash = ", hash.BTCHash(test))
				m.nonce = nonce
			}
		}

		nonce++
	}

	if m.threads > 0 {
		m.group.Done()
	}
}
