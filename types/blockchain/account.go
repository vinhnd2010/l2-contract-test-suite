package blockchain

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

type Account struct {
	pubKey     hexutil.Bytes
	withdrawTo common.Address
	tree       *MerkleTree
}

func NewAccount(pubKey hexutil.Bytes, withdrawTo common.Address) *Account {
	return &Account{
		pubKey:     pubKey,
		withdrawTo: withdrawTo,
		tree:       NewTree(AccountTreeDeep),
	}
}

// update tree, returns a new Hash
func (a *Account) Update(tokenID uint16, amount *big.Int) common.Hash {
	a.tree.Update(uint64(tokenID), common.BigToHash(amount))
	return a.tree.RootHash()
}

func (a *Account) GetPubAccountHash() common.Hash {
	return crypto.Keccak256Hash(a.pubKey, a.withdrawTo.Bytes())
}


