package blockchain

import (
	"strconv"
	"sync"
	"time"
)

var (
	MaxGoroutines, _ = strconv.Atoi(GetEnv("MAX_GOROUTINES", "1500"))
)

type Block struct {
	Transactions []Transaction `json:"transactions"`
	Hash         string        `json:"hash"`
	PreviousHash string        `json:"previousHash"`
	Timestamp    time.Time     `json:"timestamp"`
	Nonce        int           `json:"nonce"`
}

func (b *Block) CalculateHash() string {
	return CalculateHash(b.Timestamp.Format("2006-01-02 15:04:05"), strconv.Itoa(b.Nonce), b.PreviousHash, b.GetDataInBytes())
}

func (b *Block) GetDataInBytes() []byte {
	return GetTxsInBytes(b.Transactions)
}

func (b *Block) AsyncCalculateHash(wg *sync.WaitGroup, channel chan string, idx int) {
	defer wg.Done()
	channel <- b.CalculateHash()
}
