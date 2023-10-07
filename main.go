// Name: Khursheed Alam Khan
// Roll: 20i-0496
// Section: SE-A
// Assignment# 1

package main

import (
	//a1 "assignment01bca/assignment01bca"
	"fmt"

	a1 "github.com/i200496-Khursheed/assignment01bca"
)

func main() {
	// Add blocks
	a1.NewBlock("Alice to Bob", 123, "000000000")
	a1.NewBlock("Bob to Charlie", 456, a1.Blockchain[0].Hash)
	a1.NewBlock("Charlie to Delta", 789, a1.Blockchain[1].Hash)
	a1.NewBlock("Delta to Gamma", 101, a1.Blockchain[2].Hash)

	// Display blocks
	fmt.Println("Before changing a block's transaction:")
	fmt.Println("")
	a1.DisplayBlocks()

	// Change the third block's transaction
	a1.ChangeBlock(2, "Charlie to Khursheed")

	fmt.Println("-----------------------------------------")

	// Display blocks after changing a block's transaction
	fmt.Println("\nAfter changing Block 2 transaction from 'Charlie to Delta' to 'Charlie to Khursheed':")
	fmt.Println("")
	a1.DisplayBlocks()

	// Verify blockchain
	isValid := a1.VerifyChain()
	fmt.Printf("\nIs the blockchain valid? %t\n", isValid)
	fmt.Println("")
}
