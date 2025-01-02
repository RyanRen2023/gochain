package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"gochain/network"
	"net/http"
	"sync"
)

func main() {
	// Command-line flags for port and peer address
	port := flag.Int("port", 8080, "Port for this node to listen on")
	peer := flag.String("peer", "", "Address of an existing peer node to register with")
	flag.Parse()

	fmt.Printf("Peer is %s\n", *peer)

	// Create a new node
	node := network.NewNode()

	// Synchronization to ensure the server starts before registration
	var wg sync.WaitGroup
	wg.Add(1)

	// Start the HTTP server for the node
	go func() {
		defer wg.Done() // Signal that the server is ready
		err := node.StartServer(*port)
		if err != nil {
			fmt.Printf("Error starting server on port %d: %v\n", *port, err)
		}
		fmt.Println("Server initialization completed.")

	}()
	fmt.Printf("Node started on port %d\n", *port)

	// Wait for the server to initialize
	wg.Wait()

	// If a peer address is provided, register this node with the peer
	fmt.Println("Checking peer flag...")

	if *peer != "" {
		fmt.Printf("Attempting to register with peer at %s...\n", *peer)

		// Register the node with the peer
		if err := registerPeer(*peer, node, *port); err != nil {
			fmt.Printf("Failed to register with peer: %v\n", err)
		} else {
			fmt.Println("Registration with peer successful.")
		}
	} else {
		fmt.Println("No peer specified.")
	}

	fmt.Println("Peer registration block completed.")

	// Prevent the main function from exiting
	select {}
}

// Helper function to register a node with a peer
func registerPeer(peerAddress string, thisNode *network.Node, port int) error {
	// Prepare the payload for registration
	payload, err := json.Marshal(thisNode.Serialize(port))
	fmt.Printf("Payload is %s\n", payload)
	if err != nil {
		return fmt.Errorf("failed to serialize node: %w", err)
	}

	// Send the POST request to the peer's `/register` endpoint
	resp, err := http.Post(fmt.Sprintf("%s/register", peerAddress), "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to send registration request: %w", err)
	}
	defer resp.Body.Close()

	// Handle the response
	if resp.StatusCode == http.StatusCreated {
		fmt.Printf("Successfully registered with peer %s\n", peerAddress)
	} else {
		return fmt.Errorf("registration failed with status: %s", resp.Status)
	}
	return nil
}
