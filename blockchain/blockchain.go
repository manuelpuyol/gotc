package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"gotc/utils"
	"strings"
	"sync"
)

type Blockchain struct {
	Difficulty int         // number of leading 0s needed
	Head       *Block      // first block
	Tail       *Block      // last block
	NBlocks    uint        // how many blocks in the blockchain
	mutex      *sync.Mutex // mutex to sync Add
}

func NewBlockchain(difficulty int) *Blockchain {
	var mutex sync.Mutex
	return &Blockchain{difficulty, nil, nil, 0, &mutex}
}

func (bc *Blockchain) AddBlock(b *Block) bool {
	// Dont even lock if block is invalid
	if !bc.blockValid(b) {
		return false
	}

	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	// Verify again since `Lock` could have blocked and BC changed in the meantime
	if !bc.blockValid(b) {
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

func (bc *Blockchain) blockValid(b *Block) bool {
	return b.Header.Prev == bc.LastHash()
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
	j, err := json.MarshalIndent(bc.ToJSON(), "", "  ")
	utils.CheckErr(err)
	fmt.Println(string(j))
}
