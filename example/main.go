package main

import (
	"fmt"
	"log"

	blc "github.com/felipeagger/go-blockchain/blockchain"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db := blc.InitializeDatabase("./blockchain-temp.db")
	defer db.Close()

	err := blc.ExecMigrations(db)
	if err != nil {
		log.Fatal(err)
	}

	// create a new blockchain instance with a mining difficulty of 2
	blockchain, err := blc.CreateBlockchain(db, 2)
	if err != nil {
		log.Fatal(err)
	}

	// record transactions on the blockchain for Alice, Bob, and John
	err = blockchain.AddBlock("Alice", "Bob", 5)
	if err != nil {
		log.Fatal(err)
	}

	err = blockchain.AddBlock("John", "Bob", 2)
	if err != nil {
		log.Fatal(err)
	}

	// check if the blockchain is valid; expecting true
	fmt.Println("\nIsValid: ", blockchain.IsValid())
}
