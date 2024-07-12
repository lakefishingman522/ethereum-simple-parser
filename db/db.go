package db

import (
	"errors"
	"sync"

	"github.com/lakefishingman522/ethereum-simple-parser/ethereum"
)

// TransactionData represents pending transaction data of a specific subscribed address
type TransactionData struct {
	LastBlockNumber uint64                 // Stores the last block number processed for this address
	Transactions    []ethereum.Transaction // List of pending transactions for this address
}

// Create a new transaction data instance.
func NewTransactionData() *TransactionData {
	return &TransactionData{
		LastBlockNumber: 0,
		Transactions:    []ethereum.Transaction{},
	}
}

// Store defines the interface for interacting with a memory database.
type Store interface {
	GetSubscribers() (map[string]TransactionData, error)             // Get all subscribers
	GetSubscriber(address string) (TransactionData, error)           // Get a specific subscriber by address
	SetSubscriber(address string, blockNumber TransactionData) error // Add or update a subscriber
	IsSubscriber(address string) (bool, error)                       // Check if an address is a subscriber
	DeleteSubscriber(address string) error                           // Remove a subscriber
}

// MemoryDB represents an in-memory database.
type MemoryDB struct {
	subscribers map[string]TransactionData // Map from address to transaction data
	lock        sync.RWMutex               // Mutex to handle race conditions during CRUD operations
}

var (
	ErrMemoryDBNotFound = errors.New("memory database is not found")
	ErrAddressNotFound  = errors.New("address is not found in map")
)

// NewMemoryDB initializes a new MemoryDB instance.
func NewMemoryDB() *MemoryDB {
	return &MemoryDB{
		subscribers: make(map[string]TransactionData),
	}
}

// GetSubscribers returns all subscribers in memory database.
func (memory *MemoryDB) GetSubscribers() (map[string]TransactionData, error) {
	if memory == nil {
		return nil, ErrMemoryDBNotFound
	}
	memory.lock.RLock()
	defer memory.lock.RUnlock()

	return memory.subscribers, nil
}

// GetSubscriber returns transaction data for a specific address.
func (memory *MemoryDB) GetSubscriber(address string) (TransactionData, error) {
	if memory == nil {
		return *NewTransactionData(), ErrMemoryDBNotFound
	}
	memory.lock.RLock()
	defer memory.lock.RUnlock()

	value, exists := memory.subscribers[address]

	if exists {
		return value, nil
	} else {
		return *NewTransactionData(), ErrAddressNotFound
	}
}

// DeleteSubscriber removes a subscriber from the memory database.
func (memory *MemoryDB) DeleteSubscriber(address string) error {
	if memory == nil {
		return ErrMemoryDBNotFound
	}
	memory.lock.Lock() // Use write lock for deleting
	defer memory.lock.Unlock()

	delete(memory.subscribers, address)
	return nil
}

// SetSubscriber adds or updates a subscriber in the memory database.
func (memory *MemoryDB) SetSubscriber(address string, transactionData TransactionData) error {
	if memory == nil {
		return ErrMemoryDBNotFound
	}
	memory.lock.Lock() // Use write lock for adding/updating
	defer memory.lock.Unlock()

	memory.subscribers[address] = transactionData
	return nil
}

// IsSubscriber checks if an address is a subscriber.
func (memory *MemoryDB) IsSubscriber(address string) (bool, error) {
	if memory == nil {
		return false, ErrMemoryDBNotFound
	}
	memory.lock.RLock()
	defer memory.lock.RUnlock()

	_, ok := memory.subscribers[address]
	return ok, nil
}
