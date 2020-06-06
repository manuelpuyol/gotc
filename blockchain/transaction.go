package blockchain

import (
	"encoding/json"
	"fmt"
	"gotc/hash"
	"gotc/utils"
	"strconv"
)

type Transaction struct {
	Value    uint32 `json:"value"`    // value of the transaction
	Sender   string `json:"sender"`   // sender hash
	Receiver string `json:"receiver"` // receiver hash
	Hash     string // transaction hash
}

func NewTransaction(value uint32) *Transaction {
	t := Transaction{
		Value:    value,
		Sender:   hash.ByteHash(utils.RandomBytes()),
		Receiver: hash.ByteHash(utils.RandomBytes()),
	}

	t.setHash()

	return &t
}

func NewTransactionFromJSON(bytes []byte) *Transaction {
	t := Transaction{}
	err := json.Unmarshal(bytes, &t)

	utils.CheckErr(err)

	t.setHash()

	return &t
}

func (t *Transaction) setHash() {
	t.Hash = hash.StrHash(t.Sender + t.Receiver + strconv.FormatUint(uint64(t.Value), 10))
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
	j, err := json.MarshalIndent(t.ToJSON(), "", "  ")
	utils.CheckErr(err)
	fmt.Println(string(j))
}
