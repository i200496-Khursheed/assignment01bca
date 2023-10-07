// Name: Khursheed Alam Khan
// Roll: 20i-0496
// Section: SE-A
// Assignment# 1

// Inside package assignment01bca/assignment01bca.go
package assignment01bca

import (
	"crypto/sha256" // importing sha256 encryption from go crypto package
	"encoding/hex"  // encoding in hex importing
	"fmt"           // for println purpose
)

// creating a struc heterogenious array storing data sets of Trx, Nonce, hash, and previous hash
type Block struct {
	Transaction  string
	Nonce        int
	PreviousHash string
	Hash         string
}

// an empty array named Blockchain to store multiple Block instances
var Blockchain []Block

// function to create a new block using trx, nonce, and previous hash of previous block.
func NewBlock(transaction string, nonce int, previousHash string) Block {

	// Validate Transaction is Not an empty string
	if len(transaction) == 0 {
		fmt.Println("Error: Invalid transaction format")
		return Block{}
	}

	// Validate Nonce is not in Negative-Integer Range
	if nonce < 0 {
		fmt.Println("Error: Nonce must be a non-negative integer")
		return Block{}
	}

	// block instance for holding parameters for trx, nonce & previousHash of previous block
	block := Block{
		Transaction:  transaction,
		Nonce:        nonce,
		PreviousHash: previousHash,
	}

	// The current block hash is calculated using the CreateHash function
	block.Hash = block.CreateHash()

	// Adds or appends the block created to the block chain
	Blockchain = append(Blockchain, block)

	// returns the new block created
	return block
}

// Function to display all the blocks through iteration of the block chain | Shows index, Trx, nonce, previous & its current hash
func DisplayBlocks() {
	for i, block := range Blockchain {
		fmt.Printf("Block %d:\n", i)
		fmt.Printf("  Transaction: %s\n", block.Transaction)
		fmt.Printf("  Nonce: %d\n", block.Nonce)
		fmt.Printf("  Previous Hash: %s\n", block.PreviousHash)
		fmt.Printf("  Current Hash: %s\n\n", block.Hash)
	}
}

// Hash creation in SHA256 algorithm that returns the hash in Hex encoded string of the Block data structure based on trx, nonce, previous & current hash
func (b Block) CreateHash() string {
	hashInput := fmt.Sprintf("%s%d%s%s", b.Transaction, b.Nonce, b.PreviousHash, b.Hash)
	hashBytes := sha256.Sum256([]byte(hashInput))
	return hex.EncodeToString(hashBytes[:])
}

// Function to change Trx of a block at an index.
func ChangeBlock(index int, transaction string) {

	// Validate Index not negative integer and not greater than blockchain length
	if index < 0 || index >= len(Blockchain) {
		fmt.Println("Error: Invalid block index")
		return
	}

	// Validate Transaction not an Empty String
	if len(transaction) == 0 {
		fmt.Println("Error: Invalid transaction format")
		return
	}

	// updating the block at an index with new transaction | updating its hash and putting it back into the same index
	block := Blockchain[index]
	block.Transaction = transaction
	block.Hash = block.CreateHash()
	Blockchain[index] = block

	// Update hashes of subsequent blocks
	// for i := index + 1; i < len(Blockchain); i++ {
	// 	Blockchain[i].PreviousHash = Blockchain[i-1].Hash
	// 	Blockchain[i].Hash = Blockchain[i].CreateHash()
	// }

	// // Update the hash of the previous block if it's not the first block
	// if index > 0 {
	// 	Blockchain[index-1].Hash = Blockchain[index-1].CreateHash()
	// }
}

// Function to verify integrity of block | Changed previous hash returns 'false' (invalid block) | mathcing previous hash returns 'true' (valid block)
func VerifyChain() bool {
	// Validate Blockchain Length
	if len(Blockchain) < 2 {
		fmt.Println("Error: Insufficient blocks to verify")
		return false
	}

	// checks by iterating through BlockChain to check integrity compromised by checking previous hash of current block with its previous block
	// compromised integrity returns 'false' and vice versa
	for i := 1; i < len(Blockchain); i++ {
		if Blockchain[i].PreviousHash != Blockchain[i-1].Hash {
			fmt.Println("Error: Blockchain integrity compromised")
			return false
		}
	}
	return true
}

// A general purpose SHA256 hash calculating algorithm to return hash in the form of hex encoded string
// func CalculateHash(stringToHash string) string {
// 	// Validate String Length
// 	if len(stringToHash) == 0 {
// 		fmt.Println("Error: Invalid string to hash")
// 		return ""
// 	}

// 	hashBytes := sha256.Sum256([]byte(stringToHash))
// 	return hex.EncodeToString(hashBytes[:])
// }
