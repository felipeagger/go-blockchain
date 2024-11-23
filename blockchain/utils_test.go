package blockchain

import (
	"github.com/stretchr/testify/assert"
	"runtime"
	"strconv"
	"testing"
	"time"
)

func genTestData() (string, int, string, []byte, string) {
	data := make(map[string]interface{})
	data["amount"] = 1000.01
	data["from"] = "Genesis"
	data["to"] = "Alice"
	data["bool"] = true
	data["type"] = 2
	now := time.Unix(1729776083, 0).Format("2006-01-02 15:04:05")
	prevHash := "0000f63d9879d12afa64c2772c185f5b1a2547bf6c1752a80cced4013727007f"
	return now, 1, prevHash, GetDataInBytes(data), "bc6292c8921d55908f8a6f2b295a25ff30fee5f166ac7727332f9272e0fe28f1"
}

func TestCalculateHash(t *testing.T) {
	now, nonce, prevHash, data, expectedHash := genTestData()

	hash := CalculateHash(now, strconv.Itoa(nonce), prevHash, data)

	assert.Equal(t, hash, expectedHash)
}

func BenchmarkCalculateHash(b *testing.B) {
	now, nonce, prevHash, data, _ := genTestData()

	// run bench: go test -bench=. -cpuprofile cpu.prof -memprofile mem.prof -race -benchmem -count=3 -v
	for i := 0; i < b.N; i++ {
		nonce = nonce + 1
		_ = CalculateHash(now, strconv.Itoa(nonce), prevHash, data)
	}
}

func TestMemoryLeakCalcHash(t *testing.T) {
	now, nonce, prevHash, data, _ := genTestData()

	var memBefore, memAfter runtime.MemStats

	runtime.GC()
	runtime.ReadMemStats(&memBefore)

	for i := 0; i < 10000; i++ {
		nonce = nonce + 1
		_ = CalculateHash(now, strconv.Itoa(nonce), prevHash, data)
	}

	runtime.GC()
	runtime.ReadMemStats(&memAfter)

	leaks := int((memAfter.Mallocs - memBefore.Mallocs) - (memAfter.Frees - memBefore.Frees))
	if leaks > 0 {
		t.Errorf("Possível vazamento de memória: memória alocada aumentou de %v para %v bytes (leaks: %v)",
			memBefore.Alloc, memAfter.Alloc, leaks)
	}
}
