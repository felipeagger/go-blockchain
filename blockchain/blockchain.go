package blockchain

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"strings"
	"time"
)

const (
	Reward   = 3
	RewardTo = "028348e9b430ae9986a1d1f88abe6e0196bc0d3c332c2c4bf5f2852d6b742b87bd" //"alice"
)

type Blockchain struct {
	GenesisBlock Block
	Chain        []Block
	difficulty   int
	db           *sql.DB
}

func CreateBlockchain(db *sql.DB, difficulty int) (*Blockchain, error) {
	tx := Transaction{
		ID:        []byte("genesis"),
		Timestamp: time.Now(),
		Outputs:   []TxOutput{{PubKey: RewardTo, Value: 10_000_000_000}},
	}

	genesisBlock := Block{
		Hash:         "0",
		Transactions: []Transaction{tx},
		Timestamp:    time.Now(),
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

func (b *Blockchain) NewBlock(txs []Transaction) error {
	var transactions []Transaction
	txCoinbase := CoinbaseTx(RewardTo, "")

	transactions = append(transactions, txCoinbase)

	for _, tx := range txs {
		transactions = append(transactions, tx)
	}

	return b.addBlock(transactions)
}

func (b *Blockchain) addBlock(txs []Transaction) error {
	if len(b.Chain) == 0 {
		return errors.New("no blockchain detected: create using --create-genesis-block")
	}

	lastBlock := b.Chain[len(b.Chain)-1]
	newBlock := Block{
		Transactions: txs,
		PreviousHash: lastBlock.Hash,
		Timestamp:    time.Now().UTC(),
	}

	newBlock.Nonce, newBlock.Hash = newBlock.Mine(b.difficulty)
	return b.SaveNewBlock(newBlock)
}

func (b *Blockchain) IsValid() bool {
	if len(b.Chain) == 0 {
		return true
	}

	for i := range b.Chain[1:] {
		previousBlock := b.Chain[i]
		currentBlock := b.Chain[i+1]
		calculatedHash := currentBlock.CalculateHash()
		prefix := strings.Repeat("0", b.difficulty)

		//validate block
		if currentBlock.Hash != calculatedHash ||
			currentBlock.PreviousHash != previousBlock.Hash ||
			!strings.HasPrefix(currentBlock.Hash, prefix) {
			return false
		}

		//validate transactions
		for _, tx := range currentBlock.Transactions {
			if !tx.IsCoinbase() {
				if !tx.CheckIsValid() {
					return false
				}
			}
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

func (b *Blockchain) FindUnspentTransactions(address string) []Transaction {
	var unspentTxs []Transaction

	spentTXNs := make(map[string][]int)

	//for _, block := range b.Chain {
	for i := len(b.Chain) - 1; i >= 0; i-- {
		block := b.Chain[i]
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

			//Check Outputs
		outputs:
			for outIndex, output := range tx.Outputs {
				// Se essa saída já foi gasta, ignore.
				if spentIndexes, exists := spentTXNs[txID]; exists {
					for _, spentIndex := range spentIndexes {
						if spentIndex == outIndex {
							continue outputs
						}
					}
				}

				// Se a saída pertence ao endereço, adicione à lista de não gastos.
				if output.CanBeUnlocked(address) {
					unspentTxs = append(unspentTxs, tx)
				}
			}

			//Check Inputs
			if tx.IsCoinbase() == false {
				for _, in := range tx.Inputs {
					if in.CanUnlock(address) {
						inTxID := hex.EncodeToString(in.ID)
						spentTXNs[inTxID] = append(spentTXNs[inTxID], in.OutIdx)
					}
				}
			}

			if len(block.PreviousHash) == 0 {
				break
			}
		}
	}

	return unspentTxs
}

func (chain *Blockchain) FindUTXO(address string) []TxOutput {
	var UTXOs []TxOutput
	unspentTransactions := chain.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

func (b *Blockchain) FindSpendableOutputs(address string, amount uint64) (uint64, map[string][]int) {
	unspentOuts := make(map[string][]int)
	var accumulated uint64

	unspentTxs := b.FindUnspentTransactions(address)

	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
			if out.CanBeUnlocked(address) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				if accumulated >= amount {
					break
				}
			}
		}
	}

	return accumulated, unspentOuts
}

func (b *Blockchain) GetAddressBalance(address string) uint64 {
	var accumulated uint64

	unspentTxs := b.FindUnspentTransactions(address)

	for _, tx := range unspentTxs {

		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				accumulated += out.Value
			}
		}
	}

	return accumulated
}
