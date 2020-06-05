package miner

import (
	"gotc/blockchain"
	"gotc/constants"
	"sync"
)

type BarrierCTX struct {
	mutex   *sync.Mutex
	cond    *sync.Cond
	group   *sync.WaitGroup
	counter int
	threads int
}

func newBarrierCTX(threads int) *BarrierCTX {
	var mutex sync.Mutex
	var group sync.WaitGroup
	cond := sync.NewCond(&mutex)

	return &BarrierCTX{
		mutex:   &mutex,
		group:   &group,
		cond:    cond,
		counter: 0,
		threads: threads,
	}
}

type PoolCTX struct {
	missed               []*blockchain.Transaction
	transactionsPerBlock int
	shuffles             int
	processed            int
}

func newPoolCTX() *PoolCTX {
	var missed []*blockchain.Transaction

	return &PoolCTX{
		missed:               missed,
		transactionsPerBlock: constants.MaxTransactionsPerBlock,
		shuffles:             0,
		processed:            0,
	}
}
