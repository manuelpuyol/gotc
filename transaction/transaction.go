package transaction

import (
	"encoding/json"
	"fmt"
	"gotc/hash"
	"gotc/utils"
	"strconv"
)

type Transaction struct {
	Value    uint64
	Sender   string
	Receiver string
	Hash     string
}

func NewTransaction(value uint64) *Transaction {
	sender := hash.ByteHash(utils.RandomBytes())
	receiver := hash.ByteHash(utils.RandomBytes())

	h := hash.StrHash(sender + receiver + strconv.FormatUint(value, 10))

	return &Transaction{value, sender, receiver, h}
}

func (t *Transaction) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"value":    t.Value,
		"sender":   t.Sender,
		"receiver": t.Receiver,
		"hash":     t.Hash,
	}
}

func (t *Transaction) Print() {
	j, _ := json.MarshalIndent(t.ToJSON(), "", "  ")
	fmt.Println(string(j))
}
