package main

import (
	"encoding/json"
	"io/ioutil"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/l2-contract-test-suite/common/proof"
	"github.com/KyberNetwork/l2-contract-test-suite/testsample"
	"github.com/KyberNetwork/l2-contract-test-suite/types"
	"github.com/KyberNetwork/l2-contract-test-suite/types/blockchain"
)

var (
	output = "testdata/fraudProofDepositToNew.json"
)

var genesis = &blockchain.Genesis{
	AccountAlloc: map[uint32]blockchain.GenesisAccount{
		0: {
			Tokens: map[uint16]*big.Int{
				0: big.NewInt(30000),
				1: big.NewInt(2000000),
			},
			Pubkey:  testsample.PublicKeys[0],
			Address: common.HexToAddress("0xdC70a72AbF352A0E3f75d737430EB896BA9Bf9Ea"),
		},

		8: {
			Tokens: map[uint16]*big.Int{
				0: big.NewInt(50000),
				1: big.NewInt(6000000),
			},
			Pubkey:  testsample.PublicKeys[2],
			Address: common.HexToAddress("0xdC70a72AbF352A0E3f75d737430EB896BA9Bf9Ea"),
		},
		12: {
			Tokens: map[uint16]*big.Int{
				0: big.NewInt(30000),
				1: big.NewInt(1000000),
				2: big.NewInt(5000000),
			},
			Pubkey:  testsample.PublicKeys[4],
			Address: common.HexToAddress("0x052f46FeB45822E7f117536386C51B6Bd3125157"),
		},
	},
	AccountMax: 18,
	LooMax:     0,
}

type DepositFraudProofTestSuit struct {
	Msg              string
	GenesisStateHash common.Hash
	DepositOp        *types.DepositToNewOp
	Blocks           []blockchain.BlockData
}

func buildTest1() *DepositFraudProofTestSuit {
	bc := blockchain.NewBlockchain(genesis)
	genesisHash := bc.GetStateData().Hash()
	preStateData := bc.GetStateData()

	deposit := &types.DepositToNewOp{
		PubKey:     testsample.PublicKeys[6],
		WithdrawTo: common.HexToAddress("0x91F4d9EA5c1ee0fc778524b3D57fD8CF700996Cf"),
		TokenID:    2,
		Amount:     big.NewInt(45242000),
	}

	miniBlock1 := &types.MiniBlock{
		Txs: []types.Transaction{deposit},
	}

	executionProofs := bc.AddMiniBlock(miniBlock1)
	blockData := blockchain.BlockData{
		MiniBlocks:      []*types.MiniBlock{miniBlock1},
		Timestamp:       1600661872,
		MiniBlockNumber: 0,
		Proof: &blockchain.FraudProof{
			PrevStateData:      preStateData,
			PrevStateHashProof: []byte{},
			ExecutionProof:     executionProofs,
		},
	}
	blockData.Proof.MiniBlockProof = proof.BuildMiniBlockProof(blockData.MiniBlocks, blockData.MiniBlockNumber, blockData.Timestamp)
	return &DepositFraudProofTestSuit{
		Msg:              "test case simple deposit",
		GenesisStateHash: genesisHash,
		DepositOp:        deposit,
		Blocks: []blockchain.BlockData{
			blockData,
		},
	}
}

func main() {
	var testSuits []*DepositFraudProofTestSuit
	testSuits = append(testSuits, buildTest1())

	b, err := json.MarshalIndent(testSuits, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(output, b, 0644); err != nil {
		panic(err)
	}
}
