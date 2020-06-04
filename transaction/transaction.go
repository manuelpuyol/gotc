package transaction

import (
	"crypto/sha256"
	"gotc/rand"
)

type Transaction struct {
	value    uint
	sender   [sha256.Size]byte
	receiver [sha256.Size]byte
	hash     [sha256.Size]byte
}

func NewTransaction(value uint) *Transaction {
	sender := sha256.Sum256(rand.RandomBytes())
	receiver := sha256.Sum256(rand.RandomBytes())

	bytes := append(receiver[:], sender[:]...)
	bytes = append(bytes, byte(value))

	hash := sha256.Sum256(bytes)

	return &Transaction{value, sender, receiver, hash}
}
