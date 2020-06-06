package blockchain

import (
	"encoding/json"
	"fmt"
	"gotc/constants"
	"gotc/utils"
)

type Block struct {
	Header        *Header        `json:"header"`        // information about the block
	Transactions  []*Transaction `json:"transactions"`  // list of transactions
	NTransactions uint           `json:"ntransactions"` // number of transactions
	Next          *Block         // next block in the linked list
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
