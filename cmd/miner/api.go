package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	blc "github.com/felipeagger/go-blockchain/blockchain"
	"github.com/felipeagger/go-blockchain/wallet"
)

type Transaction struct {
	Seed   string  `json:"seed"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

func api() {
	http.HandleFunc("/api/blocks", func(w http.ResponseWriter, r *http.Request) {
		blocksBytes, _ := json.Marshal(blockchain.Chain)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(blocksBytes)
	})

	http.HandleFunc("/api/transaction", NewTransaction)

	http.Handle("/", http.FileServer(http.Dir("/home/felipeagger/Dados/Dev/Projects/go-blockchain/static")))

	fmt.Println("Servidor rodando em http://localhost:8088")
	http.ListenAndServe(":8088", nil)
}

func NewTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Error method now allowed", http.StatusMethodNotAllowed)
	}

	var tx Transaction
	err := json.NewDecoder(r.Body).Decode(&tx)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusUnprocessableEntity)
		return
	}

	if tx.Seed == "" || tx.To == "" || tx.Amount == 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(`{"status": "error", "msg": "invalid params"}`))
		return
	}

	err = createNewTransaction(tx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"status": "error", "msg": "%s"}`, err.Error())))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "success"}`))
}

func createNewTransaction(tx Transaction) error {
	privKey, pubKey, err := wallet.GenerateKeysFromPassword(tx.Seed)
	if err != nil {
		return err
	}

	from := wallet.PublicKeyCompressedToString(pubKey)

	newTx, err := blc.NewTransaction(blockchain, from, tx.To, blc.BtcToSatoshis(tx.Amount))
	if err != nil {
		fmt.Println(err)
		return err
	}

	newTx.Sign(privKey)

	return blockchain.NewBlock([]blc.Transaction{newTx})
}
