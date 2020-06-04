package hash

import (
	"crypto/sha256"
	"gotc/utils"
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
	val := utils.SHAToString(sha256.Sum256([]byte(test)))
	cmp := val[0:h.difficulty]

	return cmp == h.challenge
}

func BTCHash(data []byte) [sha256.Size]byte {
	hash := utils.SHAToString(sha256.Sum256(data))
	return sha256.Sum256([]byte(hash))
}
