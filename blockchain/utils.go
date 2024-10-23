package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func CalculateHash(date time.Time, pow int, prevHash string, data map[string]interface{}) string {
	_data, _ := json.Marshal(data)
	blockData := prevHash + string(_data) + date.Format("2006-01-02 15:04:05") + strconv.Itoa(pow)
	blockHash := sha256.Sum256([]byte(blockData))
	return fmt.Sprintf("%x", blockHash)
}

func reverseBlocks(blocks []Block) []Block {
	for i, j := 0, len(blocks)-1; i < j; i, j = i+1, j-1 {
		blocks[i], blocks[j] = blocks[j], blocks[i]
	}
	return blocks
}

func PrintHashRate(amount int, duration time.Duration) {
	hashRate := float64(amount) / duration.Seconds()
	fmt.Printf("\nDuration: %.9fs; HashRate: %.2f hashes/sec (%.2f kH/s)\n", duration.Seconds(), hashRate,
		hashRate/1000)
}
