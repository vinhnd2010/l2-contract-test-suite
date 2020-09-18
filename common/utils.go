package common

import (
	"crypto/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func GenerateRandomHash() (common.Hash, error) {
	var out common.Hash
	_, err := rand.Read(out[:])
	return out, err
}

func GetMiniBlockHash(miniBlocks []common.Hash) []common.Hash {
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
	return GetMiniBlockHash(newMiniBlocks)
}
