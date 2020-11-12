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
	"github.com/KyberNetwork/l2-contract-test-suite/types/test"
)

var (
	output = "testdata/fraudProofDeposit.json"
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
			Pubkey:  testsample.PublicKeys[5],
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

func buildTest1() *test.Suit {
	bc := blockchain.NewBlockchain(genesis)
	genesisHash := bc.GetStateData().Hash()
	preStateData := bc.GetStateData()

	deposit := &types.DepositOp{
		AccountID: 8,
		TokenID:   2,
		Amount:    big.NewInt(45242000),
	}

	miniBlock1 := &types.MiniBlock{
		Txs: []types.Transaction{deposit},
	}
	executionProofs := bc.AddMiniBlock(miniBlock1)

	submitBlockStep := test.SubmitBlockStep{
		BlockNumber: 1,
		MiniBlocks:  []*types.MiniBlock{miniBlock1},
		Timestamp:   1600661872,
	}

	accuseStep := test.AccuseBlockFraudProofStep{
		BlockNumber:        1,
		MiniBlockNumber:    0,
		MiniBlock:          miniBlock1,
		PrevStateData:      preStateData,
		MiniBlockProof:     proof.BuildMiniBlockProof(submitBlockStep.MiniBlocks, uint(submitBlockStep.BlockNumber), submitBlockStep.Timestamp),
		PrevStateHashProof: proof.BuildFinalStateHashProof(submitBlockStep.MiniBlocks, submitBlockStep.Timestamp),
		ExecutionProof:     executionProofs,
	}

	return &test.Suit{
		Msg:              "test case simple deposit",
		GenesisStateHash: genesisHash,
		AccountMax:       genesis.AccountMax,
		Steps: []test.Step{
			{Action: test.SubmitDeposit, Data: deposit},
			{Action: test.SubmitBlock, Data: submitBlockStep},
			{Action: test.AccuseBlockFraudProof, Data: accuseStep},
		},
	}
}

func main() {
	var testSuits []*test.Suit
	testSuits = append(testSuits, buildTest1())

	b, err := json.MarshalIndent(testSuits, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(output, b, 0644); err != nil {
		panic(err)
	}
}
