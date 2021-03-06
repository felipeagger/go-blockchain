package main

import (
	"fmt"

	blc "github.com/felipeagger/go-blockchain/blockchain"
)

func main() {
	// create a new blockchain instance with a mining difficulty of 2
	blockchain := blc.CreateBlockchain(6)

	// record transactions on the blockchain for Alice, Bob, and John
	blockchain.AddBlock("Alice", "Bob", 5)
	blockchain.AddBlock("John", "Bob", 2)

	// check if the blockchain is valid; expecting true
	fmt.Println()
	fmt.Println(blockchain.IsValid())
}
