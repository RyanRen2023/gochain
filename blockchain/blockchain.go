package blockchain

import "sync"

// Blockchain represents the chain of blocks in the blockchain.
// It includes a list of blocks and a mutex for thread-safe operations.
type Blockchain struct {
	Blocks []*Block   // A slice of pointers to blocks, representing the blockchain
	nu     sync.Mutex // Mutex to ensure thread-safe operations
}

// NewBlockchain initializes a new blockchain with a genesis block.
// It returns a pointer to the newly created Blockchain instance.
func NewBlockchain() *Blockchain {
	return &Blockchain{
		Blocks: []*Block{GenesisBlock()}, // Initialize the blockchain with the genesis block
	}
}

// AddBlock appends a new block with the given data to the blockchain.
// It ensures thread-safety by locking the mutex during the operation.
func (bc *Blockchain) AddBlock(data string) {
	bc.nu.Lock()         // Lock the mutex to prevent concurrent writes
	defer bc.nu.Unlock() // Unlock the mutex after the operation is complete

	prevBlock := bc.Blocks[len(bc.Blocks)-1]   // Get the last block in the chain
	newBlock := NewBlock(data, prevBlock.Hash) // Create a new block with the given data and the previous block's hash
	bc.Blocks = append(bc.Blocks, newBlock)    // Append the new block to the blockchain
}

// IsValid verifies the integrity of the blockchain.
// It checks whether each block is correctly linked to the previous block
// and whether all blocks satisfy the proof-of-work condition.
func (bc *Blockchain) IsValid() bool {
	for i := 1; i < len(bc.Blocks); i++ { // Start from the second block (index 1)
		currBlock := bc.Blocks[i]   // Current block
		prevBlock := bc.Blocks[i-1] // Previous block

		// Check if the current block's PrevHash matches the previous block's hash
		// and if the current block's proof-of-work is valid
		if currBlock.PrevHash != prevBlock.Hash || !currBlock.IsValid() {
			return false // Blockchain is invalid
		}
	}

	return true // Blockchain is valid
}
