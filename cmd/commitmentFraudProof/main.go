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

const commitmentBuilderTest1Output = "testdata/commitmentBuilder1.json"
const commitmentFraudProofTest1Output = "testdata/commitmentFraudProof1.json"

type CommitmentBuilderTest1 struct {
	TxData                  hexutil.Bytes `json:"txData"`
	ExpectedCommitmentInput hexutil.Bytes `json:"commitmentInput"`
	Account1                uint32        `json:"accountID1"`
	Account2                uint32        `json:"accountID2"`
	AccountPubKey1          hexutil.Bytes `json:"accountPubKey1"`
	AccountPubKey2          hexutil.Bytes `json:"accountPubKey2"`
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

func buildTest1() *test.Suit {
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
	return &test.Suit{
		Msg:              "test case when left over order at order 2",
		GenesisStateHash: genesisHash,
		Steps: []test.Step{
			{Action: test.SubmitBlock, Data: submitBlockStep},
			{Action: test.AccuseCommitmentFraudProof, Data: commitmentFPStep},
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
	if err := ioutil.WriteFile(commitmentFraudProofTest1Output, b, 0644); err != nil {
		panic(err)
	}
}
