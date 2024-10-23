package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

const Reward = 100

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

type TxOutput struct {
	Value int
	//Value would be representative of the amount of coins in a transaction
	PubKey string
	//The Pubkey is needed to "unlock" any coins within an Output. This indicated that YOU are the one that sent it.
	//You are indentifiable by your PubKey
	//PubKey in this iteration will be very straightforward, however in an actual application this is a more complex algorithm
}

//TxInput is representative of a reference to a previous TxOutput
type TxInput struct {
	ID []byte
	//ID will find the Transaction that a specific output is inside of
	Out int
	//Out will be the index of the specific output we found within a transaction.
	//For example if a transaction has 4 outputs, we can use this "Out" field to specify which output we are looking for
	Sig string
	//This would be a script that adds data to an outputs' PubKey
	//however for this tutorial the Sig will be indentical to the PubKey.
}

func CoinbaseTx(toAddress, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", toAddress)
	}
	//Since this is the "first" transaction of the block, it has no previous output to reference.
	//This means that we initialize it with no ID, and it's OutputIndex is -1
	txIn := TxInput{[]byte{}, -1, data}

	txOut := TxOutput{Reward, toAddress}

	tx := Transaction{nil, []TxInput{txIn}, []TxOutput{txOut}}

	return &tx

}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	Handle(err)

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data
}
func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}

func (tx *Transaction) IsCoinbase() bool {
	//This checks a transaction and will only return true if it is a newly minted "coin"
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}
