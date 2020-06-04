package merkle

import (
	"crypto/sha256"
	"gotc/transaction"
	"gotc/utils"
)

type Tree struct {
	leaves [][sha256.Size]byte
}

func NewTree(transactions []*transaction.Transaction) *Tree {
	var leaves [][sha256.Size]byte

	for _, t := range transactions {
		leaves = append(leaves, t.Hash)
	}

	return &Tree{leaves}
}

func (mt *Tree) GetRoot() [sha256.Size]byte {
	parents := make([][sha256.Size]byte, len(mt.leaves))
	copy(parents, mt.leaves)

	for len(parents) > 1 {
		parents = calculateNextLevel(parents)
	}

	return parents[0]
}

func calculateNextLevel(nodes [][sha256.Size]byte) [][sha256.Size]byte {
	var parents [][sha256.Size]byte

	size := len(nodes)

	for i := 0; i < size; i += 2 {
		val := utils.SHAToString(nodes[i])

		if i+i < size {
			val += utils.SHAToString(nodes[i+1])
		}

		parents = append(parents, sha256.Sum256([]byte(val)))
	}

	return parents
}
