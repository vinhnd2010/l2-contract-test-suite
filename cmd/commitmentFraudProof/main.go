package main

import (
	"encoding/json"
	"io/ioutil"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/KyberNetwork/l2-contract-test-suite/common/proof"
	"github.com/KyberNetwork/l2-contract-test-suite/testsample"
	"github.com/KyberNetwork/l2-contract-test-suite/types"
	"github.com/KyberNetwork/l2-contract-test-suite/types/blockchain"
	"github.com/KyberNetwork/l2-contract-test-suite/types/test"
)

const (
	commitmentBuilderTest1Output = "testdata/commitmentBuilder1.json"
	commitmentBuilderTest2Output = "testdata/commitmentBuilder2.json"
	commitmentBuilderTest3Output = "testdata/commitmentBuilder3.json"

	commitmentFraudProofTest1Output = "testdata/commitmentFraudProof1.json"
	commitmentFraudProofTest2Output = "testdata/commitmentFraudProof2.json"
	commitmentFraudProofTest3Output = "testdata/commitmentFraudProof3.json"
)

type CommitmentBuilderTest1 struct {
	TxData                  hexutil.Bytes `json:"txData"`
	ExpectedCommitmentInput hexutil.Bytes `json:"commitmentInput"`
	Account1                uint32        `json:"accountID1"`
	Account2                uint32        `json:"accountID2"`
	AccountPubKey1          hexutil.Bytes `json:"accountPubKey1"`
	AccountPubKey2          hexutil.Bytes `json:"accountPubKey2"`
}

type CommitmentBuilderTest2 struct {
	TxData                  hexutil.Bytes `json:"txData"`
	ExpectedCommitmentInput hexutil.Bytes `json:"commitmentInput"`
	AccountID               uint32        `json:"accountID"`
	AccountPubKey           hexutil.Bytes `json:"accountPubKey"`
	LooID                   uint64        `json:"looID"`
	LooSrcToken             uint16        `json:"looSrcToken"`
	LooDstToken             uint16        `json:"looDstToken"`
}

type CommitmentBuilderTest3 struct {
	TxData                  hexutil.Bytes `json:"txData"`
	ExpectedCommitmentInput hexutil.Bytes `json:"commitmentInput"`
	AccountID               uint32        `json:"accountID"`
	AccountPubKey           hexutil.Bytes `json:"accountPubKey"`
}

func testCommitmentBuilder1() {
	settlement1 := &types.Settlement1{
		OpType:       types.SettlementOp12,
		Token1:       0,
		Token2:       6,
		Account1:     1,
		Account2:     2,
		Amount1:      types.PackedAmount{Mantisa: 3, Exp: 18},
		Rate1:        types.PackedAmount{Mantisa: 2, Exp: 16},
		Amount2:      types.PackedAmount{Mantisa: 2, Exp: 16},
		Rate2:        types.PackedAmount{Mantisa: 5, Exp: 19},
		Fee1:         types.PackedFee{Mantisa: 3, Exp: 5},
		Fee2:         types.PackedFee{Mantisa: 2, Exp: 6},
		ValidSince1:  1605323933,
		ValidSince2:  1605323952,
		ValidPeriod1: 268435455,
		ValidPeriod2: 268430000,
	}
	data := blockchain.BuildSettlement1ZkMsg(settlement1, testsample.PublicKeys[1], testsample.PublicKeys[2])

	var err error
	var testSuits = []CommitmentBuilderTest1{
		{
			TxData: settlement1.ToBytes(), ExpectedCommitmentInput: data,
			Account1: settlement1.Account1, Account2: settlement1.Account2,
			AccountPubKey1: testsample.PublicKeys[1], AccountPubKey2: testsample.PublicKeys[2],
		},
	}
	b, err := json.MarshalIndent(testSuits, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(commitmentBuilderTest1Output, b, 0644); err != nil {
		panic(err)
	}
}

func testCommitmentBuilder2() {
	settlement2 := &types.Settlement2{
		OpType:       types.SettlementOp21,
		AccountID2:   45,
		Amount2:      types.PackedAmount{Mantisa: 2, Exp: 16},
		Rate2:        types.PackedAmount{Mantisa: 5, Exp: 19},
		Fee2:         types.PackedFee{Mantisa: 2, Exp: 6},
		ValidSince2:  1605323952,
		ValidPeriod2: 268430000,
		LooID1:       34,
	}

	data := blockchain.BuildSettlement2ZkMsg(settlement2, 4, 5, testsample.PublicKeys[3])

	var err error
	var testSuits = []CommitmentBuilderTest2{
		{
			TxData:                  settlement2.ToBytes(),
			ExpectedCommitmentInput: data,
			AccountID:               settlement2.AccountID2,
			AccountPubKey:           testsample.PublicKeys[3],
			LooID:                   settlement2.LooID1,
			LooSrcToken:             4,
			LooDstToken:             5,
		},
	}
	b, err := json.MarshalIndent(testSuits, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(commitmentBuilderTest2Output, b, 0644); err != nil {
		panic(err)
	}
}

func testCommitmentBuilder3() {
	withdraw := &types.WithdrawOp{
		TokenID:    7,
		Amount:     types.PackedAmount{Mantisa: 314, Exp: 2},
		DestAddr:   common.HexToAddress("0x85E456C9AA9e8d6f1DF6E1aae6496b25b157634F"),
		AccountID:  13,
		ValidSince: 1607871567,
		Fee:        types.PackedFee{Mantisa: 5, Exp: 4},
	}

	data := blockchain.BuildWithdrawZkMsg(withdraw, testsample.PublicKeys[7])

	var err error
	var testSuits = []CommitmentBuilderTest3{
		{
			TxData:                  withdraw.ToBytes(),
			ExpectedCommitmentInput: data,
			AccountID:               withdraw.AccountID,
			AccountPubKey:           testsample.PublicKeys[7],
		},
	}
	b, err := json.MarshalIndent(testSuits, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(commitmentBuilderTest3Output, b, 0644); err != nil {
		panic(err)
	}
}

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
			Pubkey:  testsample.PublicKeys[3],
			Address: common.HexToAddress("0x052f46FeB45822E7f117536386C51B6Bd3125157"),
		},
	},
	AccountMax: 18,
	LooMax:     0,
}

func buildCommitmentFraudProofTest1() {
	bc := blockchain.NewBlockchain(genesis)
	genesisHash := bc.GetStateData().Hash()
	miniBlock1 := &types.MiniBlock{
		Txs: []types.Transaction{
			&types.Settlement1{
				OpType:   types.SettlementOp11,
				Token1:   1,
				Token2:   2,
				Account1: 8,
				Account2: 12,
				Rate1: types.PackedAmount{
					Mantisa: 1,
					Exp:     18,
				},
				Rate2: types.PackedAmount{
					Mantisa: 1,
					Exp:     18,
				},
				Amount1: types.PackedAmount{
					Mantisa: 2,
					Exp:     6,
				},
				Amount2: types.PackedAmount{
					Mantisa: 3,
					Exp:     6,
				},
				Fee1: types.PackedFee{
					Mantisa: 7,
					Exp:     3,
				},
				Fee2: types.PackedFee{
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
	bc.AddMiniBlock(miniBlock1)

	submitBlockStep := test.SubmitBlockStep{
		MiniBlocks:  []*types.MiniBlock{miniBlock1},
		Timestamp:   1600661872,
		BlockNumber: 1,
	}
	commitmentFPStep := test.AccuseCommitmentFraudProofStep{
		BlockNumber:      1,
		MiniBlockNumber:  0,
		MiniBlock:        miniBlock1,
		PostStateData:    bc.GetStateData(),
		MiniBlockProof:   proof.BuildMiniBlockProof(submitBlockStep.MiniBlocks, 0, submitBlockStep.Timestamp),
		CommitmentProofs: []hexutil.Bytes{bc.BuildCommitmentProof(miniBlock1)},
	}
	var testSuits = []*test.Suit{
		{
			Msg:              "test case when left over order at order 2",
			GenesisStateHash: genesisHash,
			Steps: []test.Step{
				{Action: test.SubmitBlock, Data: submitBlockStep},
				{Action: test.AccuseCommitmentFraudProof, Data: commitmentFPStep},
			},
		},
	}
	b, err := json.MarshalIndent(testSuits, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(commitmentFraudProofTest1Output, b, 0644); err != nil {
		panic(err)
	}
}

var genesis2 = &blockchain.Genesis{
	AccountAlloc: map[uint32]blockchain.GenesisAccount{
		0: {
			Tokens: map[uint16]*big.Int{
				0: big.NewInt(30000),
				1: big.NewInt(2000000),
			},
			Pubkey:  testsample.PublicKeys[2],
			Address: common.HexToAddress("0xdC70a72AbF352A0E3f75d737430EB896BA9Bf9Ea"),
		},

		17: {
			Tokens: map[uint16]*big.Int{
				0: big.NewInt(70000),
				1: big.NewInt(6000000),
			},
			Pubkey:  testsample.PublicKeys[0],
			Address: common.HexToAddress("0xdC70a72AbF352A0E3f75d737430EB896BA9Bf9Ea"),
		},
		123: {
			Tokens: map[uint16]*big.Int{
				0: big.NewInt(30000),
				1: big.NewInt(1000000),
				2: big.NewInt(5000000),
			},
			Pubkey:  testsample.PublicKeys[1],
			Address: common.HexToAddress("0x052f46FeB45822E7f117536386C51B6Bd3125157"),
		},
	},
	AccountMax: 1000,
	LooMax:     289,
	LooAlloc: map[uint64]*types.LeftOverOrder{
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

func buildCommitmentFraudProofTest2() {
	bc := blockchain.NewBlockchain(genesis2)
	genesisHash := bc.GetStateData().Hash()

	miniBlock1 := &types.MiniBlock{
		Txs: []types.Transaction{
			&types.Settlement2{
				OpType:     types.SettlementOp21,
				LooID1:     243,
				AccountID2: 123,
				Rate2: types.PackedAmount{
					Mantisa: 1,
					Exp:     18,
				},
				Amount2: types.PackedAmount{
					Mantisa: 3,
					Exp:     6,
				},
				Fee2: types.PackedFee{
					Mantisa: 4,
					Exp:     2,
				},
				ValidSince2:  1600661873,
				ValidPeriod2: 8640000,
			},
		},
	}
	bc.AddMiniBlock(miniBlock1)

	submitBlockStep := test.SubmitBlockStep{
		MiniBlocks:  []*types.MiniBlock{miniBlock1},
		Timestamp:   1600661874,
		BlockNumber: 1,
	}
	commitmentFPStep := test.AccuseCommitmentFraudProofStep{
		BlockNumber:      1,
		MiniBlockNumber:  0,
		MiniBlock:        miniBlock1,
		PostStateData:    bc.GetStateData(),
		MiniBlockProof:   proof.BuildMiniBlockProof(submitBlockStep.MiniBlocks, 0, submitBlockStep.Timestamp),
		CommitmentProofs: []hexutil.Bytes{bc.BuildCommitmentProof(miniBlock1)},
	}
	var testSuits = []*test.Suit{
		{
			Msg:              "test case when left over order at order 2",
			GenesisStateHash: genesisHash,
			Steps: []test.Step{
				{Action: test.SubmitBlock, Data: submitBlockStep},
				{Action: test.AccuseCommitmentFraudProof, Data: commitmentFPStep},
			},
		},
	}
	b, err := json.MarshalIndent(testSuits, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(commitmentFraudProofTest2Output, b, 0644); err != nil {
		panic(err)
	}
}

var genesis3 = &blockchain.Genesis{
	AccountAlloc: map[uint32]blockchain.GenesisAccount{
		0: {
			Tokens:  map[uint16]*big.Int{},
			Pubkey:  testsample.PublicKeys[1],
			Address: common.HexToAddress("0xdC70a72AbF352A0E3f75d737430EB896BA9Bf9Ea"),
		},

		23: {
			Tokens: map[uint16]*big.Int{
				0: big.NewInt(2000),
				2: big.NewInt(45242000),
			},
			Pubkey:  testsample.PublicKeys[2],
			Address: common.HexToAddress("0x99aF5AF1f1a61FE1678e030916f79331a28A57E8"),
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

func buildCommitmentFraudProofTest3() {
	bc := blockchain.NewBlockchain(genesis3)
	genesisHash := bc.GetStateData().Hash()
	// create an withdraw to another user
	withdraw := &types.WithdrawOp{
		TokenID:    2,
		Amount:     types.PackedAmount{Mantisa: 4, Exp: 7},
		DestAddr:   common.HexToAddress("0x99aF5AF1f1a61FE1678e030916f79331a28A57E8"),
		AccountID:  23,
		ValidSince: 0,
		Fee:        types.PackedFee{Mantisa: 1, Exp: 2},
	}
	miniBlock := &types.MiniBlock{
		Txs: []types.Transaction{
			withdraw,
		},
	}
	bc.AddMiniBlock(miniBlock)
	submitBlockStep := test.SubmitBlockStep{
		BlockNumber: 1,
		MiniBlocks:  []*types.MiniBlock{miniBlock},
		Timestamp:   1600661872,
	}

	commitmentFPStep := test.AccuseCommitmentFraudProofStep{
		BlockNumber:      1,
		MiniBlockNumber:  0,
		MiniBlock:        miniBlock,
		PostStateData:    bc.GetStateData(),
		MiniBlockProof:   proof.BuildMiniBlockProof(submitBlockStep.MiniBlocks, 0, submitBlockStep.Timestamp),
		CommitmentProofs: []hexutil.Bytes{bc.BuildCommitmentProof(miniBlock)},
	}
	var testSuits = []*test.Suit{
		{
			Msg:              "test commitment fraud proof for withdrawing",
			GenesisStateHash: genesisHash,
			Steps: []test.Step{
				{Action: test.SubmitBlock, Data: submitBlockStep},
				{Action: test.AccuseCommitmentFraudProof, Data: commitmentFPStep},
			},
		},
	}
	b, err := json.MarshalIndent(testSuits, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(commitmentFraudProofTest3Output, b, 0644); err != nil {
		panic(err)
	}
}

func main() {
	//testCommitmentBuilder1()
	//testCommitmentBuilder2()
	//testCommitmentBuilder3()
	//buildCommitmentFraudProofTest1()
	//buildCommitmentFraudProofTest2()
	buildCommitmentFraudProofTest3()
}
