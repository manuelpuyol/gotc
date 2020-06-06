# GoTC

A **simplified** blockchain implementation, inspired in BTC

## How to run

This project has a single main file and you can run it using

```
go run main.go
```

There are some parameters too:

```
-d int
      The number of trailing 0s needed for a block to be valid (default 5)
-f string
      Path to the file which contains the transactions to be read (default "data/transactions.txt")
-m int
      The number of miners to spawn (default 1)
-o string
      Path to output the resulting blockchain (default "data/blockchain.json")
-p int
      The number of threads for each miner to run, defaults to 0 (serial implementation).
```

## Dependencies

[Permutation package](https://github.com/gitchander/permutation) - to deal with the Transactions permutations on my miner.

## How it works

1. The program reads a list of transactions from a file and stores them in a `Queue`.
2. Miners get transactions from `Queue` and start the mining process
3. In case miners couldn't find  blocks for some transactions, shuffle the queue or try to get a different number of transactions

## Mining

Mining is the core process of this project. Here is thow this is implemented:

1. Miner receives a transaction list and the number of threads to spawn
2. Miner gets a permutation of the transaction list
3. Miner calculates the merkle root of the permutation
4. Miner divides the maximum Uint32 value into nthreads blocks
5. Miner spawns threads
6. Thread runs in a loop hashing values with a nonce in their allocated block
7. If a hash has the number of trailing 0s necessary, Miner is notified and threads stop their executions
8. If not, go to step 2
9. If block is found, add it to the blockchain
10. If add failed, try to mine again with new blockchain value

## Data Structure

As you may see, most of the data is on the `blockchain` package

### Blockchain

This is the core structure, which hold all the valid blocks and transactions.
It is implemented as a linked list, which only has the `Add` method.

Since this is a blockchain, we can only add to the end of the chain, so we have a single `Mutex` to handle parallelism on the blockchain `Tail`.

```go
type Blockchain struct {
	Difficulty int         // number of trailing 0s needed
	Head       *Block      // first block
	Tail       *Block      // last block
	NBlocks    uint        // how many blocks in the blockchain
	mutex      *sync.Mutex // mutex to sync Add
}
```

### Block

Block is a single unit in the blockchain which holds transactions and has a reference to the next block.

```go
type Block struct {
	Header        *Header        // information about the block
	Transactions  []*Transaction // list of transactions
	NTransactions uint           // number of transactions
	Next          *Block         // next block in the linked list
}
```

### Header

Header holds information about a block and a reference to the last block (hash)

```go
type Header struct {
	Nonce uint32 // number used for minig
	Prev  string // hash of the last block
	Root  string // merkle root of the transaction list
	Hash  string // hash of the block
}
```

### Transaction

Transaction holds information about BTC sent from an account to another
```go
type Transaction struct {
	Value    uint32 // value of the transaction
	Sender   string // sender hash
	Receiver string // receiver hash
	Hash     string // transaction hash
}
```
