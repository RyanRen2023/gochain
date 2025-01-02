# **GoChain: A Simple Blockchain Implementation in Go**

## **Overview**
GoChain is a lightweight blockchain project implemented in Go. It demonstrates the fundamental principles of blockchain technology, including block creation, peer-to-peer networking, proof-of-work, and blockchain synchronization. This project is designed for educational purposes, providing developers with a hands-on learning experience.

---

## **Features**
- **Blockchain Implementation**:
  - Genesis block initialization.
  - Block addition and validation.
  - Proof-of-Work for block mining.
- **Networking**:
  - HTTP API for blockchain operations.
  - Peer registration and synchronization.
- **Modular Code Structure**:
  - Organized files for easy navigation and expansion.

---

## **Prerequisites**
- Go 1.18+ installed on your system.
- Familiarity with Go programming and blockchain concepts.

---

## **Setup Instructions**

### **1. Clone the Repository**
```bash
git clone https://github.com/RyanRen2023/gochain.git
cd gochain
```

### **2. Install Dependencies**
Initialize the Go module and download required dependencies:
```bash
go mod tidy
```

### **3. Run a Node**
Start a node on a specific port:
```bash
go run main.go --port 8080
```

### **4. Run Multiple Nodes**
To simulate a network, run additional nodes on different ports:
```bash
go run main.go --port 8081 --peer http://localhost:8080
```

---

## **Usage**

### **Available Endpoints**

#### **1. Retrieve Blockchain**
- **Endpoint**: `GET /blockchain`
- **Description**: Fetches the current blockchain.
- **Example**:
  ```bash
  curl -X GET http://localhost:8080/blockchain
  ```

#### **2. Add a Block**
- **Endpoint**: `POST /block`
- **Description**: Adds a new block to the blockchain.
- **Payload**:
  ```json
  {
    "Timestamp": 1629271485,
    "Data": "Sample Block Data",
    "PrevHash": "0000000000000000",
    "Hash": "0000000000000001",
    "Nonce": 123
  }
  ```
- **Example**:
  ```bash
  curl -X POST -H "Content-Type: application/json" \
    -d '{"Timestamp":1629271485,"Data":"Sample Block Data","PrevHash":"0000000000000000","Hash":"0000000000000001","Nonce":123}' \
    http://localhost:8080/block
  ```

#### **3. Register a Peer**
- **Endpoint**: `POST /register`
- **Description**: Registers a peer node.
- **Payload**:
  ```json
  {
    "address": "http://localhost:8081"
  }
  ```
- **Example**:
  ```bash
  curl -X POST -H "Content-Type: application/json" \
    -d '{"address":"http://localhost:8081"}' \
    http://localhost:8080/register
  ```

#### **4. List Peers**
- **Endpoint**: `GET /peers`
- **Description**: Lists all registered peer nodes.
- **Example**:
  ```bash
  curl -X GET http://localhost:8080/peers
  ```

---

## **Directory Structure**
The project is organized as follows:

```
gochain/
├── README.md              # Documentation
├── blockchain/            # Blockchain-related logic
│   ├── block.go           # Block structure and methods
│   ├── blockchain.go      # Blockchain structure and methods
│   └── proof.go           # Proof-of-Work logic
├── network/               # Networking-related logic
│   ├── node.go            # Node structure and logic
│   ├── server.go          # HTTP server handlers
├── utils/                 # Utility functions
│   └── utils.go           # Common utilities for the project
├── main.go                # Entry point of the application
├── go.mod                 # Go module configuration
└── server.log             # Log file for server activity (runtime-generated)
```

---

## **Logs**
All server activity is logged in `server.log`. Use the following command to view logs:
```bash
cat server.log
```

---

## **Future Improvements**
- Add consensus algorithms (e.g., Proof-of-Stake).
- Enhance blockchain security with encryption and authentication.
- Build a web-based dashboard for blockchain visualization.
- Implement advanced networking features like gossip protocol.

---

## **Contributing**
Contributions are welcome! If you’d like to contribute:
1. Fork this repository.
2. Create a new branch for your feature or fix.
3. Submit a pull request for review.

---

## **License**
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## **Acknowledgments**
- Inspired by various blockchain tutorials and projects.
- Designed as an educational resource for Go developers.


