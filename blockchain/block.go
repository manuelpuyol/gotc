package blockchain

import (
	"encoding/json"
	"fmt"
	"gotc/constants"
	"gotc/utils"
)

type Block struct {
	Header        *Header        `json:"header"`
	Transactions  []*Transaction `json:"transactions"`
	NTransactions uint           `json:"ntransactions"`
	Next          *Block
}

func NewBlock(h *Header, transactions []*Transaction) *Block {
	nTransactions := len(transactions)

	if nTransactions > constants.MaxTransactionsPerBlock {
		panic("Block with more transactions than allowed")
	}

	return &Block{h, transactions, uint(nTransactions), nil}
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
