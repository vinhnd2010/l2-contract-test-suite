package main

import (
	"encoding/json"
	"io/ioutil"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	util "github.com/KyberNetwork/l2-contract-test-suite/common"
	"github.com/KyberNetwork/l2-contract-test-suite/common/proof"
	"github.com/KyberNetwork/l2-contract-test-suite/testsample"
	"github.com/KyberNetwork/l2-contract-test-suite/types"
	"github.com/KyberNetwork/l2-contract-test-suite/types/blockchain"
	"github.com/KyberNetwork/l2-contract-test-suite/types/test"
)

const testOutput = "testdata/fraudProofExit.json"

var genesis = &blockchain.Genesis{
	AccountAlloc: map[uint32]blockchain.GenesisAccount{
		0: {
			Tokens:  map[uint16]*big.Int{},
			Pubkey:  testsample.PublicKeys[1],
			Address: common.HexToAddress("0xdC70a72AbF352A0E3f75d737430EB896BA9Bf9Ea"),
		},

		36: {
			Tokens: map[uint16]*big.Int{
				5: big.NewInt(500),
			},
			Pubkey:  testsample.PublicKeys[2],
			Address: common.HexToAddress("0xc783df8a850f42e7F7e57013759C285caa701eB6"), // address for accounts[0] from buidler config
		},
		44: {
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
		AccountID: 36,
		TokenID:   2,
		Amount:    big.NewInt(45242000),
	}

	deposit2 := &types.DepositOp{
		AccountID: 36,
		TokenID:   4,
		Amount:    big.NewInt(135000),
	}
	miniBlock1 := &types.MiniBlock{
		Txs: []types.Transaction{deposit, deposit2},
	}
	bc.AddMiniBlock(miniBlock1)
	submitBlockStep := test.SubmitBlockStep{
		BlockNumber: 1,
		MiniBlocks:  []*types.MiniBlock{miniBlock1},
		Timestamp:   1600661872,
	}
	// create exit step
	submitExitStep := test.SubmitExitStep{
		AccountID:   36,
		Timestamp:   submitBlockStep.Timestamp,
		BlockNumber: submitBlockStep.BlockNumber,
	}
	var exitProof []byte
	submitExitStep.BalanceRoot, exitProof = bc.BuildSubmitExitProof(submitExitStep.AccountID)
	// due to this only have 1 block so blockDataHash = miniBlock1.Hash()
	exitProof = append(exitProof, miniBlock1.Hash().Bytes()...)
	exitProof = append(exitProof, util.Uint8ToByte(1))
	submitExitStep.Proof = exitProof

	// create an withdraw to another user
	exit := &types.ExitOp{
		AccountID: 36,
	}
	miniBlock2 := &types.MiniBlock{
		Txs: []types.Transaction{exit},
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
	// create complete exit step
	completeExitStep := test.CompleteExitStep{
		AccountID: 36,
		TokenIDs:  []uint16{2, 4},
	}
	completeExitStep.TokenAmounts, completeExitStep.Siblings = bc.BuildCompleteExit(completeExitStep.AccountID, completeExitStep.TokenIDs)

	return &test.Suit{
		Msg:              "test case when exit with 2 tokens",
		GenesisStateHash: genesisHash,
		Steps: []test.Step{
			{Action: test.SubmitDeposit, Data: deposit},
			{Action: test.SubmitDeposit, Data: deposit2},
			{Action: test.SubmitBlock, Data: submitBlockStep},
			{Action: test.SubmitExit, Data: submitExitStep},
			{Action: test.SubmitBlock, Data: submitBlockStep2},
			{Action: test.AccuseBlockFraudProof, Data: accuseStep},
			{Action: test.CompleteExit, Data: completeExitStep},
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
