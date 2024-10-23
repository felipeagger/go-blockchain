package main

import (
	"database/sql"
	"flag"
	"fmt"
	blc "github.com/felipeagger/go-blockchain/blockchain"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"runtime"
	"strings"
	"time"
)

var (
	serverPort = blc.GetEnv("PORT", "8765")
	//Difficulty, _ = strconv.Atoi(blc.GetEnv("DIFFICULTY", "5"))
)

func main() {
	createGenesisBlock := flag.Bool("create-genesis-block", false,
		"Cria o bloco genesis da blockchain")
	difficulty := flag.Int("difficulty", 4, "Dificuldade de mineracao")
	nodes := flag.String("node-pool", "", "Nodes do Pool")
	flag.Parse()

	nodePool := strings.Split(*nodes, ",")

	db := blc.InitializeDatabase("./blockchain.db")
	defer db.Close()

	err := blc.ExecMigrations(db)
	if err != nil {
		log.Fatal(err)
	}

	err = blc.Synchronize(db, nodePool, *difficulty)
	if err != nil {
		log.Println(err)
	}

	var blockchain *blc.Blockchain
	if *createGenesisBlock {
		fmt.Println("Creating genesis block...")
		blockchain, err = blc.CreateBlockchain(db, *difficulty)
	} else {
		blockchain, err = blc.LoadBlockchain(db, *difficulty)
	}
	if err != nil {
		log.Fatal(err)
	}

	tests(db, blockchain)

	startServer(blockchain, serverPort)
}

func tests(db *sql.DB, blockchain *blc.Blockchain) {
	var memStats runtime.MemStats

	runtime.ReadMemStats(&memStats)

	fmt.Println("NumGoroutine: ", runtime.NumGoroutine())
	fmt.Printf("MemAlocada: %v bytes | MemTotalAlocada: %v bytes | MemSysUsed: %v bytes | GarbageCollections: %v\n",
		memStats.Alloc, memStats.TotalAlloc, memStats.Sys, memStats.NumGC)

	err := blockchain.AddBlock("Alice", "Bob", 50)
	if err != nil {
		log.Fatal(err)
	}

	runtime.ReadMemStats(&memStats)
	fmt.Println("NumGoroutine: ", runtime.NumGoroutine())
	fmt.Printf("MemAlocada: %v bytes | MemTotalAlocada: %v bytes | MemSysUsed: %v bytes | GarbageCollections: %v\n",
		memStats.Alloc, memStats.TotalAlloc, memStats.Sys, memStats.NumGC)

	err = blockchain.AddBlock("Bob", "John", 20)
	if err != nil {
		log.Fatal(err)
	}

	runtime.ReadMemStats(&memStats)
	fmt.Println("NumGoroutine: ", runtime.NumGoroutine())
	fmt.Printf("MemAlocada: %v bytes | MemTotalAlocada: %v bytes | MemSysUsed: %v bytes | GarbageCollections: %v\n",
		memStats.Alloc, memStats.TotalAlloc, memStats.Sys, memStats.NumGC)

	err = blockchain.AddBlock("John", "Doe", 10)
	if err != nil {
		log.Fatal(err)
	}

	runtime.ReadMemStats(&memStats)
	fmt.Println("NumGoroutine: ", runtime.NumGoroutine())
	fmt.Printf("MemAlocada: %v bytes | MemTotalAlocada: %v bytes | MemSysUsed: %v bytes | GarbageCollections: %v\n",
		memStats.Alloc, memStats.TotalAlloc, memStats.Sys, memStats.NumGC)

	err = blockchain.AddBlock("Doe", "Jane", 5)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(5 * time.Second)

	runtime.ReadMemStats(&memStats)
	fmt.Println("NumGoroutine: ", runtime.NumGoroutine())
	fmt.Printf("MemAlocada: %v bytes | MemTotalAlocada: %v bytes | MemSysUsed: %v bytes | GarbageCollections: %v\n",
		memStats.Alloc, memStats.TotalAlloc, memStats.Sys, memStats.NumGC)

	isValid := blockchain.IsValid()
	fmt.Println(isValid)

	//queryBlock, err := blc.GetBlock(db, block.Hash)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//fmt.Println(queryBlock)

	//blockchain.AddBlock("", "", 0.001)
}
