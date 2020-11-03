package main

import (
	"encoding/json"
	"io/ioutil"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/KyberNetwork/l2-contract-test-suite/testsample"
	"github.com/KyberNetwork/l2-contract-test-suite/types"
	"github.com/KyberNetwork/l2-contract-test-suite/types/blockchain"
	"github.com/KyberNetwork/l2-contract-test-suite/types/test"
)

const testOutput = "simulateData/test1.json"

var genesis = &blockchain.Genesis{
	AccountAlloc: map[uint32]blockchain.GenesisAccount{
		0: {
			Tokens: map[uint16]*big.Int{
				0: big.NewInt(30000),
				1: big.NewInt(2000000),
			},
			Pubkey:  testsample.PublicKeys[0],
			Address: common.HexToAddress("0x91F4d9EA5c1ee0fc778524b3D57fD8CF700996Cf"),
		},
	},
	AccountMax: 0,
	LooMax:     0,
}

func buildTest1() *test.Suit {
	bc := blockchain.NewBlockchain(genesis)
	genesisHash := bc.GetStateData().Hash()

	var steps []test.Step

	var depositToNewTxs []types.Transaction
	for i, account := range testsample.Accounts {
		if i == 0 { // skip accounts[0]
			continue
		}
		deposit := &types.DepositToNewOp{
			PubKey:     testsample.PublicKeys[i],
			WithdrawTo: account,
			TokenID:    0,
			Amount:     types.PackedAmount{Mantisa: 3, Exp: 8}.Big(),
		}
		depositToNewTxs = append(depositToNewTxs, deposit)
		steps = append(steps, test.Step{Action: test.SubmitDepositToNew, Data: deposit})
	}
	miniBlock1 := &types.MiniBlock{
		Txs: depositToNewTxs,
	}
	bc.AddMiniBlock(miniBlock1)
	steps = append(steps, test.Step{Action: test.SubmitBlock, Data: test.SubmitBlockStep{
		BlockNumber: 1,
		MiniBlocks:  []*types.MiniBlock{miniBlock1},
		Timestamp:   1600661872,
	}})
	// deposit steps
	{
		var depositMiniBlocks []*types.MiniBlock
		for i := range testsample.Accounts {
			var deposits []types.Transaction
			for j := 1; j < 6; j++ {
				deposit := &types.DepositOp{
					AccountID: uint32(i),
					TokenID:   uint16(j),
					Amount:    types.PackedAmount{Mantisa: 5, Exp: 10}.Big(),
				}
				deposits = append(deposits, deposit)
				steps = append(steps, test.Step{Action: test.SubmitDeposit, Data: deposit})
			}

			miniBlock := &types.MiniBlock{
				Txs: deposits,
			}
			bc.AddMiniBlock(miniBlock)
			depositMiniBlocks = append(depositMiniBlocks, miniBlock)
		}
		steps = append(steps, test.Step{Action: test.SubmitBlock, Data: test.SubmitBlockStep{
			BlockNumber: 2,
			MiniBlocks:  depositMiniBlocks,
			Timestamp:   1600661873,
		}})
	}
	// add settlement block
	{
		var settlement1Block = &types.MiniBlock{
			Txs: []types.Transaction{
				&types.Settlement1{
					OpType:       types.SettlementOp11,
					Token1:       1,
					Token2:       2,
					Account1:     4,
					Account2:     5,
					Rate1:        types.PackedAmount{Mantisa: 2, Exp: 18},
					Rate2:        types.PackedAmount{Mantisa: 5, Exp: 17},
					Amount1:      types.PackedAmount{Mantisa: 2, Exp: 6},
					Amount2:      types.PackedAmount{Mantisa: 2, Exp: 7},
					Fee1:         types.PackedFee{Mantisa: 3, Exp: 2},
					Fee2:         types.PackedFee{Mantisa: 3, Exp: 3},
					ValidSince1:  1600661873,
					ValidSince2:  1600661873,
					ValidPeriod1: 6400,
					ValidPeriod2: 86400,
				},
			},
		}
		bc.AddMiniBlock(settlement1Block)

		var settlement2Block = &types.MiniBlock{
			Txs: []types.Transaction{
				&types.Settlement1{
					OpType:       types.SettlementOp11,
					Token1:       3,
					Token2:       4,
					Account1:     4,
					Account2:     5,
					Rate1:        types.PackedAmount{Mantisa: 2, Exp: 18},
					Rate2:        types.PackedAmount{Mantisa: 5, Exp: 17},
					Amount1:      types.PackedAmount{Mantisa: 2, Exp: 6},
					Amount2:      types.PackedAmount{Mantisa: 2, Exp: 7},
					Fee1:         types.PackedFee{Mantisa: 3, Exp: 2},
					Fee2:         types.PackedFee{Mantisa: 3, Exp: 3},
					ValidSince1:  1600661873,
					ValidSince2:  1600661873,
					ValidPeriod1: 6400,
					ValidPeriod2: 86400,
				},
				&types.Settlement2{
					OpType:       types.SettlementOp21,
					LooID1:       2,
					AccountID2:   6,
					Amount2:      types.PackedAmount{Mantisa: 8, Exp: 6},
					Rate2:        types.PackedAmount{Mantisa: 2, Exp: 18},
					Fee2:         types.PackedFee{},
					ValidSince2:  1600661874,
					ValidPeriod2: 86400,
				},
			},
		}
		bc.AddMiniBlock(settlement2Block)

		var settlement3Block = &types.MiniBlock{
			Txs: []types.Transaction{
				&types.Settlement1{ // create loo with 16 e6 token 0
					OpType:       types.SettlementOp11,
					Token1:       5,
					Token2:       0,
					Account1:     6,
					Account2:     7,
					Rate1:        types.PackedAmount{Mantisa: 2, Exp: 18},
					Rate2:        types.PackedAmount{Mantisa: 5, Exp: 17},
					Amount1:      types.PackedAmount{Mantisa: 2, Exp: 6},
					Amount2:      types.PackedAmount{Mantisa: 2, Exp: 7},
					Fee1:         types.PackedFee{Mantisa: 3, Exp: 2},
					Fee2:         types.PackedFee{Mantisa: 3, Exp: 3},
					ValidSince1:  1600661875,
					ValidSince2:  1600661875,
					ValidPeriod1: 6400,
					ValidPeriod2: 86400,
				},
				&types.Settlement1{ // create loo with 1 e6 token 5
					OpType:       types.SettlementOp11,
					Token1:       5,
					Token2:       0,
					Account1:     1,
					Account2:     2,
					Rate1:        types.PackedAmount{Mantisa: 2, Exp: 18},
					Rate2:        types.PackedAmount{Mantisa: 5, Exp: 17},
					Amount1:      types.PackedAmount{Mantisa: 2, Exp: 6},
					Amount2:      types.PackedAmount{Mantisa: 2, Exp: 6},
					Fee1:         types.PackedFee{Mantisa: 3, Exp: 2},
					Fee2:         types.PackedFee{Mantisa: 3, Exp: 3},
					ValidSince1:  1600661875,
					ValidSince2:  1600661875,
					ValidPeriod1: 6400,
					ValidPeriod2: 86400,
				},
				&types.Settlement3{
					LooID1: 3,
					LooID2: 4,
				},
			},
		}
		bc.AddMiniBlock(settlement3Block)

		steps = append(steps, test.Step{Action: test.SubmitBlock, Data: test.SubmitBlockStep{
			BlockNumber: 3,
			MiniBlocks:  []*types.MiniBlock{settlement1Block, settlement2Block, settlement3Block},
			Timestamp:   1600661890,
		}})
	}

	return &test.Suit{
		Msg:              "create random data set of block combination",
		GenesisStateHash: genesisHash,
		Steps:            steps,
	}
}

func main() {
	var testSuits *test.Suit = buildTest1()
	b, err := json.MarshalIndent(testSuits, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(testOutput, b, 0644); err != nil {
		panic(err)
	}

}
