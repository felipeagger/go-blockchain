package main

import (
	"crypto/ecdsa"
	"fmt"
	"log"

	blc "github.com/felipeagger/go-blockchain/blockchain"
	"github.com/felipeagger/go-blockchain/wallet"
)

type Wallet struct {
	Seed    string  `json:"seed"`
	Address string  `json:"address"`
	Balance float64 `json:"balance"`
}

var (
	alicePrivKey, bobPrivKey, johnPrivKey *ecdsa.PrivateKey
	alicePubKey, bobPubKey, johnPubKey    *ecdsa.PublicKey
	aliceAddress, bobAddress, johnAddress string
)

func init() {
	var err error

	alicePrivKey, alicePubKey, err = wallet.GenerateKeysFromPassword("alice")
	bobPrivKey, bobPubKey, err = wallet.GenerateKeysFromPassword("bob")
	johnPrivKey, johnPubKey, err = wallet.GenerateKeysFromPassword("john")

	if err != nil {
		panic(err)
	}

	aliceAddress = wallet.PublicKeyCompressedToString(alicePubKey)
	bobAddress = wallet.PublicKeyCompressedToString(bobPubKey)
	johnAddress = wallet.PublicKeyCompressedToString(johnPubKey)
}

func tests(execTxTests bool, blockchain *blc.Blockchain) {
	if !execTxTests {
		return
	}

	//Alice
	tx1, err := blc.NewTransaction(blockchain,
		aliceAddress,
		bobAddress,
		blc.BtcToSatoshis(0.5))
	if err != nil {
		log.Fatal(err)
	}

	tx1.Sign(alicePrivKey)

	err = blockchain.NewBlock([]blc.Transaction{tx1})
	if err != nil {
		log.Fatal(err)
	}

	//Bob
	tx2, err := blc.NewTransaction(blockchain,
		bobAddress,
		johnAddress,
		blc.BtcToSatoshis(0.2))
	if err != nil {
		log.Fatal(err)
	}

	tx2.Sign(bobPrivKey)

	err = blockchain.NewBlock([]blc.Transaction{tx2})
	if err != nil {
		log.Fatal(err)
	}

	//John
	tx3, err := blc.NewTransaction(blockchain, johnAddress, "doe", blc.BtcToSatoshis(0.1))
	if err != nil {
		log.Fatal(err)
	}

	tx3.Sign(johnPrivKey)

	tx4, err := blc.NewTransaction(blockchain, bobAddress, "jane", blc.BtcToSatoshis(0.1))
	if err != nil {
		log.Fatal(err)
	}

	tx4.Sign(bobPrivKey)

	err = blockchain.NewBlock([]blc.Transaction{tx3, tx4})
	if err != nil {
		log.Fatal(err)
	}

	isValid := blockchain.IsValid()
	fmt.Println(isValid)
}

func getWalletsData() (wallets []Wallet) {
	aliceBalance := blockchain.GetAddressBalance(aliceAddress)
	bobBalance := blockchain.GetAddressBalance(bobAddress)
	johnBalance := blockchain.GetAddressBalance(johnAddress)

	wallets = append(wallets, Wallet{Seed: "alice", Balance: blc.SatoshisToBtc(aliceBalance), Address: aliceAddress})
	wallets = append(wallets, Wallet{Seed: "bob", Balance: blc.SatoshisToBtc(bobBalance), Address: bobAddress})
	wallets = append(wallets, Wallet{Seed: "john", Balance: blc.SatoshisToBtc(johnBalance), Address: johnAddress})

	return wallets
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
