package transaction

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"gotc/utils"
	"strconv"
)

type Transaction struct {
	Value    uint64
	Sender   [sha256.Size]byte
	Receiver [sha256.Size]byte
	Hash     [sha256.Size]byte
}

func NewTransaction(value uint64) *Transaction {
	sender := sha256.Sum256(utils.RandomBytes())
	receiver := sha256.Sum256(utils.RandomBytes())

	hash := sha256.Sum256(toBytes(sender, receiver, value))

	return &Transaction{value, sender, receiver, hash}
}

func (t *Transaction) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"value":    t.Value,
		"sender":   utils.SHAToString(t.Sender),
		"receiver": utils.SHAToString(t.Receiver),
		"hash":     utils.SHAToString(t.Hash),
	}
}

func (t *Transaction) Print() {
	j, _ := json.MarshalIndent(t.ToJSON(), "", "  ")
	fmt.Println(string(j))
}

func toBytes(sender, receiver [sha256.Size]byte, value uint64) []byte {
	sstr := utils.SHAToString(sender)
	rstr := utils.SHAToString(receiver)
	vstr := strconv.FormatUint(value, 10)

	str := sstr + rstr + vstr
	return []byte(str)
}
