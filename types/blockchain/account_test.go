package blockchain

import (
	"log"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func TestAccountBalanceRoot(t *testing.T) {

	var (
		pubKey1, _ = hexutil.Decode("0xb8748a745b1c75a34238d56576e41bea9207fb5e1f7da8abe741bd9dbf14dd0e0cfb7e0cf1380065477345a42aa821aa1c68e7d9eb213eee1e8f00cb707458a4")
		ethAddr    = common.HexToAddress("0xdC70a72AbF352A0E3f75d737430EB896BA9Bf9Ea")
		acc        = NewAccount(pubKey1, ethAddr)
	)
	acc.Update(10, big.NewInt(1000))
	acc.Update(1, big.NewInt(200))
	balanceRoot := acc.tree.RootHash()

	log.Printf("root tree is %s", balanceRoot.String())

}
