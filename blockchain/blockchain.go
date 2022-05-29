package blockchain

import (
	"time"
)

type Blockchain struct {
	genesisBlock Block
	chain        []Block
	difficulty   int
}

func (b *Blockchain) AddBlock(from, to string, amount float64) {
	blockData := map[string]interface{}{
		"from":   from,
		"to":     to,
		"amount": amount,
	}

	lastBlock := b.chain[len(b.chain)-1]
	newBlock := Block{
		data:         blockData,
		previousHash: lastBlock.hash,
		timestamp:    time.Now(),
	}

	newBlock.Mine(b.difficulty)
	b.chain = append(b.chain, newBlock)
}

func (b Blockchain) IsValid() bool {
	for i := range b.chain[1:] {
		previousBlock := b.chain[i]
		currentBlock := b.chain[i+1]
		if currentBlock.hash != currentBlock.CalculateHash() || currentBlock.previousHash != previousBlock.hash {
			return false
		}
	}
	return true
}

func CreateBlockchain(difficulty int) Blockchain {
	genesisBlock := Block{
		hash:      "0",
		timestamp: time.Now(),
	}

	return Blockchain{
		genesisBlock,
		[]Block{genesisBlock},
		difficulty,
	}
}
