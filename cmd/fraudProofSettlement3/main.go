package main

import (
	"encoding/json"
	"io/ioutil"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	util "github.com/KyberNetwork/l2-contract-test-suite/common"
	"github.com/KyberNetwork/l2-contract-test-suite/common/proof"
	"github.com/KyberNetwork/l2-contract-test-suite/testsample"
	"github.com/KyberNetwork/l2-contract-test-suite/types"
	"github.com/KyberNetwork/l2-contract-test-suite/types/blockchain"
)

const testOutput = "testdata/fraudProofSettlement3.json"
const benchmarkOutput = "benchmarkdata/fraudProofSettlement3.json"

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

		17: {
			Tokens: map[uint16]*big.Int{
				0: big.NewInt(70000),
				1: big.NewInt(6000000),
			},
			Pubkey:  pubKey2,
			Address: common.HexToAddress("0xdC70a72AbF352A0E3f75d737430EB896BA9Bf9Ea"),
		},
		30: {
			Tokens: map[uint16]*big.Int{
				0: big.NewInt(30000),
				1: big.NewInt(1000000),
				2: big.NewInt(5000000),
			},
			Pubkey:  pubKey3,
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

func buildTest1() *FraudProofTestSuit {
	bc := blockchain.NewBlockchain(genesis)
	genesisHash := bc.GetStateData().Hash()
	preStateData := bc.GetStateData()

	miniBlock1 := &types.MiniBlock{
		Txs: []types.Transaction{
			&types.Settlement3{
				LooID1: 243,
				LooID2: 56,
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
		Msg:              "test case when looID1 is fully filled and loo2 is partially filled",
		GenesisStateHash: genesisHash,
		Blocks: []BlockData{
			blockData,
		},
	}
}

func buildBenchmarkTest() *FraudProofTestSuit {
	var benchmarkGenesis = &blockchain.Genesis{
		AccountAlloc: map[uint32]blockchain.GenesisAccount{
			0: blockchain.GenesisAccount{
				Tokens:  map[uint16]*big.Int{},
				Pubkey:  testsample.PublicKeys[1],
				Address: common.HexToAddress("0x9aab3f75489902f3a48495025729a0af77d4b11e"),
			},

			1: blockchain.GenesisAccount{
				Tokens: map[uint16]*big.Int{
					0: types.PackedAmount{Mantisa: 2, Exp: 18}.Big(),
					1: types.PackedAmount{Mantisa: 2, Exp: 18}.Big(),
				},
				Pubkey:  testsample.PublicKeys[3],
				Address: common.HexToAddress("0x85E456C9AA9e8d6f1DF6E1aae6496b25b157634F"),
			},
			2: blockchain.GenesisAccount{
				Tokens: map[uint16]*big.Int{
					0: types.PackedAmount{Mantisa: 2, Exp: 18}.Big(),
					2: types.PackedAmount{Mantisa: 2, Exp: 18}.Big(),
				},
				Pubkey:  testsample.PublicKeys[3],
				Address: common.HexToAddress("0x85E456C9AA9e8d6f1DF6E1aae6496b25b157634F"),
			},
		},
		AccountMax: 8,
		LooAlloc:   make(map[uint64]*types.LeftOverOrder),
		LooMax:     10000,
	}

	for i := 0; i < 15; i++ {
		benchmarkGenesis.LooAlloc[uint64(i*2)] = &types.LeftOverOrder{
			AccountID:   1,
			SrcToken:    1,
			DestToken:   2,
			Amount:      big.NewInt(1),
			Fee:         big.NewInt(1),
			Rate:        util.Precision,
			ValidSince:  1601868254,
			ValidPeriod: 86400,
		}
		benchmarkGenesis.LooAlloc[uint64(i*2+1)] = &types.LeftOverOrder{
			AccountID:   2,
			SrcToken:    2,
			DestToken:   1,
			Amount:      big.NewInt(1),
			Fee:         big.NewInt(1),
			Rate:        util.Precision,
			ValidSince:  1601868254,
			ValidPeriod: 86400,
		}
	}

	bc := blockchain.NewBlockchain(benchmarkGenesis)
	genesisHash := bc.GetStateData().Hash()
	preStateData := bc.GetStateData()

	miniBlock1 := &types.MiniBlock{
		Txs: []types.Transaction{},
	}

	for i := 0; i < 15; i++ {
		miniBlock1.Txs = append(miniBlock1.Txs, &types.Settlement3{
			LooID1: uint64(i * 2),
			LooID2: uint64(i*2 + 1),
		})
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
		Msg:              "benchmark to settlement 2 loo order",
		GenesisStateHash: genesisHash,
		Blocks: []BlockData{
			blockData,
		},
	}

}

func main() {
	var testSuits []*FraudProofTestSuit
	testSuits = append(testSuits, buildTest1())
	b, err := json.MarshalIndent(testSuits, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(testOutput, b, 0644); err != nil {
		panic(err)
	}

	var testSuits2 []*FraudProofTestSuit
	testSuits2 = append(testSuits2, buildBenchmarkTest())
	b, err = json.MarshalIndent(testSuits2, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(benchmarkOutput, b, 0644); err != nil {
		panic(err)
	}
}
