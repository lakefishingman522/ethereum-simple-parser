# Ethereum Simple Parser

Ethereum Simple Parser is a Go application that connects to the Ethereum network via a specified endpoint and provides basic functionalities to interact with blocks, subscribe to addresses, retrieve transaction information, and more.

## Table of Contents
- [Installation](#installation)
- [Usage](#usage)
- [Commands](#commands)
- [Architecture](#architecture)

## Installation

1. **Clone the repository**
    ```sh
    git clone https://github.com/lakefishingman522/ethereum-simple-parser.git
    ```
   
2. **Navigate to the project directory**
    ```sh
    cd ethereum-simple-parser
    ```

3. **Install dependencies**
    ```sh
    go mod tidy
    ```

## Usage

To run the application, use the following command:

```sh
go run main.go
```

The application will start and you'll be prompted to enter commands interactively.

## Commands

Here's a list of commands you can use with the application:

- **block**: Retrieve the current block information.
  
  ```sh
  Enter Command(block, subscribe, transaction, exit): block
  ```

- **subscribe [address]**: Subscribe to a specific Ethereum address to monitor its transactions.
  
  ```sh
  Enter Command(block, subscribe, transaction, exit): subscribe 0xYourEthereumAddress
  ```

- **transaction [address]**: Retrieve transaction information for a specific Ethereum address.
  
  ```sh
  Enter Command(block, subscribe, transaction, exit): transaction 0xYourEthereumAddress
  ```

- **exit**: Exit the application.
  
  ```sh
  Enter Command(block, subscribe, transaction, exit): exit
  ```

## Architecture

The application has the following major components:
- **main.go**: Entry point of the application, which manages the command-line interface and user input.
- **db**: Package that includes an in-memory database to store Ethereum data.
- **parser**: Package that handles the interaction with the Ethereum network, updating and retrieving blockchain data.
- **ethereum**: Package that provides method to send JSON-RPC request.

### Directory Structure

```
ethereum-simple-parser/
├── db/
│   └── db.go
├── parser/
│   └── parser.go
├── ethereum/
│   └── ethereum.go
├── main.go
├── go.mod
├── go.sum
└── README.md
```