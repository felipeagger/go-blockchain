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
	Data         map[string]interface{}
	Hash         string
	PreviousHash string
	Timestamp    time.Time
	Pow          int
}

func (b *Block) CalculateHash() string {
	return CalculateHash(b.Timestamp, b.Pow, b.PreviousHash, b.Data)
}

func (b *Block) AsyncCalculateHash(wg *sync.WaitGroup, channel chan string, idx int) {
	defer wg.Done()
	channel <- b.CalculateHash()
}
