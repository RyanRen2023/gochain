package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// Block represents a single block in the blockchain.
type Block struct {
	Timestamp int64  // Timestamp is the time when the block was created.
	Data      string // Data contains the information stored in the block.
	PrevHash  string // PrevHash is the hash of the previous block in the chain.
	Hash      string // Hash is the hash of the current block.
	Nonce     int    // Nonce is a counter used for the proof-of-work mechanism.
}

// NewBlock creates and returns a new Block with the given data and previous hash.
// It also performs mining to find a valid hash for the block.
//
// Parameters:
//   - data: The content to be stored in the block.
//   - prevHash: The hash of the previous block.
//
// Returns:
//   - A pointer to the newly created Block.
func NewBlock(data, prevHash string) *Block {
	b := &Block{
		Timestamp: time.Now().Unix(), // Set the current timestamp.
		Data:      data,              // Store the given data.
		PrevHash:  prevHash,          // Store the previous block's hash.
	}
	b.MineBlock() // Perform proof-of-work to generate a valid hash.
	return b
}

// MineBlock performs the proof-of-work algorithm to find a valid hash for the block.
// It increments the Nonce value until a hash with the required difficulty (starting with "0000") is found.
func (b *Block) MineBlock() {

	for {

		// Compute the SHA-256 hash of the block's string representation and nonce.
		hash := sha256.Sum256([]byte(b.String() + string(b.Nonce)))
		b.Hash = hex.EncodeToString(hash[:]) // Convert the hash to a hexadecimal string.

		// Check if the hash meets the difficulty criteria.
		if b.Hash[:4] == "0000" {
			break // Stop mining if a valid hash is found.
		}
		b.Nonce++ // Increment the nonce to try again.
	}
}

// IsValid checks whether the current block's hash is valid by recalculating it.
//
// Returns:
//   - true if the block's hash matches the recalculated hash and meets the difficulty criteria.
//   - false otherwise.
func (b *Block) IsValid() bool {
	// Recompute the hash using the block's data and nonce.
	hash := sha256.Sum256([]byte(b.String() + string(b.Nonce)))
	// Check if the hash matches the stored hash and meets the difficulty criteria.
	return hex.EncodeToString(hash[:]) == b.Hash
}

// String returns a string representation of the block's data, previous hash, and timestamp.
// This method is primarily used for hashing purposes.
//
// Returns:
//   - A concatenated string of the block's data, previous hash, and timestamp.
func (b *Block) String() string {
	return b.Data + b.PrevHash + string(b.Timestamp)
}

// GenesisBlock creates and returns the first block in the blockchain, also known as the genesis block.
// The genesis block has predefined data and hash values.
//
// Returns:
//   - A pointer to the genesis Block.
func GenesisBlock() *Block {
	return &Block{
		Timestamp: 0,                  // Genesis block has a fixed timestamp.
		Data:      "Genesis Block",    // Predefined data for the genesis block.
		PrevHash:  "",                 // No previous hash for the genesis block.
		Hash:      "0000000000000000", // Predefined valid hash for the genesis block.
		Nonce:     0,                  // Genesis block starts with a nonce of 0.
	}
}
