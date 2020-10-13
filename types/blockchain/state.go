package blockchain

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	util "github.com/KyberNetwork/l2-contract-test-suite/common"
)

type State struct {
	accounts map[uint32]*Account
	tree     *MerkleTree
}

func NewState() *State {
	return &State{
		accounts: make(map[uint32]*Account),
		tree:     NewTree(StateTreeDeep),
	}
}

type StateData struct {
	StateRoot  common.Hash
	LOORoot    common.Hash
	AccountMax uint32
	LOOMax     uint64
}

func (sData *StateData) Hash() common.Hash {
	return crypto.Keccak256Hash(
		sData.StateRoot.Bytes(), sData.LOORoot.Bytes(),
		util.Uint32ToBytes(sData.AccountMax), util.Uint48ToBytes(sData.LOOMax),
	)
}

func NewStateFromAlloc(acountAlloc map[uint32]GenesisAccount) *State{
	var state = NewState();
	for accountID, accountAlloc := range acountAlloc {
		account := NewAccount(accountAlloc.Pubkey, accountAlloc.Address)
		for tokenID, tokenAmount := range accountAlloc.Tokens {
			account.Update(tokenID, tokenAmount)
		}
		state.accounts[accountID] = account
		accountHash := crypto.Keccak256Hash(
			account.tree.RootHash().Bytes(),
			account.GetPubAccountHash().Bytes(),
		)
		state.tree.Update(uint64(accountID), accountHash)
	}
	return state
}
