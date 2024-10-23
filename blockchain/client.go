package blockchain

import (
	"encoding/json"
	"net"
	"time"
)

func RequestLastBlock(peerAddr string) (*Block, error) {
	conn, err := net.Dial("tcp", peerAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Envia o tipo de requisição
	requestType := "get_last_block"
	err = json.NewEncoder(conn).Encode(requestType)
	if err != nil {
		return nil, err
	}

	// Lê a resposta do servidor
	var lastBlock Block
	err = json.NewDecoder(conn).Decode(&lastBlock)
	if err != nil {
		return nil, err
	}

	return &lastBlock, nil
}

func RequestBlocksFromHash(peerAddress, hash string) (*Blockchain, error) {
	conn, err := net.DialTimeout("tcp", peerAddress, 10*time.Second)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Envia o tipo de requisição
	requestType := "get_blocks"
	err = json.NewEncoder(conn).Encode(requestType)
	if err != nil {
		return nil, err
	}

	err = json.NewEncoder(conn).Encode(hash)
	if err != nil {
		return nil, err
	}

	// Lê os blocos recebidos
	var blocks []Block
	err = json.NewDecoder(conn).Decode(&blocks)
	if err != nil {
		return nil, err
	}

	return &Blockchain{Chain: blocks}, nil
}
