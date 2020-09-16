package main

import (
	"crypto/rand"
	"encoding/json"
	"io/ioutil"
	// "time"
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

const output = "testdata/merkleTxsRoot.json"

type MerkleTxsRootTestSuit struct {
	MiniBlockHashes       []common.Hash
	ExpectedBlockInfoHash common.Hash
}

type SubmitBlockTestSuit struct {
	TimeStamp            uint32
	BlockNumber          uint32
	MiniBlocks           []miniBlock
	ExpectedNewBlockRoot hexutil.Bytes
}

type miniBlock struct {
	stateHash  common.Hash
	commitment common.Hash
	txs        []byte
}

func main() {
	var err error
	var testSuits []SubmitBlockTestSuit
	var miniBlockHashes []common.Hash
	for _, miniBlockLen := range []int{1, 2, 3, 4, 5} {
		// testSuit := MerkleTxsRootTestSuit{MiniBlockHashes: make([]common.Hash, miniBlockLen)}
		var miniBlocks []miniBlock

		for i := 0; i < miniBlockLen; i++ {
			miniBlockHash, miniBlock := generateMiniBlock()
			miniBlockHashes = append(miniBlockHashes, miniBlockHash)
			miniBlocks = append(miniBlocks, miniBlock)
		}

		blockInfoHash := getMiniBlockHash(miniBlockHashes)[0]
		prevBlockRoot := common.HexToHash("0x0")
		blockNumber := uint32(1)
		blockTime := uint32(1600237638)
		blockRoot := crypto.Keccak256(prevBlockRoot.Bytes(), blockInfoHash.Bytes(), uint32ToBytes(blockTime),
			uint32ToBytes(blockNumber), uint32ToBytes(uint32(miniBlockLen)), miniBlocks[miniBlockLen-1].stateHash.Bytes())
		testSuit := SubmitBlockTestSuit{
			BlockNumber:          blockNumber,
			TimeStamp:            blockTime,
			MiniBlocks:           miniBlocks,
			ExpectedNewBlockRoot: blockRoot,
		}

		testSuits = append(testSuits, testSuit)
	}

	b, err := json.Marshal(testSuits)
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(output, b, 0644); err != nil {
		panic(err)
	}
}

func uint32ToBytes(number uint32) []byte {
	out := make([]byte, 4)
	binary.LittleEndian.PutUint32(out, number)
	return out
}

func generateRandomHash() (common.Hash, error) {
	var out common.Hash
	_, err := rand.Read(out[:])
	return out, err
}

func getMiniBlockHash(miniBlocks []common.Hash) []common.Hash {
	if len(miniBlocks) == 1 {
		return miniBlocks
	}
	var newMiniBlocks []common.Hash
	for i := 0; i < len(miniBlocks); i += 2 {
		var miniBlock common.Hash
		if i+1 == len(miniBlocks) {
			miniBlock = crypto.Keccak256Hash(miniBlocks[i].Bytes(), common.HexToHash("0x0").Bytes())
		} else {
			miniBlock = crypto.Keccak256Hash(miniBlocks[i].Bytes(), miniBlocks[i+1].Bytes())
		}
		newMiniBlocks = append(newMiniBlocks, miniBlock)
	}
	return getMiniBlockHash(newMiniBlocks)
}

func generateMiniBlock() (common.Hash, miniBlock) {
	var txs [][]byte
	for i := 0; i < 20; i++ {
		txs = append(txs, make([]byte, 6))
	}
	txRoot := crypto.Keccak256(txs...)

	var stateHash, commitment common.Hash

	stateHash, err := generateRandomHash()
	if err != nil {
		panic(err)
	}
	commitment, err = generateRandomHash()
	if err != nil {
		panic(err)
	}

	miniBlockHash := crypto.Keccak256Hash(commitment.Bytes(), stateHash.Bytes(), txRoot)
	return miniBlockHash, miniBlock{
		stateHash:  stateHash,
		commitment: commitment,
		txs:        txRoot,
	}
}
