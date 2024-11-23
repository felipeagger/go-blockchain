package wallet

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	privKeyFilename = "./priv.Key"
	myPasswordSeed  = "super secure seed pass"
	myPrivKey       = "c7f15401752b8952d63ed09eca2590492135c807deb5fa5487aad5bb3ad7e08e"
	myPubKey        = "03a7eceb218b78d4356a2a774411bb13b117a67ff70ce123ca485c9173c03a6508"
	mySomeData      = "Lorem ipsum dolor sit amet, consectetur adipiscing elit"
	mySignatureData = "11be9e75b18b942774436d0a4c74d63bd018899ab42e45699a091ff15e24971558ee63a5629852463fe2fa0d9ff057de7ac279dcd7c68ee2cc74355c1801052f"
)

func TestGenerateKeyPairs(t *testing.T) {

	privKey, pubKey, err := GenerateKeysFromPassword(myPasswordSeed)
	assert.NoError(t, err)
	assert.NotNil(t, privKey)

	privKeyStr := hex.EncodeToString(privKey.D.Bytes())
	pubKeyStr := PublicKeyCompressedToString(pubKey)

	assert.Equal(t, myPrivKey, privKeyStr)
	assert.Equal(t, myPubKey, pubKeyStr)

	err = SavePrivateKeyToFile(privKeyFilename, privKey)
	assert.NoError(t, err)
}

func TestLoadKeyPairs(t *testing.T) {

	privKey, err := LoadPrivateKeyFromFile(privKeyFilename)
	assert.NoError(t, err)
	assert.NotNil(t, privKey)

	privKeyStr := hex.EncodeToString(privKey.D.Bytes())
	pubKeyStr := PublicKeyCompressedToString(&privKey.PublicKey)

	assert.Equal(t, myPrivKey, privKeyStr)
	assert.Equal(t, myPubKey, pubKeyStr)
}

func TestCompressedPublicKeyToECDSA(t *testing.T) {
	pubKey, err := PublicKeyToECDSA(myPubKey)
	assert.NoError(t, err)
	assert.NotNil(t, pubKey)

	pubKeyStr := PublicKeyCompressedToString(pubKey)
	assert.Equal(t, myPubKey, pubKeyStr)
}

func TestSingData(t *testing.T) {
	privKey, pubKey, err := GenerateKeysFromPassword(myPasswordSeed)
	assert.NoError(t, err)
	assert.NotNil(t, privKey)

	signature, err := SignData(privKey, []byte(mySomeData))
	assert.NoError(t, err)
	assert.NotEmpty(t, signature)

	fmt.Println(hex.EncodeToString(signature))

	isValid := VerifySignature(pubKey, []byte(mySomeData), signature)
	assert.Equal(t, true, isValid)
}

func TestVerifySignatureData(t *testing.T) {
	pubKey, err := PublicKeyToECDSA(myPubKey)
	assert.NoError(t, err)
	assert.NotNil(t, pubKey)

	signature, err := hex.DecodeString(mySignatureData)

	isValid := VerifySignature(pubKey, []byte(mySomeData), signature)
	assert.Equal(t, true, isValid)
}
