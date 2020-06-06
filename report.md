# Report

This was specifically for MPCS 52060 grading.
For this report, I'm assuming a significant difficulty on the blockchain

## The problem

As described in the `README`, this project goal is to create a simplified version of a blockchain with a miner.

Mining blocks is a very well known problem which uses lots of computational process. It is also a great way of showing the impact of parallelism on non-deterministic problems.

## Parallel solution

First, please make sure to read the  `README` to have a general idea of the project and mining algorithm.

There are two main parallel components:

1. MiningPool
2. Miner

### MiningPool

The MiningPool manages multiple miners working on the same blockchain.
It starts with a list of transactions to be processed and stores them on a `Parallel Queue`, which will be read by miners.

The Pool will spawn N miners. Each one of them will read up to 5 transactions from the queue and start their mining process. After finding a valid block, the miner will try to get more transactions and start again.

A barrier is used to sync all miners, so the Pool will only return when all miners finish their job.

In case some transactions weren't able to be mined, the pool will try to shuffle them around and change the number of transactions each miner can read at a time. Then, miners will try to find blocks again. This process is done only a limited amount of times so it doesn't end on a infinite loop.

### Miner

As described in the algorithm, miners will first permute the transactions it got from the queue and then divide max Uint32 into N chunks, spawning N routines.

Each routine will test hashes while changing the nonce inside its chunk, until a valid hash is found or its chunk ends.

Again, we have a barrier here which waits for all routines to end before returning or trying the next permutation.

Since multiple miners are in play and someone may have already found a block, the miner verifies if its block could be added to the blockchain, if not (which means another block was inserted before it could finish), it will start the process again with updated parameters.

## Challenges

A big challenge here is **randomness**. Since the process deals with a double SHA256 hash, its outcome is unpredictable and the same input can return different (correct) results.
The parallelization is pretty straightforward, I believe that the complexity of the system comes with the algorithm implementation.

## Speedup

### Hotspots / Bottlenecks

Clearly, the hotspot of the program is the `Miner.findNonce` function, which runs a loop 4294967295 times (in case no block is found). It can be also called 120 times for a list of 5 transactions.

Even though it is parallelized in the program, it could still run very slowly since we are dealing with randomness. This is a kind of bottleneck which can't be actually removed, but it had a good speedup for with increasing difficulty on the blockchain.

### Limitations

Since the problem works with such large numbers, I believe that the limitation here is the computational power. If we had a machine with millions of cores, it could be fast, but the nature of the problem may still make it slow.
