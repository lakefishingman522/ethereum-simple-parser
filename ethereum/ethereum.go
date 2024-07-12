package ethereum

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// JSON-RPC response structure
type RPCResponse struct {
	Result json.RawMessage `json:"result"`
	Error  interface{}     `json:"error"`
	ID     int             `json:"id"`
}

// Simplified Ethereum block structure
type Block struct {
	Hash         string        `json:"hash"`
	Transactions []Transaction `json:"transactions"`
}

// Simplified Ethereum transaction structure
type Transaction struct {
	Hash        string `json:"hash"`
	BlockNumber string `json:"blockNumber"`
	From        string `json:"from"`
	To          string `json:"to"`
	Value       string `json:"value"`
}

// Send RPC request to get data of Ethereum network
func RPCQuery(endpoint string, method string, params []interface{}, result interface{}) error {
	var response RPCResponse

	// Make request body string using params
	requestBody := fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"method": "%s",
		"params": %s,
		"id": 1
	}`, method, toJSON(params))

	//Send HTTP Post Request to Ethereum network
	resp, err := http.Post(endpoint, "application/json", strings.NewReader(requestBody))

	// Check for errors in sending request
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//Check for the status code
	if resp.StatusCode != http.StatusOK {
		_, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil
		}
		return fmt.Errorf("received non-ok status code: %d", resp.StatusCode)
	}

	// Check for errors in decode of response body
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return err
	}

	// Check for errors in response
	if response.Error != nil {
		return errors.New("json-rpc response error")
	}

	// Parse the result
	err = json.Unmarshal(response.Result, &result)

	// Check for erros in parsing
	if err != nil {
		return err
	}

	return nil
}

// toJSON converts parameters to JSON string.
func toJSON(params []interface{}) string {
	if len(params) == 0 {
		return "[]"
	}

	var builder strings.Builder
	builder.WriteByte('[')
	for i, param := range params {
		jsonParam, _ := json.Marshal(param)
		builder.Write(jsonParam)
		if i < len(params)-1 {
			builder.WriteByte(',')
		}
	}
	builder.WriteByte(']')
	return builder.String()
}
