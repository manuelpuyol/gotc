package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"gotc/blockchain"
	"gotc/miner"
	"gotc/transaction"
	"gotc/utils"
	"os"
	"time"
)

const SpinnerDelay = 100

type CTX struct {
	transactions []*transaction.Transaction
	bc           *blockchain.Blockchain
	threads      int
}

func newCTX(difficulty, threads int, inPath string) *CTX {
	return &CTX{
		transactions: readTransactions(inPath),
		bc:           blockchain.NewBlockchain(difficulty),
		threads:      threads,
	}
}

func main() {
	difficulty := flag.Int("d", 5, "The number of trailing 0s needed for a block to be valid (Default 5)")
	inPath := flag.String(
		"f",
		"data/transactions.txt",
		"Path to the file which contains the transactions to be read (Default data/transactions.txt)",
	)
	outPath := flag.String("o", "data/blockchain.json", "Path to output the resulting blockchain (Default data/blockchain.json)")
	threads := flag.Int("p", 0, "An optional flag to run the miner in its parallel version.")
	flag.Parse()

	ctx := newCTX(*difficulty, *threads, *inPath)

	processTransactions(ctx)
	writeBlockchain(ctx.bc, *outPath)
}

func processTransactions(ctx *CTX) {
	m := miner.NewCPUMiner(ctx.transactions, ctx.bc, ctx.threads)

	go spinner(SpinnerDelay * time.Millisecond)

	b := m.Mine()

	ctx.bc.AddBlock(b)
}

func writeBlockchain(bc *blockchain.Blockchain, path string) {
	file, err := os.Create(path)
	utils.CheckErr(err)
	defer file.Close()

	j, err := json.MarshalIndent(bc.ToJSON(), "", "  ")
	utils.CheckErr(err)

	_, err = file.Write(j)
	utils.CheckErr(err)
}

func readTransactions(inPath string) []*transaction.Transaction {
	file, err := os.Open(inPath)
	utils.CheckErr(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var transactions []*transaction.Transaction
	for scanner.Scan() {
		t := transaction.NewTransactionFromJSON(scanner.Bytes())
		transactions = append(transactions, t)
	}

	return transactions
}

func spinner(delay time.Duration) {
	for {
		for _, r := range `⣽⣾⣷⣯⣟⡿⢿⣻` {
			fmt.Printf("\r Mining... %c ", r)
			time.Sleep(delay)
		}
	}
}
