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
	val := BTCHash(test)
	cmp := val[0:h.difficulty]

	return cmp == h.challenge
}

func BTCHash(data string) string {
	return StrHash(StrHash(data))
}

func StrHash(data string) string {
	return SHAToString(sha256.Sum256([]byte(data)))
}

func ByteHash(data []byte) string {
	return SHAToString(sha256.Sum256(data))
}

func SHAToString(bytes [sha256.Size]byte) string {
	return fmt.Sprintf("%x", bytes)
}
