package parser

import (
	"testing"

	"github.com/lakefishingman522/ethereum-simple-parser/db"
)

// Success and failure markers.
const (
	Success = "\u2713"
	Failed  = "\u2717"
)

func TestTxParserService(t *testing.T) {
	var (
		Endpoint = "https://cloudflare-eth.com"
		Address  = "0xc426e53c0da077676a66edf2245e990e9832d4a0"
	)

	memoryDB := db.NewMemoryDB()
	service := NewEthereumSimpleParser(Endpoint, memoryDB)

	t.Run("Block Test", func(t *testing.T) {
		testID := 0
		currentBlockNumber := service.GetCurrentBlock()
		if currentBlockNumber == 0 {
			t.Fatalf("\t%s\tTest %d:\tShould be able to return the non-zero block number ", Failed, testID)
		}
		t.Logf("\t%s\tTest %d:\tLast scanned block number is %d", Success, testID, currentBlockNumber)
	})
	t.Run("Subscribe Test with Address", func(t *testing.T) {
		testID := 1

		subscribedState := service.SubscribeAddress(Address)
		if subscribedState == false {
			t.Fatalf("\t%s\tTest %d:\t Subscription is failed", Failed, testID)
		} else {
			t.Logf("\t%s\tTest %d:\t Subscription is success", Success, testID)
		}
	})
	t.Run("Transaction Test with Address", func(t *testing.T) {
		testID := 2

		transactionData := service.GetTransactions(Address)
		if transactionData == nil {
			t.Fatalf("\t%s\tTest %d:\t Failed to get transaction", Failed, testID)
		} else {
			t.Logf("\t%s\tTest %d:\t Transaction data of address is %v", Success, testID, transactionData)
		}
	})
}
