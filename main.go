package main

import (
	"crypto/sha256"
	"fmt"
	"gotc/block"
	"gotc/blockchain"
	"gotc/header"
	"gotc/miner"
	"gotc/transaction"
	"time"
)

func spinner(delay time.Duration) {
	for {
		for _, r := range `⣽⣾⣷⣯⣟⡿⢿⣻` {
			fmt.Printf("\r Mining... %c ", r)
			time.Sleep(delay)
		}
	}
}

func mockInitialHeader() *header.Header {
	var prev [sha256.Size]byte
	root := sha256.Sum256([]byte("root"))
	return header.NewHeader(1, prev, root)
}

func setupBlockchain() *blockchain.Blockchain {
	bc := blockchain.NewBlockchain(6)

	t := transaction.NewTransaction(10)
	h := mockInitialHeader()

	ts := []*transaction.Transaction{t}
	b := block.NewBlock(h, ts)

	bc.AddBlock(b)

	return bc
}

func main() {
	bc := setupBlockchain()

	t1 := transaction.NewTransaction(20)
	t2 := transaction.NewTransaction(30)

	ts := []*transaction.Transaction{t1, t2}

	m := miner.NewCPUMiner(ts, bc, 4)

	go spinner(200 * time.Millisecond)

	b := m.Mine()

	if b != nil {
		bc.AddBlock(b)
	}

	bc.Print()
}
