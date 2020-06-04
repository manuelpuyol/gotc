package main

import (
	"crypto/sha256"
	"gotc/block"
	"gotc/blockchain"
	"gotc/header"
	"gotc/transaction"
)

func main() {
	t := transaction.NewTransaction(10)
	t.Print()

	prev := sha256.Sum256([]byte("prev"))
	root := sha256.Sum256([]byte("root"))
	h := header.NewHeader(1, prev, root)
	h.Print()

	var ts []*transaction.Transaction
	ts = append(ts, t)
	ts = append(ts, transaction.NewTransaction(20))
	b := block.NewBlock(h, ts)
	b.Print()

	bc := blockchain.NewBlockchain(2)
	bc.AddBlock(ts, 1, root)
	bc.AddBlock(ts, 2, root)
	bc.AddBlock(ts, 3, root)
	bc.Print()
}
