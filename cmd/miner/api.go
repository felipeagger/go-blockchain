package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

	http.HandleFunc("/api/wallet-balance", GetWalletsBalance)

	http.HandleFunc("/api/transaction", NewTransaction)

	http.Handle("/", http.FileServer(http.Dir("/home/felipeagger/Dados/Dev/Projects/go-blockchain/static")))

	fmt.Println("Servidor rodando em http://localhost:8088")
	http.ListenAndServe(":8088", nil)
}

func GetWalletsBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Error method now allowed", http.StatusMethodNotAllowed)
	}

	walletsData := getWalletsData()

	blocksBytes, _ := json.Marshal(walletsData)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(blocksBytes)
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
