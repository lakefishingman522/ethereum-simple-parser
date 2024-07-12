package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/lakefishingman522/ethereum-simple-parser/db"
	"github.com/lakefishingman522/ethereum-simple-parser/parser"
)

func main() {
	const endpoint = "https://cloudflare-eth.com"
	memoryDB := db.NewMemoryDB()
	ethParser := parser.NewEthereumSimpleParser(endpoint, memoryDB)

	//Create a ticket to update the transaction data of subscribed address every 5 mins
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop() // Ensure the ticker is stopped to avoid memory leaks
	stopChan := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				{
					ethParser.UpdateTransactionsData()
				}
			case <-stopChan:
				return
			}
		}
	}()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter Command(block, subscribe, transaction, exit): ")
		//Input user command
		command, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("[ERROR] main.go: Failed to read command: %v", err)
			continue
		}

		handleCommand(strings.TrimSpace(command), ethParser, stopChan)
	}
}

// handleCommand processes different user commands.
func handleCommand(input string, ethParser *parser.EthereumSimpleParser, stopChan chan struct{}) {
	args := strings.Fields(input)

	if len(args) < 1 {
		fmt.Println("[ERROR] main.go: No command entered.")
		return
	}

	switch args[0] {
	case "block":
		log.Printf("[INFO] main.go: %v", ethParser.GetCurrentBlock())
	case "subscribe":
		if len(args) < 2 {
			log.Printf("[ERROR] main.go: No address to subscribe")
		} else {
			log.Printf("[INFO] main.go: %v", ethParser.SubscribeAddress(args[1]))
		}
	case "transaction":
		if len(args) < 2 {
			log.Printf("[ERROR] main.go: No address to get transactions")
		} else {
			log.Printf("[INFO] main.go: %v", ethParser.GetTransactions(args[1]))
		}
	case "exit":
		//close the update routine
		close(stopChan)
		os.Exit(0)
	default:
		fmt.Printf("[ERROR] main.go: Invalid action: %v. Please enter one of (block, subscribe, transaction, exit)\n", args[0])
	}
}
