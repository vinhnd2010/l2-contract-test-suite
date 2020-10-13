package blockchain

import (
	"log"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func TestStateRoot(t *testing.T) {

	var (
		pubKey1, _   = hexutil.Decode("0xb8748a745b1c75a34238d56576e41bea9207fb5e1f7da8abe741bd9dbf14dd0e0cfb7e0cf1380065477345a42aa821aa1c68e7d9eb213eee1e8f00cb707458a4")
		accountAlloc = map[uint32]GenesisAccount{
			0: {
				Tokens: map[uint16]*big.Int{
					0: big.NewInt(30000),
					1: big.NewInt(2000000),
				},
				Pubkey:  pubKey1,
				Address: common.HexToAddress("0xdC70a72AbF352A0E3f75d737430EB896BA9Bf9Ea"),
			},
		}
	)

	state := NewStateFromAlloc(accountAlloc);
	stateRoot:= state.tree.RootHash();
	log.Printf("root tree is %s", stateRoot.String())

}
