package miner

/*
#cgo LDFLAGS: -L${SRCDIR}/../ -lgpu

#include<stdlib.h>
u_int32_t cmine(const char *str, int difficulty);
*/
import "C"

import (
	"fmt"
	"gotc/blockchain"
	"gotc/constants"
	"gotc/hash"
	"gotc/merkle"
	"gotc/sync"
	"strconv"
	"sync/atomic"
	"unsafe"

	"github.com/gitchander/permutation"
)

type Miner struct {
	transactions []*blockchain.Transaction
	bc           *blockchain.Blockchain
	prev         string
	found        int32
	nonce        uint32
	gpu          bool
	id           int
	barrier      *sync.Barrier
}

func NewMiner(bc *blockchain.Blockchain, threads int, gpu bool, id int) *Miner {
	return &Miner{
		bc:      bc,
		prev:    bc.LastHash(),
		found:   constants.NotFound,
		nonce:   0,
		gpu:     gpu,
		id:      id,
		barrier: sync.NewBarrier(threads),
	}
}

func (m *Miner) Reset(t []*blockchain.Transaction) {
	m.transactions = t
	m.prev = m.bc.LastHash()
	m.found = constants.NotFound
	m.nonce = 0

	m.barrier.Reset()
}

func (m *Miner) Mine() bool {
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

func (m *Miner) sendBlock() bool {
	mt := merkle.NewTree(m.transactions)
	h := blockchain.NewHeader(m.nonce, m.prev, mt.GetRoot())

	res := m.bc.AddBlock(blockchain.NewBlock(h, m.transactions))

	if res && !constants.Silent {
		fmt.Println("\nMiner", m.id, "found a block")
		fmt.Println("Nonce = ", m.nonce)
		fmt.Println("Hash = ", h.Hash)
		fmt.Println("Prev = ", h.Prev)
	}

	return res
}

func (m *Miner) checkPermutation() {
	mt := merkle.NewTree(m.transactions)
	root := mt.GetRoot()
	prefix := m.prev + root

	if m.gpu {
		m.checkGPU(prefix)
	} else {
		m.checkCPU(prefix)
	}
}

func (m *Miner) checkGPU(prefix string) {
	str := C.CString(prefix)
	defer C.free(unsafe.Pointer(str))

	difficulty := C.int(m.bc.Difficulty)
	nonce := C.cmine(str, difficulty)

	h := hash.NewHash(m.bc.Difficulty)

	test := prefix + strconv.FormatUint(uint64(nonce), 10)

	// parallelism is on GPU, so no need to be atomic
	// also check if hash is correct, since GPU always returns a value
	if h.IsValid(test) {
		m.nonce = uint32(nonce)
		m.found = constants.Found
	}
}

func (m *Miner) checkCPU(prefix string) {
	if m.barrier.Threads > 0 {
		m.barrier.Start()
		var id uint32
		for id = 0; id < uint32(m.barrier.Threads); id++ {
			go findNonce(id, m, prefix)
		}
		m.barrier.Wait()
	} else {
		findNonce(0, m, prefix)
	}
}

func findNonce(id uint32, m *Miner, prefix string) {
	bucket := constants.MaxUint32

	if m.barrier.Threads > 0 {
		bucket /= uint32(m.barrier.Threads)
	}

	nonce := id * bucket

	var max uint32
	if m.barrier.Threads == 0 || id == uint32(m.barrier.Threads-1) {
		max = constants.MaxUint32
	} else {
		max = (id + 1) * bucket
	}

	h := hash.NewHash(m.bc.Difficulty)

	for nonce < max && m.found == constants.NotFound {
		test := prefix + strconv.FormatUint(uint64(nonce), 10)

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
