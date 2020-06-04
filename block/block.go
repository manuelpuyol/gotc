package block

import (
	"gotc/header"
	"gotc/transaction"
)

type Block struct {
	Header        *header.Header             `json:"header"`
	Transactions  []*transaction.Transaction `json:"transactions"`
	NTransactions uint                       `json:"ntransactions"`
}

func NewBlock(h *header.Header, transactions []*transaction.Transaction) *Block {
	return &Block{h, transactions, uint(len(transactions))}
}
