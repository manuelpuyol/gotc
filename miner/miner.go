package miner

import (
	"fmt"
	"gotc/blockchain"
	"gotc/constants"
	"gotc/hash"
	"gotc/merkle"
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
	barrier      *BarrierCTX
}

func NewCPUMiner(bc *blockchain.Blockchain, threads, id int) Miner {
	return &CPUMiner{
		bc:      bc,
		prev:    bc.LastHash(),
		found:   constants.NotFound,
		nonce:   0,
		id:      id,
		barrier: newBarrierCTX(threads),
	}
}

func (m *CPUMiner) Reset(t []*blockchain.Transaction) {
	m.transactions = t
	m.prev = m.bc.LastHash()
	m.found = constants.NotFound
	m.nonce = 0

	m.barrier.counter = 0
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
		m.Mine()
	}

	return false
}

func (m *CPUMiner) sendBlock() bool {
	mt := merkle.NewTree(m.transactions)
	h := blockchain.NewHeader(m.nonce, m.prev, mt.GetRoot())

	return m.bc.AddBlock(blockchain.NewBlock(h, m.transactions))
}

func (m *CPUMiner) checkPermutation() {
	mt := merkle.NewTree(m.transactions)
	root := mt.GetRoot()
	prefix := m.prev + root

	if m.barrier.threads > 0 {
		m.barrier.group.Add(m.barrier.threads)
		var id uint64
		for id = 0; id < uint64(m.barrier.threads); id++ {
			go findNonce(id, m, prefix)
		}
		m.barrier.group.Wait()
	} else {
		findNonce(0, m, prefix)
	}
}

func findNonce(id uint64, m *CPUMiner, prefix string) {
	bucket := constants.MaxUint64

	if m.barrier.threads > 0 {
		bucket /= uint64(m.barrier.threads)
	}

	nonce := id * bucket

	var max uint64
	if m.barrier.threads == 0 || id == uint64(m.barrier.threads-1) {
		max = constants.MaxUint64
	} else {
		max = (id + 1) * bucket
	}

	h := hash.NewHash(m.bc.Difficulty)

	for nonce < max && m.found == constants.NotFound {
		test := prefix + strconv.FormatUint(nonce, 10)

		if h.IsValid(test) {
			if atomic.CompareAndSwapInt32(&m.found, constants.NotFound, constants.Found) {
				fmt.Println("\nMiner", m.id, "routine", id, "found a block")
				fmt.Println("Nonce = ", nonce)
				fmt.Println("Hash = ", hash.BTCHash(test))
				m.nonce = nonce
			}
		}

		nonce++
	}

	if m.barrier.threads > 0 {
		// Barrier
		m.barrier.mutex.Lock()
		m.barrier.counter++
		if m.barrier.counter == m.barrier.threads {
			m.barrier.cond.Broadcast()
		} else {
			for m.barrier.counter != m.barrier.threads {
				m.barrier.cond.Wait()
			}
		}
		m.barrier.mutex.Unlock()

		// Everyone finished
		m.barrier.group.Done()
	}
}
