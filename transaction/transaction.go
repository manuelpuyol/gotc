package transaction

import (
	"crypto/sha256"
	"gotc/rand"
)

type Transaction struct {
	Value    uint              `json:"value"`
	Sender   [sha256.Size]byte `json:"sender"`
	Receiver [sha256.Size]byte `json:"receiver"`
	Hash     [sha256.Size]byte `json:"hash"`
}

func NewTransaction(value uint) *Transaction {
	sender := sha256.Sum256(rand.RandomBytes())
	receiver := sha256.Sum256(rand.RandomBytes())

	hash := sha256.Sum256(toBytes(sender, receiver, value))

	return &Transaction{value, sender, receiver, hash}
}

func toBytes(sender, receiver [sha256.Size]byte, value uint) []byte {
	bytes := append(receiver[:], sender[:]...)
	return append(bytes, byte(value))
}
