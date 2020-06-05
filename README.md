## structs

- Blockchain
  - has multiple blocks (parallel linked list)
  - holds the number of blocks in the chain
  - holds the number of leading 0s necessary (difficulty)

- Block
  - has a header
  - has multiple transactions (can only have up to 5 transactions)
  - holds the number of transactions

- Header
  - holds the nonce (used for mining)
  - holds the previous block hash
  - holds the merkle tree root (calculated from the block transactions)

- Transaction
  - holds the transaction value
  - holds the sender hash (random sha256 for now)
  - holds the receiver hash (random sha256 for now)
  - holds the transaction hash (calculated doing a sha256 of sender + receiver + value)

- Miner
  - mines (more on that later)

- CLI
  - user can create a transaction (should be stored in a file)
  - user can start the miner
  - user can dump the blockchain (should be stored in a file)
  - user can change the number of leading 0s
  - user can change the number of threads (1 for serial, n for parallel)

## mining

1. Get available transactions
2. Create a permutation of them
3. Calculate merkle root
4. Create a block header with merkle root, previous block hash, and arbitrary nonce
5. Run a double sha256
6. Check if the calculated hash has `n` leading 0s (n configured by the user)
7. If yes -> return block
8. If not and nonce not reached max -> go to step 4
9. If reached max nonce -> go to step 2

Here, I'll parallelize from step 4, after getting the merkle root, we can divide MAX_NONCE into N chunks and each goroutine will try to mine using their chunk of nonces.
Then, we can have a barrier to wait for all routines to test their nonces before trying next permutation.

**Bonus**: I'd like to try to implement the miner for the GPU, since it would run a lot faster 
