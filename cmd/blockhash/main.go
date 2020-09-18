package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	util "github.com/KyberNetwork/l2-contract-test-suite/common"
)

const output = "testdata/submitBlock.json"

type SubmitBlockTestSuit struct {
	TimeStamp            uint32
	BlockNumber          uint32
	MiniBlocks           []hexutil.Bytes
	ExpectedNewBlockRoot hexutil.Bytes
}

type miniBlock struct {
	StateHash  common.Hash
	Commitment common.Hash
	Txs        [][]byte
}

func main() {
	var err error
	var testSuits []SubmitBlockTestSuit
	var miniBlockHashes []common.Hash
	for _, miniBlockLen := range []int{1} {
		// testSuit := MerkleTxsRootTestSuit{MiniBlockHashes: make([]common.Hash, miniBlockLen)}
		var miniBlocks []miniBlock
		var miniBlockDataArr []hexutil.Bytes

		for i := 0; i < miniBlockLen; i++ {
			miniBlockHash, miniBlockData, miniBlockStruct := generateMiniBlock()
			miniBlockHashes = append(miniBlockHashes, miniBlockHash)
			miniBlocks = append(miniBlocks, miniBlockStruct)
			miniBlockDataArr = append(miniBlockDataArr, miniBlockData)
		}

		blockInfoHash := util.GetMiniBlockHash(miniBlockHashes)
		prevBlockRoot := common.HexToHash("0x0")
		blockNumber := uint32(1)
		blockTime := uint32(1600237638)
		blockRoot := crypto.Keccak256(
			prevBlockRoot.Bytes(),
			blockInfoHash.Bytes(),
			util.Uint32ToBytes(blockTime),
			util.Uint32ToBytes(blockNumber),
			[]byte{util.Uint8ToByte(uint8(miniBlockLen))},
			miniBlocks[miniBlockLen-1].StateHash.Bytes(),
		)

		testSuit := SubmitBlockTestSuit{
			BlockNumber:          blockNumber,
			TimeStamp:            blockTime,
			MiniBlocks:           miniBlockDataArr,
			ExpectedNewBlockRoot: blockRoot,
		}
		testSuits = append(testSuits, testSuit)
	}

	b, err := json.MarshalIndent(testSuits, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(output, b, 0644); err != nil {
		panic(err)
	}
}

func generateMiniBlock() (common.Hash, []byte, miniBlock) {
	var stateHash, commitment common.Hash

	stateHash, err := util.GenerateRandomHash()
	if err != nil {
		panic(err)
	}
	commitment, err = util.GenerateRandomHash()
	if err != nil {
		panic(err)
	}

	var txs []byte
	for i := 0; i < 20; i++ {
		txs = append(txs, make([]byte, 6)...)
	}
	txRoot := crypto.Keccak256(txs)

	var miniBlockData []byte
	miniBlockData = append(miniBlockData, commitment.Bytes()...)
	miniBlockData = append(miniBlockData, stateHash.Bytes()...)
	miniBlockData = append(miniBlockData, txs...)

	miniBlockHash := crypto.Keccak256Hash(commitment.Bytes(), stateHash.Bytes(), txRoot)
	return miniBlockHash, miniBlockData, miniBlock{
		StateHash:  stateHash,
		Commitment: commitment,
	}
}
