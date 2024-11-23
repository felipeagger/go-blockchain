package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

const (
	salt       = "unique-salt"
	iterations = 4096
)

func GenerateKeysFromPassword(password string) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	seed := pbkdf2.Key([]byte(password), []byte(salt), iterations, 32, sha256.New)

	// Generate private key
	curve := elliptic.P256()
	n := curve.Params().N
	d := new(big.Int).SetBytes(seed)
	d.Mod(d, new(big.Int).Sub(n, big.NewInt(1))) // d = seed % (n-1)
	d.Add(d, big.NewInt(1))                      // d = d + 1 (para evitar zero)

	// Construir a chave privada
	privKey := &ecdsa.PrivateKey{
		D: d,
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
		},
	}

	privKey.Curve = curve

	// Derivar a chave pública
	privKey.PublicKey.X, privKey.PublicKey.Y = curve.ScalarBaseMult(d.Bytes())

	return privKey, &privKey.PublicKey, nil
}

// Função para salvar a chave privada em um arquivo
func SavePrivateKeyToFile(filename string, privKey *ecdsa.PrivateKey) error {
	privKeyBytes := privKey.D.Bytes()
	return os.WriteFile(filename, privKeyBytes, 0600)
}

func LoadPrivateKeyFromFile(filename string) (*ecdsa.PrivateKey, error) {
	// Ler os bytes da chave privada salvos no arquivo
	privKeyBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error on read file: %w", err)
	}

	// Recriar a chave privada
	privKey := new(ecdsa.PrivateKey)
	privKey.D = bigFromBytes(privKeyBytes)
	privKey.Curve = elliptic.P256()
	privKey.PublicKey.X, privKey.PublicKey.Y = privKey.Curve.ScalarBaseMult(privKey.D.Bytes())

	return privKey, nil
}

func PublicKeyCompressedToString(pubKey *ecdsa.PublicKey) string {
	prefix := byte(0x02)
	if pubKey.Y.Bit(0) == 1 {
		prefix = 0x03
	}

	key := append([]byte{prefix}, pubKey.X.Bytes()...)
	return hex.EncodeToString(key)
}

// Assinar dados usando a chave privada
func SignData(privKey *ecdsa.PrivateKey, data []byte) ([]byte, error) {
	hash := sha256.Sum256(data)
	r, s, err := ecdsa.Sign(rand.Reader, privKey, hash[:])
	if err != nil {
		return nil, err
	}

	signature := append(r.Bytes(), s.Bytes()...)
	return signature, nil
}

func VerifySignature(pubKey *ecdsa.PublicKey, data []byte, signature []byte) bool {
	hash := sha256.Sum256(data)

	// Split the signature into r and s
	r := new(big.Int).SetBytes(signature[:len(signature)/2])
	s := new(big.Int).SetBytes(signature[len(signature)/2:])

	// Verify the signature
	valid := ecdsa.Verify(pubKey, hash[:], r, s)
	return valid
}

func PublicKeyToECDSA(pubKey string) (*ecdsa.PublicKey, error) {
	pubKeyBytes, err := hex.DecodeString(pubKey)
	if err != nil {
		return nil, err
	}

	if len(pubKeyBytes) != 33 {
		return nil, err // Invalid length for a compressed public key
	}

	curve := elliptic.P256()
	x := new(big.Int).SetBytes(pubKeyBytes[1:])
	prefix := pubKeyBytes[0]

	// Calculate y² = x³ - 3x + b (mod p)
	ySquared := new(big.Int).Exp(x, big.NewInt(3), curve.Params().P) // x³
	threeX := new(big.Int).Mul(big.NewInt(3), x)                     // 3x
	threeX.Mod(threeX, curve.Params().P)                             // mod p
	ySquared.Sub(ySquared, threeX)                                   // x³ - 3x
	ySquared.Add(ySquared, curve.Params().B)                         // x³ - 3x + b
	ySquared.Mod(ySquared, curve.Params().P)                         // mod p

	// Compute the modular square root of y²
	y := new(big.Int).ModSqrt(ySquared, curve.Params().P)
	if y == nil {
		return nil, fmt.Errorf("Invalid compressed public key")
	}

	// Choose the correct y based on prefix
	if (prefix == 0x03) != (y.Bit(0) == 1) {
		y.Sub(curve.Params().P, y) // Negate y mod P
	}

	return &ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}, nil
}

// Função auxiliar para criar um big.Int a partir de bytes
func bigFromBytes(b []byte) *big.Int {
	return new(big.Int).SetBytes(b)
}
