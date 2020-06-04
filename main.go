package main

import (
	"crypto/sha256"
	"fmt"
	"gotc/block"
	"gotc/header"
	"gotc/merkle"
	"gotc/transaction"
)

func main() {
	t := transaction.NewTransaction(10)
	t.Print()

	var prev [sha256.Size]byte
	root := sha256.Sum256([]byte("root"))
	h := header.NewHeader(1, prev, root)
	h.Print()

	var ts []*transaction.Transaction
	ts = append(ts, t)
	ts = append(ts, transaction.NewTransaction(20))
	b := block.NewBlock(h, ts)
	b.Print()

	// bc := blockchain.NewBlockchain(2)
	// bc.AddBlock(b)

	// h2 := header.NewHeader(2, h.Hash, root)
	// b2 := block.NewBlock(h2, ts)

	// bc.AddBlock(b2)
	// bc.Print()

	mt := merkle.NewTree(ts)
	fmt.Printf("%x", mt.GetRoot())
}
