package main

import (
	"encoding/json"
	"io/ioutil"
	"math/big"

	"github.com/KyberNetwork/l2-contract-test-suite/common/proof"
	"github.com/KyberNetwork/l2-contract-test-suite/testsample"
	"github.com/KyberNetwork/l2-contract-test-suite/types"
	"github.com/KyberNetwork/l2-contract-test-suite/types/blockchain"
	"github.com/KyberNetwork/l2-contract-test-suite/types/test"
	"github.com/ethereum/go-ethereum/common"
)

const testOutput = "testdata/fraudProofWithdraw.json"

var genesis = &blockchain.Genesis{
	AccountAlloc: map[uint32]blockchain.GenesisAccount{
		0: {
			Tokens:  map[uint16]*big.Int{},
			Pubkey:  testsample.PublicKeys[1],
			Address: common.HexToAddress("0xdC70a72AbF352A0E3f75d737430EB896BA9Bf9Ea"),
		},

		23: {
			Tokens: map[uint16]*big.Int{
				0: big.NewInt(2000),
			},
			Pubkey:  testsample.PublicKeys[2],
			Address: common.HexToAddress("0xdC70a72AbF352A0E3f75d737430EB896BA9Bf9Ea"),
		},
		47: {
			Tokens:  map[uint16]*big.Int{},
			Pubkey:  testsample.PublicKeys[3],
			Address: common.HexToAddress("0x052f46FeB45822E7f117536386C51B6Bd3125157"),
		},
	},
	AccountMax: 1000,
	LooMax:     289,
	LooAlloc: map[uint64]*types.LeftOverOrder{
		56: {
			AccountID:   30,
			SrcToken:    2,
			DestToken:   1,
			Amount:      big.NewInt(4321),
			Fee:         big.NewInt(600),
			Rate:        types.PackedAmount{Mantisa: 4, Exp: 18}.Big(),
			ValidSince:  1601436626,
			ValidPeriod: 823000,
		},
		243: {
			AccountID:   17,
			SrcToken:    1,
			DestToken:   2,
			Amount:      big.NewInt(34500),
			Fee:         big.NewInt(67432),
			Rate:        types.PackedAmount{Mantisa: 2, Exp: 17}.Big(),
			ValidSince:  1601436627,
			ValidPeriod: 823000,
		},
	},
}

func buildTest1() *test.Suit {
	bc := blockchain.NewBlockchain(genesis)
	genesisHash := bc.GetStateData().Hash()
	// create an deposit to user
	deposit := &types.DepositOp{
		AccountID: 23,
		TokenID:   2,
		Amount:    big.NewInt(45242000),
	}
	miniBlock1 := &types.MiniBlock{
		Txs: []types.Transaction{
			deposit,
		},
	}
	bc.AddMiniBlock(miniBlock1)
	submitBlockStep := test.SubmitBlockStep{
		BlockNumber: 1,
		MiniBlocks:  []*types.MiniBlock{miniBlock1},
		Timestamp:   1600661872,
	}
	// create an withdraw to another user
	withdraw := &types.WithdrawOp{
		TokenID:    2,
		Amount:     types.PackedAmount{Mantisa: 4, Exp: 7},
		DestAddr:   common.HexToAddress("0x99aF5AF1f1a61FE1678e030916f79331a28A57E8"),
		AccountID:  23,
		ValidSince: 0,
		Fee:        types.PackedFee{Mantisa: 1, Exp: 2},
	}
	miniBlock2 := &types.MiniBlock{
		Txs: []types.Transaction{
			withdraw,
		},
	}
	prevStateData := bc.GetStateData()
	executionProofs := bc.AddMiniBlock(miniBlock2)
	submitBlockStep2 := test.SubmitBlockStep{
		BlockNumber: 2,
		MiniBlocks:  []*types.MiniBlock{miniBlock2},
		Timestamp:   1600661872,
	}

	accuseStep := test.AccuseBlockFraudProofStep{
		BlockNumber:        2,
		MiniBlockNumber:    0,
		MiniBlock:          miniBlock2,
		PrevStateData:      prevStateData,
		MiniBlockProof:     proof.BuildMiniBlockProof(submitBlockStep2.MiniBlocks, 0, submitBlockStep2.Timestamp),
		PrevStateHashProof: proof.BuildFinalStateHashProof(submitBlockStep.MiniBlocks, submitBlockStep.Timestamp),
		ExecutionProof:     executionProofs,
	}

	//blockData2.Proof.MiniBlockProof = proof.BuildMiniBlockProof(blockData2.MiniBlocks, uint(blockData2.MiniBlockNumber), blockData2.Timestamp)

	return &test.Suit{
		Msg:              "test case when withdraw",
		GenesisStateHash: genesisHash,
		AccountMax:       genesis.AccountMax,
		Steps: []test.Step{
			{Action: test.SubmitDeposit, Data: deposit},
			{Action: test.SubmitBlock, Data: submitBlockStep},
			{Action: test.SubmitBlock, Data: submitBlockStep2},
			{Action: test.AccuseBlockFraudProof, Data: accuseStep},
			{Action: test.CompleteWithdraw, Data: withdraw},
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
	if err := ioutil.WriteFile(testOutput, b, 0644); err != nil {
		panic(err)
	}
}
