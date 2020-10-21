package main

import (
	"encoding/json"
	"io/ioutil"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/KyberNetwork/l2-contract-test-suite/common/proof"
	"github.com/KyberNetwork/l2-contract-test-suite/types"
	"github.com/KyberNetwork/l2-contract-test-suite/types/blockchain"
)

const output = "testdata/fraudProofSettlement1.json"
const benchmarkOutput = "benchmarkdata/fraudProofSettlement1.json"

type FraudProofTestSuit struct {
	Msg              string
	GenesisStateHash common.Hash
	Blocks           []BlockData
}
type BlockData struct {
	MiniBlocks      []*types.MiniBlock
	Timestamp       uint32
	MiniBlockNumber int
	Proof           *FraudProof
}

type FraudProof struct {
	PrevStateData      *blockchain.StateData
	PrevStateHashProof hexutil.Bytes
	MiniBlockProof     hexutil.Bytes
	ExecutionProof     []hexutil.Bytes
}

var (
	pubKey1, _ = hexutil.Decode("0xb8748a745b1c75a34238d56576e41bea9207fb5e1f7da8abe741bd9dbf14dd0e0cfb7e0cf1380065477345a42aa821aa1c68e7d9eb213eee1e8f00cb707458a4")
	pubKey2, _ = hexutil.Decode("0xe61f3aab7e1bd78495524c955a6e3f89152ee3811fe52b85882002c465a235f7dc9bc9ed7b58277d5f9036c85e47958c65bc81104718a9364a294d96b4d277da")
	pubKey3, _ = hexutil.Decode("0x5bb440955b11980eaad949aa3f1fb05825c53cefb211b0f515415107a3aaf9dec1820b7899ad2a62a1c4aacf320b1a528c8c98aa558ee777e60110be62626e42")
)

var genesis = &blockchain.Genesis{
	AccountAlloc: map[uint32]blockchain.GenesisAccount{
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
}

func buildTest1() *FraudProofTestSuit {
	bc := blockchain.NewBlockchain(genesis)
	genesisHash := bc.GetStateData().Hash()

	preStateData := bc.GetStateData()

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

	executionProofs := bc.AddMiniBlock(miniBlock1)
	blockData := BlockData{
		MiniBlocks:      []*types.MiniBlock{miniBlock1},
		Timestamp:       1600661872,
		MiniBlockNumber: 0,
		Proof: &FraudProof{
			PrevStateData:      preStateData,
			PrevStateHashProof: []byte{},
			ExecutionProof:     executionProofs,
		},
	}
	blockData.Proof.MiniBlockProof = proof.BuildMiniBlockProof(blockData.MiniBlocks, uint(blockData.MiniBlockNumber), blockData.Timestamp)
	return &FraudProofTestSuit{
		Msg:              "test case when left over order at order 2",
		GenesisStateHash: genesisHash,
		Blocks: []BlockData{
			blockData,
		},
	}
}

func buildTestForSecondBlock() *FraudProofTestSuit {
	bc := blockchain.NewBlockchain(genesis)
	genesisHash := bc.GetStateData().Hash()
	preStateData := bc.GetStateData()

	miniBlock1 := &types.MiniBlock{Txs: nil}
	bc.AddMiniBlock(miniBlock1)

	blockData1 := BlockData{
		MiniBlocks:      []*types.MiniBlock{miniBlock1},
		Timestamp:       1600661870,
		MiniBlockNumber: -1,
	}

	miniBlock2 := &types.MiniBlock{
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
	executionProofs := bc.AddMiniBlock(miniBlock2)
	blockData2 := BlockData{
		MiniBlocks:      []*types.MiniBlock{miniBlock2},
		Timestamp:       1600661872,
		MiniBlockNumber: 0,
		Proof: &FraudProof{
			PrevStateData:      preStateData,
			PrevStateHashProof: proof.BuildFinalStateHashProof(blockData1.MiniBlocks, blockData1.Timestamp),
			ExecutionProof:     executionProofs,
		},
	}
	blockData2.Proof.MiniBlockProof = proof.BuildMiniBlockProof(blockData2.MiniBlocks, uint(blockData2.MiniBlockNumber), blockData2.Timestamp)
	return &FraudProofTestSuit{
		Msg:              "test case when miniBlock is at block number = 2",
		GenesisStateHash: genesisHash,
		Blocks:           []BlockData{blockData1, blockData2},
	}
}

func buildTestForSecondMiniBlock() *FraudProofTestSuit {
	bc := blockchain.NewBlockchain(genesis)
	genesisHash := bc.GetStateData().Hash()

	miniBlock1 := &types.MiniBlock{Txs: nil}
	bc.AddMiniBlock(miniBlock1)

	preStateData := bc.GetStateData()
	miniBlock2 := &types.MiniBlock{
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
	executionProofs := bc.AddMiniBlock(miniBlock2)

	blockData := BlockData{
		MiniBlocks:      []*types.MiniBlock{miniBlock1, miniBlock2},
		Timestamp:       1600661872,
		MiniBlockNumber: 1,
		Proof: &FraudProof{
			PrevStateData:  preStateData,
			ExecutionProof: executionProofs,
		},
	}
	blockData.Proof.PrevStateHashProof = proof.BuildPrevStateHashMiniBlockProof(blockData.MiniBlocks, uint(blockData.MiniBlockNumber-1))
	blockData.Proof.MiniBlockProof = proof.BuildMiniBlockProof(blockData.MiniBlocks, uint(blockData.MiniBlockNumber), blockData.Timestamp)
	return &FraudProofTestSuit{
		Msg:              "test case when mimiBlockNumber = 2",
		GenesisStateHash: genesisHash,
		Blocks: []BlockData{
			blockData,
		},
	}
}

func buildTest2() *FraudProofTestSuit {
	bc := blockchain.NewBlockchain(genesis)
	genesisHash := bc.GetStateData().Hash()
	preStateData := bc.GetStateData()

	var txs []types.Transaction
	for i := 0; i < 15; i++ {
		txs = append(txs, &types.Settlement1{
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
				Exp:     3,
			},
			Amount2: types.PackedAmount{
				Mantisa: 3,
				Exp:     3,
			},
			Fee1: types.PackedFee{
				Mantisa: 1,
				Exp:     2,
			},
			Fee2: types.PackedFee{
				Mantisa: 1,
				Exp:     2,
			},
			ValidSince1:  1600661872,
			ValidSince2:  1600661873,
			ValidPeriod1: 86400,
			ValidPeriod2: 86400,
		})
	}

	miniBlock1 := &types.MiniBlock{
		Txs: txs,
	}

	executionProofs := bc.AddMiniBlock(miniBlock1)
	blockData := BlockData{
		MiniBlocks:      []*types.MiniBlock{miniBlock1},
		Timestamp:       1600661872,
		MiniBlockNumber: 0,
		Proof: &FraudProof{
			PrevStateData:      preStateData,
			PrevStateHashProof: []byte{},
			ExecutionProof:     executionProofs,
		},
	}
	blockData.Proof.MiniBlockProof = proof.BuildMiniBlockProof(blockData.MiniBlocks, uint(blockData.MiniBlockNumber), blockData.Timestamp)
	return &FraudProofTestSuit{
		Msg:              "test case with 15 orders",
		GenesisStateHash: genesisHash,
		Blocks: []BlockData{
			blockData,
		},
	}
}

func main() {
	var testSuits []*FraudProofTestSuit
	testSuits = append(testSuits, buildTest1())
	testSuits = append(testSuits, buildTestForSecondBlock())
	testSuits = append(testSuits, buildTestForSecondMiniBlock())

	b, err := json.MarshalIndent(testSuits, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(output, b, 0644); err != nil {
		panic(err)
	}

	var testSuits2 []*FraudProofTestSuit
	testSuits2 = append(testSuits2, buildTest2())
	b, err = json.MarshalIndent(testSuits2, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(benchmarkOutput, b, 0644); err != nil {
		panic(err)
	}

}
