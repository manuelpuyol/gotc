package blockchain

import (
	"crypto/sha256"
	"gotc/block"
	"gotc/header"
	"gotc/transaction"
)

type Blockchain struct {
	Difficulty uint `json:"difficulty"`
	Head       *block.Block
	Tail       *block.Block
	NBlocks    uint `json:"nblocks"`
}

func NewBlockchain(difficulty uint) *Blockchain {
	return &Blockchain{Difficulty: difficulty, NBlocks: 0}
}

func (bc *Blockchain) AddBlock(transactions []*transaction.Transaction, nonce uint, root [sha256.Size]byte) bool {
	var lastHash [sha256.Size]byte

	if bc.Head != nil {
		lastHash = bc.Tail.Header.Hash
	}

	header := header.NewHeader(nonce, lastHash, root)
	block := block.NewBlock(header, transactions)

	if bc.Head == nil {
		bc.Head = block
	}
	if bc.Tail != nil {
		bc.Tail.Next = block
	}
	bc.Tail = block
	bc.NBlocks++

	return true
}
