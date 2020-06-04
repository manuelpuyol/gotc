package transaction

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"gotc/rand"
)

type Transaction struct {
	Value    uint
	Sender   [sha256.Size]byte
	Receiver [sha256.Size]byte
	Hash     [sha256.Size]byte
}

func NewTransaction(value uint) *Transaction {
	sender := sha256.Sum256(rand.RandomBytes())
	receiver := sha256.Sum256(rand.RandomBytes())

	hash := sha256.Sum256(toBytes(sender, receiver, value))

	return &Transaction{value, sender, receiver, hash}
}

func (t *Transaction) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"value":    t.Value,
		"sender":   fmt.Sprintf("%x", t.Sender),
		"receiver": fmt.Sprintf("%x", t.Receiver),
		"hash":     fmt.Sprintf("%x", t.Hash),
	}
}

func (t *Transaction) Print() {
	j, _ := json.MarshalIndent(t.ToJSON(), "", "  ")
	fmt.Println(string(j))
}

func toBytes(sender, receiver [sha256.Size]byte, value uint) []byte {
	bytes := append(receiver[:], sender[:]...)
	return append(bytes, byte(value))
}
