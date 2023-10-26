// Name: Khursheed Alam Khan
// Roll Number: 20i-0496
// Section: SE-A
// Assignment: 2
package assignment01bca

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
)

// MerkleNode represents a node in the Merkle Tree.
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Hash  string
}

// Block represents a block in the blockchain.
type Block struct {
	Transactions       []Transaction
	PreviousHash       string
	MerkleRoot         string
	PreviousMerkleRoot string
	Hash               string
}

// Transaction represents a transaction with data and nonce.
type Transaction struct {
	Data  string
	Nonce int
}

var Blockchain []Block
var TransactionPool []Transaction
var NumberOfTransactionsPerBlock = 2
var BlockHashMin = "0000000000000000000000000000000000000000000000000000000000000000"
var BlockHashMax = "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"

// CreateHash calculates the hash of a block using the previous hash, Merkle root..
func (b *Block) CreateHash() string {
	data := b.PreviousHash + b.MerkleRoot
	return CalculateHash(data)
}

// NewBlock creates a new block with the given transactions and adds it to the blockchain.
func NewBlock(transactions []Transaction) {
	if len(transactions) == 0 {
		fmt.Println("Error: No transactions provided")
		return
	}

	var previousMerkleRoot string
	if len(Blockchain) > 0 {
		previousMerkleRoot = Blockchain[len(Blockchain)-1].MerkleRoot
	} else {
		previousMerkleRoot = ""
	}

	merkleRoot := BuildMerkleTree(transactions)

	block := Block{
		Transactions:       transactions,
		PreviousHash:       getLastBlockHash(),
		MerkleRoot:         merkleRoot,
		PreviousMerkleRoot: previousMerkleRoot,
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

	previousMerkleRoot := block.MerkleRoot

	block.Transactions[transactionIndex].Data = newTransactionData
	block.Transactions[transactionIndex].Nonce = newNonce

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
	if len(Blockchain) < 2 {
		fmt.Println("Less blocks to verify")
		return false
	}

	for i := 1; i < len(Blockchain); i++ {

		// Verify the Merkle Tree
		previousMerkleRoot := BuildMerkleTree(Blockchain[i-1].Transactions)
		if Blockchain[i].PreviousMerkleRoot != previousMerkleRoot {
			fmt.Println("Block", i, "is compromised. Merkle Tree mismatch.")
			return false
		}

		if Blockchain[i].PreviousHash != Blockchain[i-1].Hash || !isBlockHashInRange(Blockchain[i].Hash) {
			fmt.Println("Block", i, "is compromised. Block hash mismatch.")
			return false
		}
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

// BuildMerkleTree builds the Merkle Tree and returns its root hash.
func BuildMerkleTree(transactions []Transaction) string {
	if len(transactions) == 0 {
		return ""
	}

	var nodes []MerkleNode

	for _, trx := range transactions {
		dataToHash := trx.Data + strconv.Itoa(trx.Nonce)
		hash := CalculateHash(dataToHash)
		nodes = append(nodes, MerkleNode{Hash: hash})
	}

	for len(nodes) > 1 {
		var newLevel []MerkleNode

		for i := 0; i < len(nodes); i += 2 {
			left := &nodes[i]
			right := left
			if i+1 < len(nodes) {
				right = &nodes[i+1]
			}
			hash := CalculateHash(left.Hash + right.Hash)
			node := MerkleNode{Left: left, Right: right, Hash: hash}
			newLevel = append(newLevel, node)
		}

		nodes = newLevel
	}

	return nodes[0].Hash
}
