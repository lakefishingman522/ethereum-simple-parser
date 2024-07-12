package parser

import (
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/lakefishingman522/ethereum-simple-parser/db"
	"github.com/lakefishingman522/ethereum-simple-parser/ethereum"
)

// Parser interface for an Ethereum blockchain parser that queries transactions for subscribed addresses
type Parser interface {
	GetCurrentBlock() uint64
	GetTransactions(address string) []ethereum.Transaction
	SubscribeAddress(address string) bool
	UpdateTransactionsData()
}

// EthereumSimpleParser is a simplified Ethereum blockchain parser
type EthereumSimpleParser struct {
	endpoint string
	store    db.Store
}

// NewEthereumSimpleParser initializes a new instance of EthereumSimpleParser
func NewEthereumSimpleParser(endpoint string, store db.Store) *EthereumSimpleParser {
	return &EthereumSimpleParser{endpoint: endpoint, store: store}
}

// GetCurrentBlock returns the current block number of the Ethereum network
func (parser *EthereumSimpleParser) GetCurrentBlock() uint64 {
	if parser == nil {
		log.Println("[ERROR] parser.go: Parser is nil")
		return 0
	}

	var blockNumberHex string
	if err := ethereum.RPCQuery(parser.endpoint, "eth_blockNumber", nil, &blockNumberHex); err != nil {
		log.Printf("[ERROR] parser.go: GetCurrentBlock RPCQuery: %v", err)
		return 0
	}

	blockNumber, err := ParseHexUint64(blockNumberHex)
	if err != nil {
		log.Printf("[ERROR] parser.go: Failed to parse hex uint64: %v", err)
		return 0
	}

	return blockNumber
}

// UpdateTransactionsData updates the transactions data for all subscribed addresses
func (parser *EthereumSimpleParser) UpdateTransactionsData() {
	if parser == nil {
		return
	}

	subscribers, err := parser.store.GetSubscribers()
	if err != nil {
		return
	}

	var blockNumberHex string
	if err := ethereum.RPCQuery(parser.endpoint, "eth_blockNumber", nil, &blockNumberHex); err != nil {
		return
	}

	currentBlockNumber, err := ParseHexUint64(blockNumberHex)
	if err != nil {
		return
	}

	minBlockNumber := findMinBlockNumber(subscribers)
	if minBlockNumber == math.MaxUint64 {
		return
	}
	updateSubscribersData(parser, subscribers, minBlockNumber, currentBlockNumber)
}

func findMinBlockNumber(subscribers map[string]db.TransactionData) uint64 {
	minBlockNumber := uint64(math.MaxUint64)
	for _, data := range subscribers {
		if data.LastBlockNumber < minBlockNumber {
			minBlockNumber = data.LastBlockNumber
		}
	}
	return minBlockNumber
}

func updateSubscribersData(parser *EthereumSimpleParser, subscribers map[string]db.TransactionData, startBlock, endBlock uint64) {
	for blockNumber := startBlock; blockNumber <= endBlock; blockNumber++ {
		var block ethereum.Block
		err := ethereum.RPCQuery(parser.endpoint, "eth_getBlockByNumber", ParseToAnySlice(fmt.Sprintf("0x%x", blockNumber), true), &block)
		if err != nil {
			return
		}

		for _, transaction := range block.Transactions {
			for address, data := range subscribers {
				if transaction.From == address || transaction.To == address {
					data.Transactions = append(data.Transactions, transaction)
				}
				data.LastBlockNumber = blockNumber
				subscribers[address] = data
			}
		}
	}
}

// GetTransactions retrieves transactions for a specific subscribed address
func (parser *EthereumSimpleParser) GetTransactions(address string) []ethereum.Transaction {
	if parser == nil {
		log.Println("[ERROR] parser.go: Parser is nil")
		return nil
	}

	if address == "" {
		log.Println("[ERROR] parser.go: Address is not defined")
		return nil
	}

	ok, err := parser.store.IsSubscriber(address)
	if err != nil || !ok {
		log.Printf("[ERROR] parser.go: Address %v is not subscribed", address)
		return nil
	}

	currentBlock := parser.GetCurrentBlock()
	if currentBlock == 0 {
		log.Println("[ERROR] parser.go: Failed to get current block number")
		return nil
	}

	data, err := parser.store.GetSubscriber(address)
	if err != nil {
		log.Printf("[ERROR] parser.go: Failed to get subscription data for address %v: %v", address, err)
		return nil
	}

	if data.LastBlockNumber > currentBlock {
		log.Println("[ERROR] parser.go: Subscribed block number is larger than current block number")
		return nil
	}

	return parser.retrieveTransactions(address, data.LastBlockNumber, currentBlock)
}

func (parser *EthereumSimpleParser) retrieveTransactions(address string, startBlock, endBlock uint64) []ethereum.Transaction {
	transactions := []ethereum.Transaction{}
	for blockNumber := startBlock; blockNumber <= endBlock; blockNumber++ {
		var block ethereum.Block
		err := ethereum.RPCQuery(parser.endpoint, "eth_getBlockByNumber", ParseToAnySlice(fmt.Sprintf("0x%x", blockNumber), true), &block)
		if err != nil {
			log.Printf("[ERROR] parser.go: Failed to get block data: %v", err)
			parser.store.SetSubscriber(address, db.TransactionData{LastBlockNumber: blockNumber})
			return transactions
		}
		log.Printf("[INFO] parser.go: Retrieve transactions of block %d", blockNumber)
		for _, transaction := range block.Transactions {
			if transaction.From == address || transaction.To == address {
				transactions = append(transactions, transaction)
			}
		}
	}

	parser.store.SetSubscriber(address, db.TransactionData{LastBlockNumber: endBlock + 1})
	return transactions
}

// SubscribeAddress subscribes to an Ethereum address
func (parser *EthereumSimpleParser) SubscribeAddress(address string) bool {
	if address == "" {
		log.Println("[ERROR] parser.go: Address is not defined")
		return false
	}
	currentBlockNumber := parser.GetCurrentBlock()
	if currentBlockNumber == 0 {
		log.Printf("[ERROR] Failed to get last block number")
		return false
	}
	if err := parser.store.SetSubscriber(address, db.TransactionData{
		LastBlockNumber: currentBlockNumber,
	}); err != nil {
		return false
	}
	return true
}

// ParseHexUint64 converts a hex-encoded string to a uint64
func ParseHexUint64(hexStr string) (uint64, error) {
	return strconv.ParseUint(hexStr[2:], 16, 64)
}

// ParseToAnySlice converts variadic arguments to a slice of empty interfaces
func ParseToAnySlice(params ...interface{}) []interface{} {
	return append([]interface{}{}, params...)
}
