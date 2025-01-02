package network

import (
	"encoding/json"
	"fmt"
	"gochain/blockchain"
	"io"
	"log"
	"net/http"
	"os"
)

// StartServer starts the HTTP server for the node.
// It sets up handlers for blockchain-related API endpoints and listens for incoming requests.
func (n *Node) StartServer(port int) error {
	// Set up logging to a file
	logFile, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return err
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Set up HTTP handlers for various endpoints
	http.HandleFunc("/blockchain", n.logMiddleware(n.handleGetBlockchain)) // Fetch the current blockchain
	http.HandleFunc("/block", n.logMiddleware(n.handleReceiveBlock))       // Receive a new block from peers
	http.HandleFunc("/register", n.logMiddleware(n.handleRegisterNode))    // Register a new peer node
	http.HandleFunc("/peers", n.logMiddleware(n.handleGetPeers))           // Fetch the list of peers

	// Start the HTTP server on the specified port
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting server on %s...\n", addr)

	go func() {
		err := http.ListenAndServe(addr, nil)
		if err != nil && err != http.ErrServerClosed {
			log.Printf("Server failed to start: %v\n", err)
		}
	}()

	// Log server readiness
	log.Println("Server started successfully and is running in the background.")
	return nil
}

// logMiddleware wraps an HTTP handler to log incoming requests and responses.
func (n *Node) logMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
		log.Printf("Request processed: %s %s", r.Method, r.URL.Path)
	}
}

// Handle GET /peers - Returns the list of registered peers.
func (n *Node) handleGetPeers(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling /peers request...")

	n.Mu.Lock()
	defer n.Mu.Unlock()

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Log the list of peers
	log.Printf("Current peers: %v", n.Peers)

	// Encode the list of peers and write it to the response
	if err := json.NewEncoder(w).Encode(n.Peers); err != nil {
		log.Printf("Error encoding peers: %v", err)
		http.Error(w, "Failed to fetch peers", http.StatusInternalServerError)
	}
	log.Println("Finished processing /peers request.")
}

// Handle GET /blockchain - Returns the current blockchain in JSON format.
func (n *Node) handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling /blockchain request...")

	n.Mu.Lock()
	defer n.Mu.Unlock()

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Log the blockchain state
	log.Printf("Current blockchain: %v", n.Chain.Blocks)

	// Encode the blockchain and send it in the response
	if err := json.NewEncoder(w).Encode(n.Chain.Blocks); err != nil {
		log.Printf("Error encoding blockchain: %v", err)
		http.Error(w, "Failed to fetch blockchain", http.StatusInternalServerError)
	}
	log.Println("Finished processing /blockchain request.")
}

// Handle POST /block - Receives a new block from peers and attempts to add it to the blockchain.
func (n *Node) handleReceiveBlock(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling /block request...")

	var block blockchain.Block

	// Decode the incoming block
	body, _ := io.ReadAll(r.Body)
	log.Printf("Received block payload: %s", string(body))
	if err := json.Unmarshal(body, &block); err != nil {
		log.Printf("Error decoding block: %v", err)
		http.Error(w, "Invalid block data", http.StatusBadRequest)
		return
	}

	n.Mu.Lock()
	defer n.Mu.Unlock()

	// Validate the block and add it to the blockchain
	if len(n.Chain.Blocks) == 0 || block.PrevHash == n.Chain.Blocks[len(n.Chain.Blocks)-1].Hash {
		n.Chain.Blocks = append(n.Chain.Blocks, &block)
		log.Printf("Block added successfully: %+v", block)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "Block added successfully")
	} else {
		log.Printf("Invalid block: %+v", block)
		http.Error(w, "Block is invalid", http.StatusBadRequest)
	}
	log.Println("Finished processing /block request.")
}

type RegisterRequest struct {
	Address string `json:"address"`
}

// Handle POST /register - Registers a new peer node.
func (n *Node) handleRegisterNode(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling /register request...")

	body, _ := io.ReadAll(r.Body)
	log.Printf("Raw request body: %s", string(body))

	var req RegisterRequest
	if err := json.Unmarshal(body, &req); err != nil || req.Address == "" {
		log.Printf("Invalid payload received: %v, body: %s", err, string(body))
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	log.Printf("Received registration from: %s", req.Address)

	n.Mu.Lock()
	defer n.Mu.Unlock()
	for _, peer := range n.Peers {
		if peer == req.Address {
			log.Printf("Peer %s is already registered", req.Address)
			http.Error(w, "Peer already registered", http.StatusConflict)
			return
		}
	}

	// Add the new peer
	n.Peers = append(n.Peers, req.Address)
	log.Printf("Peer %s registered successfully", req.Address)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Peer registered: %s", req.Address)
	log.Println("Finished processing /register request.")
}
