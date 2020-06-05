package block

import (
	"encoding/json"
	"fmt"
	"gotc/header"
	"gotc/transaction"
	"gotc/utils"
)

type Block struct {
	Header        *header.Header             `json:"header"`
	Transactions  []*transaction.Transaction `json:"transactions"`
	NTransactions uint                       `json:"ntransactions"`
	Next          *Block
}

func NewBlock(h *header.Header, transactions []*transaction.Transaction) *Block {
	return &Block{h, transactions, uint(len(transactions)), nil}
}

func (b *Block) ToJSON() map[string]interface{} {
	var transactionsJSON []map[string]interface{}
	for _, t := range b.Transactions {
		transactionsJSON = append(transactionsJSON, t.ToJSON())
	}

	return map[string]interface{}{
		"header":        b.Header.ToJSON(),
		"transactions":  transactionsJSON,
		"ntransactions": b.NTransactions,
	}
}

func (b *Block) Print() {
	j, err := json.MarshalIndent(b.ToJSON(), "", "  ")
	utils.CheckErr(err)
	fmt.Println(string(j))
}
