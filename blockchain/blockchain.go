package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"gotc/block"
	"strings"
	"sync"
)

type Blockchain struct {
	Difficulty int
	Head       *block.Block
	Tail       *block.Block
	NBlocks    uint
	mutex      *sync.Mutex
}

func NewBlockchain(difficulty int) *Blockchain {
	var mutex sync.Mutex
	return &Blockchain{difficulty, nil, nil, 0, &mutex}
}

func (bc *Blockchain) AddBlock(b *block.Block) bool {
	if b.Header.Prev != bc.LastHash() {
		return false
	}

	bc.mutex.Lock()
	defer bc.mutex.Unlock()

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

func (bc *Blockchain) LastHash() string {
	if bc.Head != nil {
		return bc.Tail.Header.Hash
	}
	return strings.Repeat("0", sha256.BlockSize)
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
