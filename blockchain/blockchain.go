package blockchain

import (
        "time"
)


type Block struct {
	data         map[string]interface{}
	hash         string
	previousHash string
	timestamp    time.Time
	pow          int
}

type Blockchain struct {
	genesisBlock Block
	chain        []Block
	difficulty   int
}


func (b *Blockchain) addBlock(from, to string, amount float64) {
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

	newBlock.mine(b.difficulty)
	b.chain = append(b.chain, newBlock)
}

func (b Blockchain) isValid() bool {
	for i := range b.chain[1:] {
			previousBlock := b.chain[i]
			currentBlock := b.chain[i+1]
			if currentBlock.hash != currentBlock.calculateHash() || currentBlock.previousHash != previousBlock.hash {
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