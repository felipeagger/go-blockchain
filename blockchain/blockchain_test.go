package blockchain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func getTestData() []Block {
	var blocks []Block

	blocks = append(blocks, Block{Hash: "0", Transactions: []Transaction{
		{
			ID:        []byte("genesis"),
			Timestamp: time.Now(),
			Outputs: []TxOutput{
				{
					Value:  10_000_000_000,
					PubKey: RewardTo,
				},
			},
		},
	}})

	//alice to bob
	blocks = append(blocks, Block{
		Hash:         "0000e6657636d3f0b8c32618f2d6b2526febcd9a6cb1d196ad97699037123b40",
		PreviousHash: "0",
		Transactions: []Transaction{
			{
				ID:        []byte("VOtL246MPvboxfIcjdmkpoSsf08zYNGyZiO6uhqLCIQ="),
				Timestamp: time.Now(),
				Inputs: []TxInput{
					{
						ID:        []byte("Z2VuZXNpcw=="),
						OutIdx:    0,
						Signature: "80740746c45978415b949e359dd0cd19558cdac74eedfe86d9ec5814512f296cb81b65119b421582c9ad90a6d65987c23b79711aa6e3a2b97ca9f9cb6b09579d",
						PubKey:    RewardTo,
					},
				},
				Outputs: []TxOutput{
					{
						Value:  5_000_000,
						PubKey: "03a44d59c5bfbb36b173e729752dc1748a3f2555991b13f575853960158f91fa80",
					},
					{
						Value:  10_000_000_000 - 5_000_000,
						PubKey: RewardTo,
					},
				},
			},
		}})

	//bob to john
	blocks = append(blocks, Block{
		Hash:         "00003cfc1fb4202323095efeae9573b605309b8bb30045de8c56eb29f17484e4",
		PreviousHash: "0000e6657636d3f0b8c32618f2d6b2526febcd9a6cb1d196ad97699037123b40",
		Transactions: []Transaction{
			{
				ID:        []byte("HX91UCstxSSL06fzz6WXz2Iub+MXTITgraXON1mPnOs="),
				Timestamp: time.Now(),
				Inputs: []TxInput{
					{
						ID:        []byte("VOtL246MPvboxfIcjdmkpoSsf08zYNGyZiO6uhqLCIQ="),
						OutIdx:    0,
						Signature: "0cf5eb2b071e07205ceee47589b06c642b1142178efc77ff7e5581dc880344b7f4b1f64b5e3d048fb8bf2b225a3b4ef27aa76ef30b65e7735d54b6393f649ba8",
						PubKey:    "03a44d59c5bfbb36b173e729752dc1748a3f2555991b13f575853960158f91fa80",
					},
				},
				Outputs: []TxOutput{
					{
						Value:  5_000_000,
						PubKey: "02d14c44e48412e9e7a9b153d02f3933e74e5ed4045304040ba5ac923939fffe5d",
					},
				},
			},
		}})

	return blocks
}

func TestFindUnspentTransactions(t *testing.T) {
	blockchain := Blockchain{Chain: getTestData()}

	bobAddress := "03a44d59c5bfbb36b173e729752dc1748a3f2555991b13f575853960158f91fa80"

	result := blockchain.FindUnspentTransactions(bobAddress)

	assert.Len(t, result, 0)
}
