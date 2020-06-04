package merkle

import (
	"crypto/sha256"
	"fmt"
	"gotc/transaction"
)

type Tree struct {
	leaves      [][sha256.Size]byte
	permutation int
}

func NewTree(transactions []*transaction.Transaction) *Tree {
	var leaves [][sha256.Size]byte

	for _, t := range transactions {
		leaves = append(leaves, t.Hash)
	}

	return &Tree{leaves, 0}
}

func (mt *Tree) GetRoot() [sha256.Size]byte {
	parents := make([][sha256.Size]byte, len(mt.leaves))
	copy(parents, mt.leaves)

	fmt.Println("parents len = ", len(parents))
	fmt.Println("leaves len = ", len(mt.leaves))
	for len(parents) > 1 {
		parents = calculateNextLevel(parents)
	}

	return parents[0]
}

func calculateNextLevel(nodes [][sha256.Size]byte) [][sha256.Size]byte {
	var parents [][sha256.Size]byte

	size := len(nodes)

	for i := 0; i < size; i += 2 {
		val := fmt.Sprintf("%x", nodes[i])

		if i+i < size {
			val += fmt.Sprintf("%x", nodes[i+1])
		}

		parents = append(parents, sha256.Sum256([]byte(val)))
	}

	return parents
}
