package blockchain

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

func Synchronize(db *sql.DB, nodePool []string, difficulty int) error {
	fmt.Println("Synchronize blocks on pool...")

	lastBlock, err := GetLastBlock(db)
	if err != nil {
		if err.Error() == "block not found" {
			lastBlock = Block{Hash: "0"}
		} else {
			return err
		}
	}

	for _, nodeAddr := range nodePool {
		if strings.TrimSpace(nodeAddr) == "" {
			continue
		}

		err = syncWithPeer(strings.TrimSpace(nodeAddr), db, lastBlock.Hash, difficulty)
		if err != nil {
			fmt.Println("Failed to synchronize blocks on node:", nodeAddr)
			continue
		}
	}

	return err
}

func syncWithPeer(nodeAddr string, db *sql.DB, localLastBlockHash string, difficulty int) error {
	// Solicita o último bloco do peer
	blockRes, err := RequestLastBlock(nodeAddr)
	if err != nil {
		return err
	}

	if blockRes.Hash != localLastBlockHash {
		// Solicita os blocos faltantes
		remoteNodeBlockchain, err := RequestBlocksFromHash(nodeAddr, localLastBlockHash)
		if err != nil {
			return err
		}

		//valida blockchain
		remoteNodeBlockchain.difficulty = difficulty
		if !remoteNodeBlockchain.IsValid() {
			return errors.New("remote node blockchain is invalid")
		}

		for _, block := range remoteNodeBlockchain.Chain {

			existLocalBlk, _ := GetBlock(db, block.Hash)
			if existLocalBlk.Hash == "" {
				fmt.Println("Novo bloco recebido:", block.Hash)

				err = InsertBlock(db, block)
				if err != nil {
					fmt.Println("Erro ao inserir bloco recebido:", block.Hash)
				}
			}
		}
	}

	return nil
}
