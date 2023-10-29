// Name: Khursheed Alam Khan
// Roll Number: 20i-0496
// Section: SE-A
// Assignment: 2
package assignment01bca

// Importing Packages
import (
	"crypto/sha256" // provides SHA256 hash encoding function
	"encoding/hex"  // encodes & decodes hexadecimal strings
	"fmt"           // printing output on console
	"strconv"       // converts string to integer and vice versa
)

// MerkleNode represents a node in the Merkle Tree. Leaf nodes show transaction, non-leaf nodes have hashes
type MerkleNode struct {
	Left  *MerkleNode // left child of current merkle node
	Right *MerkleNode // right child of current merkle node
	Hash  string      // concatenates left + right child hash
}

// Block represents a block in the blockchain.
type Block struct {
	Transactions       []Transaction // stores transaction (data + nonce) for multiple transaction in a single block
	PreviousHash       string        // previous hash of previous block
	MerkleRoot         string        // hash of merkle root of current block
	PreviousMerkleRoot string        // previous merkle root hash to check for tampering of transaction
	Hash               string        // Current block hash (current block merkle root hash + previous block hash)
}

// Transaction represents a transaction with data and nonce.
type Transaction struct {
	Data  string // transaction data
	Nonce int    // nonce for the transaction data
}

var Blockchain []Block                                                                // block slice
var TransactionPool []Transaction                                                     // Transaction Pool to temporalily store Transactions
var NumberOfTransactionsPerBlock = 2                                                  // total number of transactions to produce a new block
var BlockHashMin = "0000000000000000000000000000000000000000000000000000000000000000" // lowest nonce target required
var BlockHashMax = "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff" // highest nonce traget required

// CreateHash calculates the hash of a block using the previous hash, Merkle root..
func (b *Block) CreateHash() string {
	data := b.PreviousHash + b.MerkleRoot // stores previous block hash and merkle root hash in data
	return CalculateHash(data)            // calculates new hash of current block
}

// NewBlock creates a new block with the given transactions and adds it to the blockchain.
func NewBlock(transactions []Transaction) {
	// validation step to check if transaction is empty | Display Error if yes means length of transaction slice is 0
	if len(transactions) == 0 {
		fmt.Println("Error: No transactions provided")
		return
	}

	merkleRoot := BuildMerkleTree(transactions)

	block := Block{
		Transactions:       transactions,
		PreviousHash:       getLastBlockHash(),
		MerkleRoot:         merkleRoot,
		PreviousMerkleRoot: merkleRoot, // Store the current block's Merkle Root in both MerkleRoot and PreviousMerkleRoot
	}

	block.Hash = block.CreateHash()
	Blockchain = append(Blockchain, block)
}

// ChangeBlock allows changing a transaction within a specified block, recalculates Merkle root, and checks for tampering.
func ChangeBlock(blockIndex int, transactionIndex int, newTransactionData string, newNonce int) {
	if blockIndex < 0 || blockIndex >= len(Blockchain) {
		fmt.Println("Error: Invalid block index")
		return
	}

	block := &Blockchain[blockIndex]

	if transactionIndex < 0 || transactionIndex >= len(block.Transactions) {
		fmt.Println("Error: Invalid transaction index")
		return
	}

	// Store the current Merkle Root before changing the transaction
	previousMerkleRoot := block.MerkleRoot

	// Change the transaction data and nonce
	block.Transactions[transactionIndex].Data = newTransactionData
	block.Transactions[transactionIndex].Nonce = newNonce

	// Recalculate the Merkle Root for the current block
	block.MerkleRoot = BuildMerkleTree(block.Transactions)
	block.Hash = block.CreateHash()

	if previousMerkleRoot != block.MerkleRoot {
		fmt.Println("Merkle Tree in Block", blockIndex, "has been tampered with.")
	}
}

// AddTransactionToPool adds a new transaction to the transaction pool and creates a new block when the pool is full.
func AddTransactionToPool(transactionData string, nonce int) {
	transaction := Transaction{
		Data:  transactionData,
		Nonce: nonce,
	}
	TransactionPool = append(TransactionPool, transaction)

	if len(TransactionPool) >= NumberOfTransactionsPerBlock {
		NewBlock(TransactionPool)
		TransactionPool = nil
	}
}

// DisplayBlocks prints the details of all blocks in the blockchain.
func DisplayBlocks() {
	for i, block := range Blockchain {
		fmt.Printf("Block %d:\n", i)
		fmt.Printf("  Previous Hash: %s\n", block.PreviousHash)
		fmt.Printf("  Merkle Root: %s\n", block.MerkleRoot)
		fmt.Printf("  Current Hash: %s\n\n", block.Hash)
	}
}

// DisplayTransactionsInBlock prints the details of transactions in a specified block.
func DisplayTransactionsInBlock(index int) {
	if index < 0 || index >= len(Blockchain) {
		fmt.Println("Error: Invalid block index")
		return
	}

	block := Blockchain[index]
	fmt.Printf("Transactions in Block %d:\n", index)
	for i, transaction := range block.Transactions {
		fmt.Printf("Transaction %d:\n", i)
		fmt.Printf("  Data: %s\n", transaction.Data)
		fmt.Printf("  Nonce: %d\n", transaction.Nonce)
	}
}

// VerifyChain checks the integrity of the blockchain by comparing block hashes and Merkle Tree hashes.
func VerifyChain() bool {
	if len(Blockchain) < 1 {
		fmt.Println("No blocks to verify")
		return false
	}

	// Start from index 0
	for i := 0; i < len(Blockchain); i++ {
		// For the first block, check if PreviousMerkleRoot Hash is equal to its current MerkleRoot Hash?
		if i == 0 {
			if Blockchain[i].PreviousMerkleRoot != Blockchain[i].MerkleRoot {
				fmt.Println("Block 0 is compromised. Merkle Root Hash Mismatch!")
				return false
			}
		} else {
			// Verify the Merkle Tree
			if Blockchain[i].PreviousMerkleRoot != Blockchain[i].MerkleRoot {
				fmt.Println("Block", i, "is compromised. Merkle Tree mismatch.")
				return false
			}

			if Blockchain[i].PreviousHash != Blockchain[i-1].Hash || !isBlockHashInRange(Blockchain[i].Hash) {
				fmt.Println("Block", i, "is compromised. Block hash mismatch.")
				return false
			}
		}
	}

	// Verify the Merkle Tree of the last block
	lastBlockIndex := len(Blockchain) - 1
	previousMerkleRoot := BuildMerkleTree(Blockchain[lastBlockIndex].Transactions)
	if Blockchain[lastBlockIndex].PreviousMerkleRoot != previousMerkleRoot {
		fmt.Println("Last block is compromised. Merkle Tree mismatch.")
		return false
	} else {
		fmt.Println("All TRX Valid.")
	}

	return true
}

// CalculateHash calculates the SHA-256 hash of a given string.
func CalculateHash(stringToHash string) string {
	hashBytes := sha256.Sum256([]byte(stringToHash))
	return hex.EncodeToString(hashBytes[:])
}

// getLastBlockHash retrieves the hash of the last block in the blockchain.
func getLastBlockHash() string {
	if len(Blockchain) == 0 {
		return "0"
	}
	return Blockchain[len(Blockchain)-1].Hash
}

// isBlockHashInRange checks if a block's hash is within a specified range.
func isBlockHashInRange(hash string) bool {
	return hash >= BlockHashMin && hash <= BlockHashMax
}

// BuildMerkleTree builds the Merkle Tree and returns its root hash. Takes slice of transaction object
func BuildMerkleTree(transactions []Transaction) string {
	if len(transactions) == 0 {
		return ""
	}

	var nodes []MerkleNode // slice of merkle nodes to hold leaf nodes of merkle tree

	for _, trx := range transactions { // iterates over all transactions in each leaf node and assigns the current transaction to the trx vairable
		dataToHash := trx.Data + strconv.Itoa(trx.Nonce) // put trx data and its nonce in dataToHash variable
		hash := CalculateHash(dataToHash)                // dataToHash is hashed and put into hash variable of the concatenated dataToHash
		nodes = append(nodes, MerkleNode{Hash: hash})    // a node is created into the merkle tree with the hash calculated
	}

	for len(nodes) > 1 { // iteration loop to run until more than 1 node i.e. until root node found
		var newLevel []MerkleNode // newLelvel MerkleNode slice to hold new level of merkle tree nodes

		for i := 0; i < len(nodes); i += 2 { // increments of 2 to start i with left node and later check i+1 for right node to form pairs.
			left := &nodes[i] // left is pointed towards the ith node
			right := left     // sets right node to left node as well in-case of odd number of nodes

			if i+1 < len(nodes) {
				right = &nodes[i+1] // points right towards the next node if available within the same node slice
			}
			hash := CalculateHash(left.Hash + right.Hash)            // internal node hash caluclated by concatenating left and right nodes hash
			node := MerkleNode{Left: left, Right: right, Hash: hash} // new node created with the left and right node set and concatenated hash calculated
			newLevel = append(newLevel, node)
		}

		nodes = newLevel // nodes updated to contain newLevel nodes created
	}

	return nodes[0].Hash // the entire Merkle tree returned by its root hash.
}
