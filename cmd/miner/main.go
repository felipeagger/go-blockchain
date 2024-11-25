package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	blc "github.com/felipeagger/go-blockchain/blockchain"
	_ "github.com/mattn/go-sqlite3"
)

var (
	serverPort    = blc.GetEnv("PORT", "8765")
	Difficulty, _ = strconv.Atoi(blc.GetEnv("DIFFICULTY", "5"))
	blockchain    *blc.Blockchain
)

func main() {
	createGenesisBlock := flag.Bool("create-genesis-block", false,
		"Cria o bloco genesis da blockchain")
	difficulty := flag.Int("difficulty", 4, "Dificuldade de mineracao")
	nodes := flag.String("node-pool", "", "Nodes do Pool")
	txTests := flag.Bool("tests-txs", false, "Transacoes de teste")
	flag.Parse()

	nodePool := strings.Split(*nodes, ",")

	db := blc.InitializeDatabase("./blockchain.db")
	defer db.Close()

	err := blc.ExecMigrations(db)
	if err != nil {
		log.Fatal(err)
	}

	if *createGenesisBlock {
		fmt.Println("Creating genesis block...")
		blockchain, err = blc.CreateBlockchain(db, *difficulty)
	} else {
		err = blc.Synchronize(db, nodePool, *difficulty)
		if err != nil {
			log.Println(err)
		}

		blockchain, err = blc.LoadBlockchain(db, *difficulty)
	}
	if err != nil {
		log.Fatal(err)
	}

	go api()

	tests(*txTests, blockchain)

	isValid := blockchain.IsValid()
	fmt.Println(isValid)

	startServer(blockchain, serverPort)
}
