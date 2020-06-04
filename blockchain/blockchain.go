package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"gotc/block"
	"sync"
)

type Blockchain struct {
	Difficulty uint
	Head       *block.Block
	Tail       *block.Block
	NBlocks    uint
	mutex      *sync.Mutex
}

func NewBlockchain(difficulty uint) *Blockchain {
	var mutex sync.Mutex
	return &Blockchain{difficulty, nil, nil, 0, &mutex}
}

func (bc *Blockchain) AddBlock(b *block.Block) bool {
	var lastHash [sha256.Size]byte

	if bc.Head != nil {
		lastHash = bc.Tail.Header.Hash
	}

	if b.Header.Prev != lastHash {
		return false
	}

	bc.mutex.Lock()
	if bc.Head == nil {
		bc.Head = b
	}
	if bc.Tail != nil {
		bc.Tail.Next = b
	}
	bc.Tail = b
	bc.NBlocks++
	bc.mutex.Unlock()

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
