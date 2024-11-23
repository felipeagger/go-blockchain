package wallet

// Wallet represent a key pairs
type Wallet struct {
	PrivateKey []byte // Chave privada da carteira.
	PublicKey  []byte // Chave pública derivada da privada.
	Address    string // Endereço público da carteira (derivado do hash da chave pública).
}
