package main

import (
	"crypto/rand"
	"encoding/json"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const output = "testdata/merkleTxsRoot.json"

type MerkleTxsRootTestSuit struct {
	MiniBlockHashes       []common.Hash
	ExpectedBlockInfoHash common.Hash
}

func main() {
	var err error
	var testSuits []MerkleTxsRootTestSuit
	for _, miniBlockLen := range []int{1, 2, 3, 4, 5} {
		testSuit := MerkleTxsRootTestSuit{MiniBlockHashes: make([]common.Hash, miniBlockLen)}
		for i := 0; i < miniBlockLen; i++ {
			if testSuit.MiniBlockHashes[i], err = generateRandomHash(); err != nil {
				panic(err)
			}
		}
		testSuit.ExpectedBlockInfoHash = getMiniBlockHash(testSuit.MiniBlockHashes)[0]
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
