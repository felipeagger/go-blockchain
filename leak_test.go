package go_blockchain

import (
	"fmt"
	blc "github.com/felipeagger/go-blockchain/blockchain"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
	"os"
	"runtime"
	"testing"
)

func doBlockchainOperations() (bool, int, error) {
	dbFile := "./blockchain-tests.db"
	os.Remove(dbFile)

	db := blc.InitializeDatabase(dbFile)
	defer db.Close()

	err := blc.ExecMigrations(db)
	if err != nil {
		return false, 0, err
	}

	blockchain, err := blc.CreateBlockchain(db, 2)
	if err != nil {
		return false, 0, err
	}

	err = blockchain.AddBlock("Alice", "Bob", 5)
	if err != nil {
		return false, len(blockchain.Chain), err
	}

	err = blockchain.AddBlock("John", "Bob", 2)
	if err != nil {
		return false, len(blockchain.Chain), err
	}

	err = blockchain.AddBlock("Bob", "Doe", 1)
	if err != nil {
		return false, len(blockchain.Chain), err
	}

	isValid := blockchain.IsValid()
	return isValid, len(blockchain.Chain), err
}

func TestBlockchain(t *testing.T) {
	isValid, size, err := doBlockchainOperations()
	assert.Nil(t, err)
	assert.Equal(t, isValid, true)
	assert.Equal(t, size, 4)

	goleak.VerifyNone(t)
}

func TestMemoryLeak(t *testing.T) {
	runtime.GC()
	var memBefore runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	_, _, _ = doBlockchainOperations()

	runtime.GC()
	var memAfter runtime.MemStats
	runtime.ReadMemStats(&memAfter)

	fmt.Println("Antes de executar a função:")
	fmt.Printf("Memória alocada: %v bytes | Total alocada: %v bytes | Sistema: %v bytes | GC: %v\n",
		memBefore.Alloc, memBefore.TotalAlloc, memBefore.Sys, memBefore.NumGC)
	fmt.Println("Depois de executar a função:")
	fmt.Printf("Memória alocada: %v bytes | Total alocada: %v bytes | Sistema: %v bytes | GC: %v\n",
		memAfter.Alloc, memAfter.TotalAlloc, memAfter.Sys, memAfter.NumGC)

	margin := memBefore.Alloc / 20 // 5%

	if memAfter.Alloc > memBefore.Alloc+margin {
		t.Errorf("Possível vazamento de memória: memória alocada aumentou de %v para %v bytes (excedendo a margem de 5%%)", memBefore.Alloc, memAfter.Alloc)
	} else {
		t.Logf("Aumento de memória dentro da margem de 5%%")
	}
}

func TestMain(m *testing.M) {
	// Verifica se todas as goroutines foram finalizadas ao terminar os testes
	goleak.VerifyTestMain(m)
	os.Exit(m.Run())
}
