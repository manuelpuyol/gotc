package miner

import (
	"fmt"
	"gotc/blockchain"
	"gotc/constants"
	"gotc/hash"
	"gotc/merkle"
	"strconv"
	"sync"
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
	ctx          *MinerCTX
}

type MinerCTX struct {
	mutex   *sync.Mutex
	cond    *sync.Cond
	group   *sync.WaitGroup
	counter int
	threads int
}

func newMinerCTX(threads int) *MinerCTX {
	var mutex sync.Mutex
	var group sync.WaitGroup
	cond := sync.NewCond(&mutex)

	return &MinerCTX{
		mutex:   &mutex,
		group:   &group,
		cond:    cond,
		counter: 0,
		threads: threads,
	}
}

func NewCPUMiner(bc *blockchain.Blockchain, threads int) Miner {
	return &CPUMiner{
		bc:    bc,
		prev:  bc.LastHash(),
		found: constants.NotFound,
		nonce: 0,
		ctx:   newMinerCTX(threads),
	}
}

func (m *CPUMiner) Reset(t []*blockchain.Transaction) {
	m.transactions = t
	m.prev = m.bc.LastHash()
	m.found = constants.NotFound
	m.nonce = 0

	m.ctx.counter = 0
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

	if m.ctx.threads > 0 {
		m.ctx.group.Add(m.ctx.threads)
		var id uint64
		for id = 0; id < uint64(m.ctx.threads); id++ {
			go findNonce(id, m, prefix)
		}
		m.ctx.group.Wait()
	} else {
		findNonce(0, m, prefix)
	}
}

func findNonce(id uint64, m *CPUMiner, prefix string) {
	bucket := constants.MaxUint64

	if m.ctx.threads > 0 {
		bucket /= uint64(m.ctx.threads)
	}

	nonce := id * bucket

	var max uint64
	if m.ctx.threads == 0 || id == uint64(m.ctx.threads-1) {
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

	if m.ctx.threads > 0 {
		// Barrier
		m.ctx.mutex.Lock()
		m.ctx.counter++
		if m.ctx.counter == m.ctx.threads {
			m.ctx.cond.Broadcast()
		} else {
			for m.ctx.counter != m.ctx.threads {
				m.ctx.cond.Wait()
			}
		}
		m.ctx.mutex.Unlock()

		// Everyone finished
		m.ctx.group.Done()
	}
}
