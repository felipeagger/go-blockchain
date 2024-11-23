package blockchain

import (
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func genTestMine() (*Block, int, int, string) {
	now := time.Unix(1729776083, 0) //.Format("2006-01-02 15:04:05")
	tx, _ := NewTransaction(nil, "Genesis", "Alice", 100_000_000)
	tx.Timestamp = now
	prevHash := "0000f63d9879d12afa64c2772c185f5b1a2547bf6c1752a80cced4013727007f"
	difficulty := 4
	expectedNonce := 43683
	expectedHash := "00003fe697fe9cac490a587ebc51d01be0a0c8147b2cd6ffbcc74f1de3f49d65"
	return &Block{
		Transactions: []Transaction{tx},
		Hash:         "",
		PreviousHash: prevHash,
		Timestamp:    now,
		Nonce:        1,
	}, difficulty, expectedNonce, expectedHash
}

func TestBlock_Mine(t *testing.T) {
	block, difficulty, expectedNonce, expectedHash := genTestMine()

	nonce, hash := block.Mine(difficulty)

	assert.Equal(t, nonce, expectedNonce)
	assert.Equal(t, hash, expectedHash)
}

func TestMineMemoryLeak(t *testing.T) {
	block, difficulty, _, _ := genTestMine()

	runtime.GC()
	var memBefore runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	for i := 0; i < 5; i++ {
		_, _ = block.Mine(difficulty)
	}

	time.Sleep(1 * time.Second)

	runtime.GC()
	var memAfter runtime.MemStats
	runtime.ReadMemStats(&memAfter)

	leaks := int((memAfter.Mallocs - memBefore.Mallocs) - (memAfter.Frees - memBefore.Frees))
	if leaks > 0 {
		t.Errorf("Possível vazamento de memória: memória alocada aumentou de %v para %v bytes (leaks: %v)",
			memBefore.Alloc, memAfter.Alloc, leaks)
	}
}

func TestMain(m *testing.M) {
	// Verifica se todas as goroutines foram finalizadas ao terminar os testes
	goleak.VerifyTestMain(m)
	os.Exit(m.Run())
}
