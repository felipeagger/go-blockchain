package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	blc "github.com/felipeagger/go-blockchain/blockchain"
	"github.com/felipeagger/go-blockchain/wallet"
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

func tests(execTxTests bool, blockchain *blc.Blockchain) {
	if !execTxTests {
		return
	}

	alicePrivKey, alicePubKey, err := wallet.GenerateKeysFromPassword("alice")
	bobPrivKey, bobPubKey, err := wallet.GenerateKeysFromPassword("bob")
	johnPrivKey, johnPubKey, err := wallet.GenerateKeysFromPassword("john")

	aliceAddress := wallet.PublicKeyCompressedToString(alicePubKey)
	bobAddress := wallet.PublicKeyCompressedToString(bobPubKey)
	johnAddress := wallet.PublicKeyCompressedToString(johnPubKey)

	//Alice
	tx1, err := blc.NewTransaction(blockchain,
		aliceAddress,
		bobAddress,
		blc.BtcToSatoshis(0.5))
	if err != nil {
		log.Fatal(err)
	}

	tx1.Sign(alicePrivKey)

	err = blockchain.NewBlock([]blc.Transaction{tx1})
	if err != nil {
		log.Fatal(err)
	}

	//Bob
	tx2, err := blc.NewTransaction(blockchain,
		bobAddress,
		johnAddress,
		blc.BtcToSatoshis(0.2))
	if err != nil {
		log.Fatal(err)
	}

	tx2.Sign(bobPrivKey)

	err = blockchain.NewBlock([]blc.Transaction{tx2})
	if err != nil {
		log.Fatal(err)
	}

	//John
	tx3, err := blc.NewTransaction(blockchain, johnAddress, "doe", blc.BtcToSatoshis(0.1))
	if err != nil {
		log.Fatal(err)
	}

	tx3.Sign(johnPrivKey)

	tx4, err := blc.NewTransaction(blockchain, bobAddress, "jane", blc.BtcToSatoshis(0.1))
	if err != nil {
		log.Fatal(err)
	}

	tx4.Sign(bobPrivKey)

	err = blockchain.NewBlock([]blc.Transaction{tx3, tx4})
	if err != nil {
		log.Fatal(err)
	}

	isValid := blockchain.IsValid()
	fmt.Println(isValid)
}
