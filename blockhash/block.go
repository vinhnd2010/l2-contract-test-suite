package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	randomMiniBlocks := make([][]byte, 55)
	miniBlockHash := getMiniBlockHash(randomMiniBlocks)

	fmt.Println(len(miniBlockHash))
	fmt.Println(miniBlockHash)

}

func getMiniBlockHash(miniBlocks [][]byte) [][]byte {
	fmt.Println("Round ", len(miniBlocks))
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
