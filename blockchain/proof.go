package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
)

// ProofOfWork performs the Proof-of-Work algorithm to find a hash
// that meets the required difficulty level. It repeatedly increments
// the nonce until a valid hash is found.
//
// Parameters:
//   - data: The input data to hash.
//   - difficulty: The required number of leading characters in the hash that must match a specific prefix.
//
// Returns:
//   - hash: The valid hash that satisfies the difficulty criteria.
//   - nonce: The nonce value that produces the valid hash.
func ProofOfWork(data string, difficulty int) (string, int) {
	var nonce int                              // Nonce is a counter used for the Proof-of-Work process.
	var hash string                            // Stores the resulting hash.
	prefix := string(make([]rune, difficulty)) // Generate a prefix of null characters (e.g., "\u0000") based on the difficulty level.

	for {
		// Calculate the hash using the given data and nonce.
		hash = CalculateHash(data, nonce)

		// Check if the hash satisfies the difficulty criteria.
		if hash[:difficulty] == prefix {
			break // Exit the loop when a valid hash is found.
		}
		nonce++ // Increment the nonce and try again.
	}
	return hash, nonce // Return the valid hash and nonce.
}

// CalculateHash computes a SHA-256 hash for the input data and nonce.
// It concatenates the data and nonce, hashes the result, and converts it to a hexadecimal string.
//
// Parameters:
//   - data: The input data to hash.
//   - nonce: The nonce value to include in the hash calculation.
//
// Returns:
//   - A hexadecimal string representation of the computed hash.
func CalculateHash(data string, nonce int) string {
	hash := sha256.Sum256([]byte(data + string(nonce))) // Combine data and nonce, then compute SHA-256 hash.
	return hex.EncodeToString(hash[:])                  // Convert the hash to a hexadecimal string.
}
