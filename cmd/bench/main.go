package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"

	blc "github.com/felipeagger/go-blockchain/blockchain"
)

var difficulty int

func init() {
	runtime.GOMAXPROCS(4)
	difficulty, _ = strconv.Atoi(os.Getenv("DIFFICULTY"))
	if difficulty == 0 {
		difficulty = 5
	}
	fmt.Println("gomaxprocs: 4")
	fmt.Println("difficulty: " + strconv.Itoa(difficulty))
}

func main() {
	Mine(difficulty)
}

func Mine(difficulty int) {
	block, _, _, _ := genTestMine()

	now := time.Now()
	nonce, hash := block.Mine(difficulty)
	curTime := time.Since(now)

	fmt.Printf("nonce: %d, hash: %v", nonce, hash)
	fmt.Printf("\ntempo total: %dms (%.2fs)\n", curTime.Milliseconds(), curTime.Seconds())
}

func genTestMine() (*blc.Block, int, int, string) {
	now := time.Unix(1729776083, 0) //.Format("2006-01-02 15:04:05")
	tx, _ := blc.NewTransaction(nil, "Genesis", "Alice", 100_000_000)
	tx.Timestamp = now
	prevHash := "0000f63d9879d12afa64c2772c185f5b1a2547bf6c1752a80cced4013727007f"
	difficulty := 5

	return &blc.Block{
		Transactions: []blc.Transaction{tx},
		Hash:         "",
		PreviousHash: prevHash,
		Timestamp:    now,
		Nonce:        1,
	}, difficulty, 0, ""
}
