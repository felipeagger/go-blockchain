package blockchain

import (
	"database/sql"
	"strings"
	"time"
)

type Blockchain struct {
	GenesisBlock Block
	Chain        []Block
	difficulty   int
	db           *sql.DB
}

func (b *Blockchain) AddBlock(from, to string, amount float64) error {
	blockData := map[string]interface{}{
		"amount": amount,
		"from":   from,
		"to":     to,
	}

	lastBlock := b.Chain[len(b.Chain)-1]
	newBlock := Block{
		Data:         blockData,
		PreviousHash: lastBlock.Hash,
		Timestamp:    time.Now().UTC(),
	}

	//newBlock.Hash = newBlock.Mine(b.difficulty)
	newBlock.Pow, newBlock.Hash = newBlock.Mine(b.difficulty)
	return b.SaveNewBlock(newBlock)
}

func (b *Blockchain) IsValid() bool {
	for i := range b.Chain[1:] {
		previousBlock := b.Chain[i]
		currentBlock := b.Chain[i+1]
		calculatedHash := currentBlock.CalculateHash()
		prefix := strings.Repeat("0", b.difficulty)

		if currentBlock.Hash != calculatedHash ||
			currentBlock.PreviousHash != previousBlock.Hash ||
			!strings.HasPrefix(currentBlock.Hash, prefix) {
			return false
		}
	}
	return true
}

func (b *Blockchain) SaveNewBlock(block Block) error {
	err := InsertBlock(b.db, block)
	if err != nil {
		return err
	}

	b.Chain = append(b.Chain, block)
	return nil
}

func CreateBlockchain(db *sql.DB, difficulty int) (*Blockchain, error) {
	data := make(map[string]interface{})
	data["amount"] = 1000.0
	data["from"] = "Genesis"
	data["to"] = "Alice"

	genesisBlock := Block{
		Hash:      "0",
		Data:      data,
		Timestamp: time.Now(),
	}

	err := InsertBlock(db, genesisBlock)
	if err != nil {
		return nil, err
	}

	return &Blockchain{
		genesisBlock,
		[]Block{genesisBlock},
		difficulty,
		db,
	}, nil
}
