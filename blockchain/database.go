package blockchain

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

func InitializeDatabase(dbPath string) *sql.DB {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func ExecMigrations(db *sql.DB) error {
	fmt.Println("Executing migrations...")

	createBlocksTableSQL := `
    CREATE TABLE IF NOT EXISTS blocks (
        hash TEXT PRIMARY KEY NOT NULL,
        data TEXT NOT NULL,
        previous_hash TEXT NOT NULL,
        timestamp DATETIME NOT NULL,
        pow INTEGER NOT NULL
    );`
	//_, err := db.Exec(createBlocksTableSQL)
	//if err != nil {
	//	log.Fatal(err)
	//}

	statement, err := db.Prepare(createBlocksTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	_, err = statement.Exec()
	return err
}

func InsertBlock(db *sql.DB, block Block) error {
	jsonData, err := json.Marshal(block.Data)
	if err != nil {
		log.Fatal(err)
	}

	insertSQL := `INSERT INTO blocks (hash, data, previous_hash, timestamp, pow) VALUES (?, ?, ?, ?, ?)`
	_, err = db.Exec(insertSQL, block.Hash, string(jsonData), block.PreviousHash, block.Timestamp, block.Pow)
	return err
}

func GetBlock(db *sql.DB, hash string) (Block, error) {
	var block Block
	var jsonData string

	row := db.QueryRow(`SELECT hash, data, previous_hash, timestamp, pow FROM blocks where hash = ?`, hash)

	err := row.Scan(&block.Hash, &jsonData, &block.PreviousHash, &block.Timestamp, &block.Pow)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return block, fmt.Errorf("block not found with this hash: %s", hash)
		}
		return block, err
	}

	err = json.Unmarshal([]byte(jsonData), &block.Data)
	return block, err
}

func LoadBlockchain(db *sql.DB, difficulty int) (*Blockchain, error) {
	_blockchain := &Blockchain{}
	var blocks []Block

	query := `SELECT hash, data, previous_hash, timestamp, pow FROM blocks ORDER BY timestamp DESC LIMIT 100`

	rows, err := db.Query(query)
	if err != nil {
		return _blockchain, err
	}
	defer rows.Close()

	for rows.Next() {
		var block Block
		var jsonData string

		err := rows.Scan(&block.Hash, &jsonData, &block.PreviousHash, &block.Timestamp, &block.Pow)
		if err != nil {
			return _blockchain, err
		}

		err = json.Unmarshal([]byte(jsonData), &block.Data)
		if err != nil {
			return _blockchain, err
		}

		blocks = append(blocks, block)
	}

	if err := rows.Err(); err != nil {
		return _blockchain, err
	}

	_blockchain.Chain = reverseBlocks(blocks)
	_blockchain.difficulty = difficulty
	_blockchain.db = db

	return _blockchain, nil
}

func GetLastBlock(db *sql.DB) (Block, error) {
	var block Block
	var jsonData string

	query := `SELECT hash, data, previous_hash, timestamp, pow FROM blocks ORDER BY timestamp DESC LIMIT 1`

	row := db.QueryRow(query)
	err := row.Scan(&block.Hash, &jsonData, &block.PreviousHash, &block.Timestamp, &block.Pow)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return block, fmt.Errorf("block not found")
		}
		return block, err
	}

	err = json.Unmarshal([]byte(jsonData), &block.Data)
	return block, err
}
