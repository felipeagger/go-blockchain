package main

import (
	"encoding/json"
	"fmt"
	blc "github.com/felipeagger/go-blockchain/blockchain"
	"net"
	"sync"
)

var mu sync.Mutex

func startServer(blockchain *blc.Blockchain, port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Servidor TCP escutando na porta:", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}
		go handleConnection(conn, blockchain) // Lida com a conexão em uma nova goroutine
	}
}

func handleConnection(conn net.Conn, blockchain *blc.Blockchain) {
	defer conn.Close()
	fmt.Println("conexão aceita!")

	var requestType string
	err := json.NewDecoder(conn).Decode(&requestType)
	if err != nil {
		fmt.Println("Erro ao ler requisição:", err)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	switch requestType {
	case "get_last_block":
		if len(blockchain.Chain) > 0 {
			lastBlock := blockchain.Chain[len(blockchain.Chain)-1]
			err := json.NewEncoder(conn).Encode(lastBlock)
			if err != nil {
				fmt.Println("Erro ao enviar o último bloco:", err)
			}
		} else {
			conn.Write([]byte("Nenhum bloco disponível"))
		}

	case "get_blocks":
		var startHash string
		err := json.NewDecoder(conn).Decode(&startHash)
		if err != nil {
			fmt.Println("Erro ao ler hash:", err)
			return
		}

		var blocksToSend []blc.Block
		found := false

		// Enviar blocos a partir do hash fornecido
		for _, block := range blockchain.Chain {
			if block.Hash == startHash {
				found = true
			}
			if found {
				blocksToSend = append(blocksToSend, block)
			}
		}

		err = json.NewEncoder(conn).Encode(blocksToSend)
		if err != nil {
			fmt.Println("Erro ao enviar blocos:", err)
		}

	default:
		fmt.Println("Tipo de requisição desconhecido:", requestType)
	}
}
