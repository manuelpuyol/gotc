package hash

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

type Hash struct {
	difficulty int
	challenge  string
}

func NewHash(difficulty int) *Hash {
	challenge := strings.Repeat("0", difficulty)

	return &Hash{difficulty, challenge}
}

func (h *Hash) IsValid(test string) bool {
	val := fmt.Sprintf("%x", sha256.Sum256([]byte(test)))
	cmp := val[0:h.difficulty]

	return cmp == h.challenge
}

func BTCHash(data []byte) [sha256.Size]byte {
	hash := fmt.Sprintf("%x", sha256.Sum256(data))
	return sha256.Sum256([]byte(hash))
}
