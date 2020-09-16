package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/crypto"
)

type BlockInfo struct {
	Timestamp   int
	BlockNumber int
	K           int
	StateHashK  string
}

func main() {
	var miniBlocks [][]byte
	for i := 0; i <= 100; i++ {
		miniBlocks = append(miniBlocks, generateMiniBlock())
	}
	blockInfoHash := getMiniBlockHash(miniBlocks)

	fmt.Println(blockInfoHash)

	generateJson()
}

func getMiniBlockHash(miniBlocks [][]byte) [][]byte {
	var newMiniBlocks [][]byte

	if len(miniBlocks) == 1 {
		return miniBlocks
	}

	for i := 0; i < len(miniBlocks)-1; i += 2 {
		var miniBlock []byte
		if i+1 == len(miniBlocks) {
			miniBlock = crypto.Keccak256(miniBlocks[i], miniBlocks[i+1])
		} else {
			miniBlock = crypto.Keccak256(miniBlocks[i], []byte{0})
		}
		newMiniBlocks = append(newMiniBlocks, miniBlock)
	}

	return getMiniBlockHash(newMiniBlocks)
}

func generateMiniBlock() []byte {
	var txs [][]byte
	for i := 0; i <= 20; i++ {
		txs = append(txs, make([]byte, 6))
	}
	txRoot := crypto.Keccak256(txs...)

	stateHash := make([]byte, 32)
	commitment := make([]byte, 32)

	miniBlock := crypto.Keccak256(commitment, stateHash, txRoot)
	return miniBlock
}

func generateJson() {
	data := []BlockInfo{
		BlockInfo{
			Timestamp:   1600226307,
			BlockNumber: 10870686,
			K:           1,
			StateHashK:  hex.EncodeToString(make([]byte, 32)),
		},
		BlockInfo{
			Timestamp:   1600226835,
			BlockNumber: 10870720,
			K:           2,
			StateHashK:  hex.EncodeToString(make([]byte, 32)),
		},
	}

	file, _ := json.MarshalIndent(data, "", "")
	_ = ioutil.WriteFile("state.json", file, 0644)
}
