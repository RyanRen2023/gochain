package network

import (
	"encoding/json"
	"fmt"
	"gochain/blockchain"
	"net/http"
	"sync"
)

// Node represents a single node in the blockchain network.
type Node struct {
	Chain *blockchain.Blockchain // The blockchain instance for this node
	Peers []string               // List of peer node addresses
	Mu    sync.Mutex             // Mutex to ensure thread-safe operations
}

// NewNode creates and initializes a new node with its own blockchain instance.
func NewNode() *Node {
	return &Node{
		Chain: blockchain.NewBlockchain(), // Initialize the blockchain with a genesis block
		Peers: []string{},
	}
}

// AddPeer adds a new peer to this node's list of peers.
// Duplicate peers are ignored.
func (n *Node) AddPeer(peer string) {
	n.Mu.Lock()
	defer n.Mu.Unlock()

	// Check for duplicates
	for _, p := range n.Peers {
		if p == peer {
			return
		}
	}

	// Add the new peer
	n.Peers = append(n.Peers, peer)
}

// Sync synchronizes this node's blockchain with its peers.
func (n *Node) Sync() {
	n.Mu.Lock()
	defer n.Mu.Unlock()

	for _, peerAddress := range n.Peers {
		// Fetch the blockchain from the peer
		resp, err := http.Get(fmt.Sprintf("%s/blockchain", peerAddress))
		if err != nil {
			fmt.Printf("Failed to fetch blockchain from peer %s: %v\n", peerAddress, err)
			continue
		}
		defer resp.Body.Close()

		// Decode the peer's blockchain
		var peerBlocks []*blockchain.Block
		if err := json.NewDecoder(resp.Body).Decode(&peerBlocks); err != nil {
			fmt.Printf("Failed to decode blockchain from peer %s: %v\n", peerAddress, err)
			continue
		}

		// Replace this node's blockchain if the peer's is longer
		if len(peerBlocks) > len(n.Chain.Blocks) {
			n.Chain.Blocks = peerBlocks
			fmt.Printf("Synchronized with peer %s: updated blockchain\n", peerAddress)
		}
	}
}

// Serialize prepares a Node's data for JSON serialization.
func (n *Node) Serialize(port int) map[string]interface{} {
	n.Mu.Lock()
	defer n.Mu.Unlock()

	return map[string]interface{}{
		"address": fmt.Sprintf("http://localhost:%d", port),
	}
}

// PrintPeers prints the list of peers.
func (n *Node) PrintPeers() {
	n.Mu.Lock()
	defer n.Mu.Unlock()

	fmt.Println("Peers:")
	for _, peer := range n.Peers {
		fmt.Println(" -", peer)
	}
}
