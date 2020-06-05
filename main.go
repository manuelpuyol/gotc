package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"gotc/blockchain"
	"gotc/constants"
	"gotc/miner"
	"gotc/utils"
	"math/rand"
	"os"
)

type CTX struct {
	m                    miner.Miner
	transactions         []*blockchain.Transaction
	missed               []*blockchain.Transaction
	bc                   *blockchain.Blockchain
	transactionsPerBlock int
	shuffles             int
	processed            int
	miners               int
}

func newCTX(difficulty, threads, miners int, inPath string) *CTX {
	var missed []*blockchain.Transaction
	bc := blockchain.NewBlockchain(difficulty)
	m := miner.NewCPUMiner(bc, threads)

	return &CTX{
		m:                    m,
		transactions:         readTransactions(inPath),
		missed:               missed,
		bc:                   bc,
		transactionsPerBlock: constants.MaxTransactionsPerBlock,
		shuffles:             0,
		processed:            0,
		miners:               miners,
	}
}

func main() {
	difficulty := flag.Int("d", 5, "The number of trailing 0s needed for a block to be valid")
	inPath := flag.String("f", "data/transactions.txt", "Path to the file which contains the transactions to be read")
	outPath := flag.String("o", "data/blockchain.json", "Path to output the resulting blockchain")
	miners := flag.Int("m", 1, "The number of miners to spawn")
	threads := flag.Int("p", 0, "The number of threads for each miner to run, defaults to 0 (serial implementation).")

	flag.Parse()

	ctx := newCTX(*difficulty, *threads, *miners, *inPath)

	go utils.Spinner("Mining...")

	minedAll := processTransactions(ctx)

	if !minedAll {
		fmt.Println("\nCould not find blocks for the following transactions:")

		for _, t := range ctx.transactions {
			t.Print()
		}
	}

	writeBlockchain(ctx.bc, *outPath)
}

func processTransactions(ctx *CTX) bool {
	processed := 0
	transactionsCount := len(ctx.transactions)

	for processed < transactionsCount {
		transactions := getTransactions(ctx)
		ctx.m.Reset(transactions)

		if !ctx.m.Mine() {
			ctx.missed = append(ctx.missed, transactions...)
		}
		processed += ctx.transactionsPerBlock
	}

	if len(ctx.missed) > 0 {
		return retryMissedTransactions(ctx)
	}

	return true
}

func retryMissedTransactions(ctx *CTX) bool {
	size := len(ctx.missed)
	maxShuffles := size * size

	if size > ctx.transactionsPerBlock && ctx.shuffles < maxShuffles {
		ctx.shuffles++
		return suffleAndProcess(ctx)
	}
	if ctx.transactionsPerBlock > 0 {
		ctx.shuffles = 0
		return splitAndProcess(ctx)
	}

	return false
}

func suffleAndProcess(ctx *CTX) bool {
	transactions := ctx.missed
	rand.Shuffle(len(transactions), func(i, j int) {
		transactions[i], transactions[j] = transactions[j], transactions[i]
	})

	var missed []*blockchain.Transaction
	ctx.transactions = transactions
	ctx.missed = missed

	return processTransactions(ctx)
}

func splitAndProcess(ctx *CTX) bool {
	transactions := ctx.missed
	var missed []*blockchain.Transaction
	ctx.transactions = transactions
	ctx.missed = missed

	ctx.transactionsPerBlock--

	if ctx.transactionsPerBlock == 0 {
		return false
	}

	return processTransactions(ctx)
}

func getTransactions(ctx *CTX) []*blockchain.Transaction {
	end := ctx.processed + ctx.transactionsPerBlock
	transactionsCount := len(ctx.transactions)

	if end > transactionsCount {
		end = transactionsCount
	}

	return ctx.transactions[ctx.processed:end]
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

func readTransactions(inPath string) []*blockchain.Transaction {
	file, err := os.Open(inPath)
	utils.CheckErr(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var transactions []*blockchain.Transaction
	for scanner.Scan() {
		t := blockchain.NewTransactionFromJSON(scanner.Bytes())
		transactions = append(transactions, t)
	}

	return transactions
}
