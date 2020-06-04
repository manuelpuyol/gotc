package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"gotc/block"
)

type Blockchain struct {
	Difficulty uint
	Head       *block.Block
	Tail       *block.Block
	NBlocks    uint
}

func NewBlockchain(difficulty uint) *Blockchain {
	return &Blockchain{difficulty, nil, nil, 0}
}

func (bc *Blockchain) AddBlock(b *block.Block) bool {
	var lastHash [sha256.Size]byte

	if bc.Head != nil {
		lastHash = bc.Tail.Header.Hash
	}

	if b.Header.Prev != lastHash {
		return false
	}

	if bc.Head == nil {
		bc.Head = b
	}
	if bc.Tail != nil {
		bc.Tail.Next = b
	}
	bc.Tail = b
	bc.NBlocks++

	return true
}

func (bc *Blockchain) ToJSON() map[string]interface{} {
	var blocksJSON []map[string]interface{}
	curr := bc.Head
	for curr != nil {
		blocksJSON = append(blocksJSON, curr.ToJSON())
		curr = curr.Next
	}

	return map[string]interface{}{
		"difficulty": bc.Difficulty,
		"nblocks":    bc.NBlocks,
		"blocks":     blocksJSON,
	}
}

func (bc *Blockchain) Print() {
	j, _ := json.MarshalIndent(bc.ToJSON(), "", "  ")
	fmt.Println(string(j))
}
