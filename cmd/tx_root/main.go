package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/common"

	util "github.com/KyberNetwork/l2-contract-test-suite/common"
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
			if testSuit.MiniBlockHashes[i], err = util.GenerateRandomHash(); err != nil {
				panic(err)
			}
		}
		testSuit.ExpectedBlockInfoHash = util.GetMiniBlockHash(testSuit.MiniBlockHashes)[0]
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
