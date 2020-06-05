package miner

import (
	"bufio"
	"encoding/json"
	"gotc/blockchain"
	"gotc/queue"
	"gotc/utils"
	"math/rand"
	"os"
	"sync"
)

type Pool struct {
	size    int
	inPath  string
	outPath string
	miners  []Miner
	bc      *blockchain.Blockchain
	ctx     *PoolCTX
	barrier *BarrierCTX
	Queue   *queue.Queue
}

func NewPool(size, threads int, inPath, outPath string, bc *blockchain.Blockchain) *Pool {
	var cond sync.Cond
	q := queue.NewQueue(&cond)

	var miners []Miner
	for i := 0; i < size; i++ {
		miners = append(miners, NewCPUMiner(bc, threads, i))
	}

	return &Pool{
		size:    size,
		inPath:  inPath,
		outPath: outPath,
		miners:  miners,
		bc:      bc,
		ctx:     newPoolCTX(),
		barrier: newBarrierCTX(threads),
		Queue:   q,
	}
}

func (p *Pool) Prepare() {
	file, err := os.Open(p.inPath)
	utils.CheckErr(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		t := blockchain.NewTransactionFromJSON(scanner.Bytes())
		p.Queue.Enqueue(t)
	}
}

func (p *Pool) Process() bool {
	transactionsCount := p.Queue.Size

	m := p.miners[0]

	for p.ctx.processed < transactionsCount {
		transactions := p.getTransactions()
		m.Reset(transactions)

		if !m.Mine() {
			p.ctx.missed = append(p.ctx.missed, transactions...)
		}

		p.ctx.processed += p.ctx.transactionsPerBlock
	}

	if len(p.ctx.missed) > 0 {
		return p.retryMissedTransactions()
	}

	return true
}

func (p *Pool) Finish() {
	file, err := os.Create(p.outPath)
	utils.CheckErr(err)
	defer file.Close()

	j, err := json.MarshalIndent(p.bc.ToJSON(), "", "  ")
	utils.CheckErr(err)

	_, err = file.Write(j)
	utils.CheckErr(err)
}

func (p *Pool) getTransactions() []*blockchain.Transaction {
	var transactions []*blockchain.Transaction

	for i := 0; i < p.ctx.transactionsPerBlock; i++ {
		t := p.Queue.Dequeue()

		if t == nil {
			break
		}

		transactions = append(transactions, t)
	}

	return transactions
}

func (p *Pool) retryMissedTransactions() bool {
	size := len(p.ctx.missed)
	maxShuffles := size * size

	if size > p.ctx.transactionsPerBlock && p.ctx.shuffles < maxShuffles {
		p.ctx.shuffles++
		return p.suffleAndProcess()
	}
	if p.ctx.transactionsPerBlock > 0 {
		p.ctx.shuffles = 0
		return p.splitAndProcess()
	}

	return false
}

func (p *Pool) suffleAndProcess() bool {
	transactions := p.ctx.missed
	rand.Shuffle(len(transactions), func(i, j int) {
		transactions[i], transactions[j] = transactions[j], transactions[i]
	})

	var missed []*blockchain.Transaction
	p.ctx.missed = missed

	for _, t := range transactions {
		p.Queue.Enqueue(t)
	}

	return p.Process()
}

func (p *Pool) splitAndProcess() bool {
	p.ctx.transactionsPerBlock--

	if p.ctx.transactionsPerBlock == 0 {
		return false
	}

	transactions := p.ctx.missed
	var missed []*blockchain.Transaction
	p.ctx.missed = missed

	for _, t := range transactions {
		p.Queue.Enqueue(t)
	}

	return p.Process()
}
