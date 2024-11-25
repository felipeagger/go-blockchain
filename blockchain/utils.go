package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetTxsInBytes(txs []Transaction) []byte {
	sizeOfData := len(txs)
	if sizeOfData == 0 {
		return []byte{}
	}

	sort.Slice(txs, func(i, j int) bool {
		return txs[i].Timestamp.Before(txs[j].Timestamp)
	})

	data, _ := json.Marshal(txs)
	return data
}

func GetDataInBytes(data map[string]interface{}) []byte {
	sizeOfData := len(data)
	if sizeOfData == 0 {
		return []byte{}
	}

	keys := make([]string, 0, sizeOfData)
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var builder strings.Builder
	for _, key := range keys {
		value := data[key]

		builder.WriteString(key)
		builder.WriteString(":")
		builder.WriteString(fmt.Sprintf("%v", value))
		builder.WriteString(",")
	}

	return []byte(builder.String())
}

func CalculateHash(date string, nonce string, prevHash string, data []byte) string {
	var builder strings.Builder

	//blockData := prevHash + string(data) + date.Format("2006-01-02 15:04:05") + strconv.Itoa(nonce)
	builder.WriteString(prevHash)
	builder.Write(data)
	builder.WriteString(date)
	builder.WriteString(nonce)

	blockHash := sha256.Sum256([]byte(builder.String()))
	return hex.EncodeToString(blockHash[:]) //string(blockHash[:]) //bytesToString(blockHash[:])
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

func BtcToSatoshis(btc float64) uint64 {
	const satoshisPerBTC = 100_000_000
	return uint64(btc * satoshisPerBTC)
}

func SatoshisToBtc(sats uint64) float64 {
	const satoshisPerBTC = 100_000_000
	return float64(sats / satoshisPerBTC)
}
