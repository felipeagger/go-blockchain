package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

const MAX_GOROUTINES = 1500

type Block struct {
	data         map[string]interface{}
	hash         string
	previousHash string
	timestamp    time.Time
	pow          int
}

func (b Block) CalculateHash() string {
	data, _ := json.Marshal(b.data)
	blockData := b.previousHash + string(data) + b.timestamp.String() + strconv.Itoa(b.pow)
	blockHash := sha256.Sum256([]byte(blockData))
	return fmt.Sprintf("%x", blockHash)
}

func (b Block) AsyncCalculateHash(wg *sync.WaitGroup, channel chan string, idx int) {
	defer wg.Done()
	data, _ := json.Marshal(b.data)
	blockData := b.previousHash + string(data) + b.timestamp.String() + strconv.Itoa(b.pow)
	blockHash := sha256.Sum256([]byte(blockData))
	channel <- fmt.Sprintf("%x", blockHash)
}

/*func (b *Block) SequentialMine(difficulty int) {
	started := time.Now()
	qtdHashes := 0
	for !strings.HasPrefix(b.hash, strings.Repeat("0", difficulty)) {
		b.pow++
		b.hash = b.CalculateHash()
		qtdHashes = qtdHashes + 1
	}
	fmt.Printf("\nMine Duration: %v - QtdHashes: %v", time.Now().Sub(started), qtdHashes)
}*/

func (b *Block) Mine(difficulty int) {
	started := time.Now()
	qtdHashes := 0
	newHash := b.hash
	var wg sync.WaitGroup

	for !strings.HasPrefix(newHash, strings.Repeat("0", difficulty)) {

		wg.Add(MAX_GOROUTINES)
		channel := make(chan string, MAX_GOROUTINES)
		for i := 1; i <= MAX_GOROUTINES; i++ {
			b.pow++
			qtdHashes = qtdHashes + 1
			go b.AsyncCalculateHash(&wg, channel, i)
		}

		go func() {
			wg.Wait()
			close(channel)
		}()

		for hash := range channel {
			if strings.HasPrefix(hash, strings.Repeat("0", difficulty)) {
				newHash = hash
				break
			}
		}

	}

	b.hash = newHash
	fmt.Printf("\nMine Duration: %v - QtdHashes: %v", time.Now().Sub(started), qtdHashes)
}
