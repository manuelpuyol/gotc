package miner

import (
	"fmt"
	"gotc/blockchain"
	"gotc/constants"
	"gotc/hash"
	"gotc/merkle"
	"gotc/sync"
	"strconv"
	"sync/atomic"

	"github.com/gitchander/permutation"
)

type Miner interface {
	Mine() bool
	Reset(t []*blockchain.Transaction)
}

type CPUMiner struct {
	transactions []*blockchain.Transaction
	bc           *blockchain.Blockchain
	prev         string
	found        int32
	nonce        uint64
	id           int
	barrier      *sync.Barrier
}

func NewCPUMiner(bc *blockchain.Blockchain, threads, id int) Miner {
	return &CPUMiner{
		bc:      bc,
		prev:    bc.LastHash(),
		found:   constants.NotFound,
		nonce:   0,
		id:      id,
		barrier: sync.NewBarrier(threads),
	}
}

func (m *CPUMiner) Reset(t []*blockchain.Transaction) {
	m.transactions = t
	m.prev = m.bc.LastHash()
	m.found = constants.NotFound
	m.nonce = 0

	m.barrier.Reset()
}

func (m *CPUMiner) Mine() bool {
	p := permutation.New(blockchain.Slice(m.transactions))

	for m.found == constants.NotFound && p.Next() {
		m.checkPermutation()
	}

	if m.found == constants.Found {
		if m.sendBlock() {
			return true
		}

		// another miner beat me to it, so I have to mine again
		m.Reset(m.transactions)
		return m.Mine()
	}

	return false
}

func (m *CPUMiner) sendBlock() bool {
	mt := merkle.NewTree(m.transactions)
	h := blockchain.NewHeader(m.nonce, m.prev, mt.GetRoot())

	res := m.bc.AddBlock(blockchain.NewBlock(h, m.transactions))

	if res {
		fmt.Println("\nMiner", m.id, "found a block")
		fmt.Println("Nonce = ", m.nonce)
		fmt.Println("Hash = ", h.Hash)
		fmt.Println("Prev = ", h.Prev)
	}

	return res
}

func (m *CPUMiner) checkPermutation() {
	mt := merkle.NewTree(m.transactions)
	root := mt.GetRoot()
	prefix := m.prev + root

	if m.barrier.Threads > 0 {
		m.barrier.Start()
		var id uint64
		for id = 0; id < uint64(m.barrier.Threads); id++ {
			go findNonce(id, m, prefix)
		}
		m.barrier.Wait()
	} else {
		findNonce(0, m, prefix)
	}
}

func findNonce(id uint64, m *CPUMiner, prefix string) {
	bucket := constants.MaxUint64

	if m.barrier.Threads > 0 {
		bucket /= uint64(m.barrier.Threads)
	}

	nonce := id * bucket

	var max uint64
	if m.barrier.Threads == 0 || id == uint64(m.barrier.Threads-1) {
		max = constants.MaxUint64
	} else {
		max = (id + 1) * bucket
	}

	h := hash.NewHash(m.bc.Difficulty)

	for nonce < max && m.found == constants.NotFound {
		test := prefix + strconv.FormatUint(nonce, 10)

		if h.IsValid(test) {
			if atomic.CompareAndSwapInt32(&m.found, constants.NotFound, constants.Found) {
				m.nonce = nonce
			}
		}

		nonce++
	}

	if m.barrier.Threads > 0 {
		m.barrier.Done()
	}
}
