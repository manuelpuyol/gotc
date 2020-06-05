package merkle

import (
	"gotc/hash"
	"gotc/transaction"
)

type Tree struct {
	leaves []string
}

func NewTree(transactions []*transaction.Transaction) *Tree {
	var leaves []string

	for _, t := range transactions {
		leaves = append(leaves, t.Hash)
	}

	return &Tree{leaves}
}

func (mt *Tree) GetRoot() string {
	parents := make([]string, len(mt.leaves))
	copy(parents, mt.leaves)

	for len(parents) > 1 {
		parents = calculateNextLevel(parents)
	}

	return parents[0]
}

func calculateNextLevel(nodes []string) []string {
	var parents []string

	size := len(nodes)

	for i := 0; i < size; i += 2 {
		val := nodes[i]

		if i+i < size {
			val += nodes[i+1]
		}

		parents = append(parents, hash.StrHash(val))
	}

	return parents
}
