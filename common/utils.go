package common

import (
	"crypto/rand"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	Precision = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
)

func GenerateRandomHash() (common.Hash, error) {
	var out common.Hash
	_, err := rand.Read(out[:])
	return out, err
}

func GenerateRandomPubKey() hexutil.Bytes {
	out := make(hexutil.Bytes, 64)
	_, err := rand.Read(out[:])
	if err != nil {
		panic(err)
	}
	return out
}

func AddAmount(beforeValue common.Hash, value *big.Int) common.Hash {
	return common.BigToHash(new(big.Int).Add(beforeValue.Big(), value))
}

func SubAmount(beforeValue common.Hash, value *big.Int) common.Hash {
	if beforeValue.Big().Cmp(value) < 0 {
		panic("insufficient funds")
	}
	return common.BigToHash(new(big.Int).Sub(beforeValue.Big(), value))
}

func CalAmountOut(amount *big.Int, rate *big.Int) *big.Int {
	return new(big.Int).Div(new(big.Int).Mul(amount, rate), Precision)
}

func CalAmountIn(amount *big.Int, rate *big.Int) *big.Int {
	tmp := new(big.Int).Mul(amount, Precision)
	tmp = new(big.Int).Add(tmp, rate)
	tmp = new(big.Int).Sub(tmp, big.NewInt(1))
	return tmp.Div(tmp, rate)
}

func GetMiniBlockHash(miniBlocks []common.Hash) common.Hash {
	if len(miniBlocks) == 1 {
		return miniBlocks[0]
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
