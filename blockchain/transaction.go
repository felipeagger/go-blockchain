package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/felipeagger/go-blockchain/wallet"
)

type Transaction struct {
	ID        []byte     `json:"id"`
	Timestamp time.Time  `json:"timestamp"`
	Inputs    []TxInput  `json:"inputs"`
	Outputs   []TxOutput `json:"outputs"`
}

type TxInput struct {
	ID []byte `json:"id"`
	//ID will find the Transaction that a specific output is inside of
	OutIdx int `json:"outIdx"`
	//Out will be the index of the specific output we found within a transaction.
	//For example if a transaction has 4 outputs, we can use this "Out" field to specify which output we are looking for
	Signature string `json:"signature"`
	//Digital signature to authenticate the use of the value.
	PubKey string `json:"pubKey"`
	//Sender's public key.
}

type TxOutput struct {
	Value uint64 `json:"value"`
	//Value would be representative of the amount of coins in a transaction
	PubKey string `json:"pubKey"`
	//Receiver wallet...You are indentifiable by your PubKey
	//PubKey in this iteration will be very straightforward, however in an actual application this is a more complex algorithm
}

func NewTransaction(chain *Blockchain, from, to string, amount uint64) (tx Transaction, err error) {
	var inputs []TxInput
	var outputs []TxOutput

	//STEP 1
	acc, validOutputs := chain.FindSpendableOutputs(from, amount)
	if acc < amount {
		return tx, fmt.Errorf("not enough funds!")
	}

	//STEP 3
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			return tx, err
		}

		for _, out := range outs {
			input := TxInput{txID, out, from, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{amount, to})

	//STEP 4: remaining change
	if acc > amount {
		outputs = append(outputs, TxOutput{acc - amount, from})
	}

	tx = Transaction{nil, time.Now(), inputs, outputs}
	tx.SetID()

	return tx, nil
}

func (tx *Transaction) SetID() error {
	var encoded bytes.Buffer
	var hash [32]byte

	encoder := json.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	if err != nil {
		return err
	}

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
	return nil
}

func CoinbaseTx(toAddress, data string) Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", toAddress)
	}
	//Since this is the "first" transaction of the block, it has no previous output to reference.
	//This means that we initialize it with no ID, and it's OutputIndex is -1
	txIn := TxInput{[]byte{}, -1, data, "genesis"}

	txOut := TxOutput{BtcToSatoshis(RewardGenesis), toAddress}

	return Transaction{nil, time.Now(), []TxInput{txIn}, []TxOutput{txOut}}

}

func (in *TxInput) CanUnlock(data string) bool {
	return in.Signature == data
}

func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}

func (tx *Transaction) IsCoinbase() bool {
	//This checks a transaction and will only return true if it is a newly minted "coin"
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].OutIdx == -1
}

func (tx *Transaction) Sign(privKey *ecdsa.PrivateKey) error {
	hash := tx.CalculateHash()

	signature, err := wallet.SignData(privKey, hash)
	if err != nil {
		return err
	}

	// Create hash signature with priv key
	/*r, s, err := ecdsa.Sign(rand.Reader, privKey, hash)
	if err != nil {
		return err
	}

	// Compact signature
	signature := append(r.Bytes(), s.Bytes()...)*/

	// Add signature and pub key at txInput
	for idx, _ := range tx.Inputs {
		tx.Inputs[idx].Signature = hex.EncodeToString(signature)
		tx.Inputs[idx].PubKey = wallet.PublicKeyCompressedToString(&privKey.PublicKey)
	}

	return nil
}

func (tx *Transaction) CheckIsValid() bool {
	hash := tx.CalculateHash()

	for idx, _ := range tx.Inputs {
		pubKey, err := wallet.PublicKeyToECDSA(tx.Inputs[idx].PubKey)
		if err != nil {
			return false
		}

		signature, err := hex.DecodeString(tx.Inputs[idx].Signature)
		if err != nil {
			return false
		}

		isValid := wallet.VerifySignature(pubKey, hash, signature)
		if !isValid {
			return false
		}
	}

	return true
}

func (tx *Transaction) CalculateHash() []byte {
	data := []byte{}

	for _, input := range tx.Inputs {
		data = append(data, input.ID...)
		data = append(data, byte(input.OutIdx))
	}
	for _, output := range tx.Outputs {
		data = append(data, byte(output.Value))
		data = append(data, output.PubKey...)
	}

	firstHash := sha256.Sum256(data)
	return firstHash[:]
}
