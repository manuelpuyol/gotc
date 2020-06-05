package miner

import (
	"bufio"
	"gotc/blockchain"
	"gotc/constants"
	"gotc/queue"
	"gotc/utils"
	"os"
	"sync"
)

type PoolCTX struct {
	missed               []*blockchain.Transaction
	bc                   *blockchain.Blockchain
	transactionsPerBlock int
	shuffles             int
	processed            int
}

func newPoolCTX(bc *blockchain.Blockchain) *PoolCTX {
	var missed []*blockchain.Transaction

	return &PoolCTX{
		missed:               missed,
		bc:                   bc,
		transactionsPerBlock: constants.MaxTransactionsPerBlock,
		shuffles:             0,
		processed:            0,
	}
}

type Pool struct {
	miners []Miner
	size   int
	q      *queue.Queue
	ctx    PoolCTX
}

func NewPool(size, threads int, bc *blockchain.Blockchain) *Pool {
	var cond sync.Cond
	q := queue.NewQueue(&cond)

	var miners []Miner
	for i := 0; i < size; i++ {
		miners = append(miners, NewCPUMiner(bc, threads))
	}

	return &Pool{
		miners: miners,
		size:   size,
		q:      q,
	}
}

func (p *Pool) Prepare(inPath string) {
	file, err := os.Open(inPath)
	utils.CheckErr(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		t := blockchain.NewTransactionFromJSON(scanner.Bytes())
		p.q.Enqueue(t)
	}
}

func (p *Pool) Process() {
}
