package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/i200496-Khursheed/assignment01bca/assignment01bca"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\nMenu:")
		fmt.Println("1. Add a new transaction")
		fmt.Println("2. Display blocks")
		fmt.Println("3. Display transactions in a block")
		fmt.Println("4. Verify blockchain")
		fmt.Println("5. Change a transaction in a block")
		fmt.Println("6. Exit")
		fmt.Print("Enter your choice: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			fmt.Print("Enter transaction: ")
			scanner.Scan()
			transaction := scanner.Text()
			fmt.Print("Enter nonce: ")
			var nonce int
			fmt.Scanln(&nonce)

			assignment01bca.AddTransactionToPool(transaction, nonce)

		case 2:
			fmt.Println("Displaying blocks:")
			assignment01bca.DisplayBlocks()
		case 3:
			fmt.Print("Enter block index to display transactions: ")
			var index int
			fmt.Scanln(&index)
			assignment01bca.DisplayTransactionsInBlock(index)

		case 4:
			isValid := assignment01bca.VerifyChain()
			fmt.Printf("Is the blockchain valid? %t\n", isValid)

		case 5:
			fmt.Print("Enter block index to change transaction: ")
			var blockIndex int
			fmt.Scanln(&blockIndex)

			fmt.Print("Enter transaction index to change: ")
			var transactionIndex int
			fmt.Scanln(&transactionIndex)

			fmt.Print("Enter new transaction data: ")
			scanner.Scan()
			newTransactionData := scanner.Text()

			fmt.Print("Enter new nonce: ")
			var newNonce int
			fmt.Scanln(&newNonce)

			assignment01bca.ChangeBlock(blockIndex, transactionIndex, newTransactionData, newNonce)

		case 6:
			fmt.Println("Exiting the program.")
			return
		default:
			fmt.Println("Invalid choice. Please select a valid option.")
		}
	}
}
