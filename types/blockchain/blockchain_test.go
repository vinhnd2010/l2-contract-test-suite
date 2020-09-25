package blockchain

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/l2-contract-test-suite/types"
)

func TestNewBlockchain(t *testing.T) {
	var (
		pubKey1, _ = hexutil.Decode("0xb8748a745b1c75a34238d56576e41bea9207fb5e1f7da8abe741bd9dbf14dd0e0cfb7e0cf1380065477345a42aa821aa1c68e7d9eb213eee1e8f00cb707458a4")
		pubKey2, _ = hexutil.Decode("0xe61f3aab7e1bd78495524c955a6e3f89152ee3811fe52b85882002c465a235f7dc9bc9ed7b58277d5f9036c85e47958c65bc81104718a9364a294d96b4d277da")
		pubKey3, _ = hexutil.Decode("0x5bb440955b11980eaad949aa3f1fb05825c53cefb211b0f515415107a3aaf9dec1820b7899ad2a62a1c4aacf320b1a528c8c98aa558ee777e60110be62626e42")
	)
	bc := NewBlockchain(&Genesis{
		AccountAlloc: map[uint32]GenesisAccount{
			0: {
				Tokens: map[uint16]*big.Int{
					0: big.NewInt(30000),
					1: big.NewInt(2000000),
				},
				Pubkey:  pubKey1,
				Address: common.HexToAddress("0xdC70a72AbF352A0E3f75d737430EB896BA9Bf9Ea"),
			},

			8: {
				Tokens: map[uint16]*big.Int{
					0: big.NewInt(50000),
					1: big.NewInt(6000000),
				},
				Pubkey:  pubKey2,
				Address: common.HexToAddress("0xdC70a72AbF352A0E3f75d737430EB896BA9Bf9Ea"),
			},
			12: {
				Tokens: map[uint16]*big.Int{
					0: big.NewInt(30000),
					1: big.NewInt(1000000),
					2: big.NewInt(5000000),
				},
				Pubkey:  pubKey3,
				Address: common.HexToAddress("0x052f46FeB45822E7f117536386C51B6Bd3125157"),
			},
		},
		AccountMax: 18,
		LooMax:     0,
	})
	fmt.Println(bc.GetStateData().Hash().Hex())

	blk := &types.MiniBlock{
		Txs: []types.Transaction{
			&types.Settlement1{
				OpType:   types.SettlementOp11,
				Token1:   1,
				Token2:   2,
				Account1: 8,
				Account2: 12,
				Rate1: types.Amount{
					Mantisa: 1,
					Exp:     18,
				},
				Rate2: types.Amount{
					Mantisa: 1,
					Exp:     18,
				},
				Amount1: types.Amount{
					Mantisa: 2,
					Exp:     6,
				},
				Amount2: types.Amount{
					Mantisa: 3,
					Exp:     6,
				},
				Fee1: types.Fee{
					Mantisa: 7,
					Exp:     3,
				},
				Fee2: types.Fee{
					Mantisa: 4,
					Exp:     2,
				},
				ValidSince1:  1600661872,
				ValidSince2:  1600661873,
				ValidPeriod1: 86400,
				ValidPeriod2: 86400,
			},
		},
	}
	proofs := bc.AddMiniBlock(blk)
	fmt.Println(proofs)
}
