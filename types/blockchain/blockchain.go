package blockchain

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	util "github.com/KyberNetwork/l2-contract-test-suite/common"
	"github.com/KyberNetwork/l2-contract-test-suite/types"
)

const AccountTreeDeep = 11
const StateTreeDeep = 33
const LOOTreeDeep = 45
const FeeTokenIndex = 0
const AdminIndex uint32 = 0

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
	return a.tree.rootHash()
}

func (a *Account) GetPubAccountHash() common.Hash {
	return crypto.Keccak256Hash(a.pubKey, a.withdrawTo.Bytes())
}

type Blockchain struct {
	state      *State
	accountMax uint32
	looState   *LeftOverOrderList
	looMax     uint64
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

//TODO: add loo to genesis later
type Genesis struct {
	AccountAlloc map[uint32]GenesisAccount
	AccountMax   uint32

	LooMax uint64
}

type GenesisAccount struct {
	Tokens  map[uint16]*big.Int
	Pubkey  hexutil.Bytes
	Address common.Address
}

func NewBlockchain(genesis *Genesis) *Blockchain {
	if genesis == nil {
		return &Blockchain{
			state:      NewState(),
			accountMax: 0,
			looState:   NewLOOList(),
			looMax:     0,
		}
	}

	state := NewState()
	for accountID, accountAlloc := range genesis.AccountAlloc {
		account := NewAccount(accountAlloc.Pubkey, accountAlloc.Address)
		for tokenID, tokenAmount := range accountAlloc.Tokens {
			account.Update(tokenID, tokenAmount)
		}
		state.accounts[accountID] = account
		accountHash := crypto.Keccak256Hash(
			account.tree.rootHash().Bytes(),
			account.GetPubAccountHash().Bytes(),
		)
		state.tree.Update(uint64(accountID), accountHash)
	}
	return &Blockchain{
		state:      state,
		accountMax: genesis.AccountMax,
		looState:   NewLOOList(),
		looMax:     genesis.LooMax,
	}
}

//func (bc *Blockchain) AddBlock(block *types.Blo)

func (bc *Blockchain) AddMiniBlock(block *types.MiniBlock) []hexutil.Bytes {
	var (
		proofs   []hexutil.Bytes
		totalFee = big.NewInt(0)
	)

	for _, tx := range block.Txs {
		switch obj := tx.(type) {
		case *types.Settlement1:
			proof, fee := bc.handleSettlement1(obj)
			proofs = append(proofs, proof)
			totalFee = totalFee.Add(totalFee, fee)
		default:
			panic("unsupported type")
		}
	}
	proofs = append(proofs, bc.handleTotalFee(totalFee))
	block.StateHash = bc.GetStateData().Hash()
	//TODO: build commitment here
	block.Commitment = common.HexToHash(zeroHash)
	return proofs
}

func (bc *Blockchain) handleSettlement1(op *types.Settlement1) (proof hexutil.Bytes, fee *big.Int) {
	account := bc.state.accounts[op.Account1]
	if account == nil {
		panic("empty account")
	}

	amount1, amount2, fee1, fee2, loo := op.GetSettlementValue()
	_, accountSiblings := bc.state.tree.GetProof(uint64(op.Account1))
	proof = appendSiblings(proof, accountSiblings)
	pubAccountHash := account.GetPubAccountHash()
	proof = append(proof, pubAccountHash.Bytes()...)
	// update balance of token
	token1Amount, token1Siblings := account.tree.GetProof(uint64(op.Token1))
	account.tree.Update(uint64(op.Token1), util.SubAmount(token1Amount, amount1))
	proof = appendTokenProof(proof, token1Amount, token1Siblings)

	token2Amount, token2Siblings := account.tree.GetProof(uint64(op.Token2))
	account.tree.Update(uint64(op.Token2), util.AddAmount(token2Amount, amount2))
	proof = appendTokenProof(proof, token2Amount, token2Siblings)

	token0Amount, token0Siblings := account.tree.GetProof(FeeTokenIndex)
	account.tree.Update(FeeTokenIndex, util.SubAmount(token0Amount, fee1))
	proof = appendTokenProof(proof, token0Amount, token0Siblings)

	// update root to merkle tree
	accountHash := crypto.Keccak256Hash(account.tree.rootHash().Bytes(), pubAccountHash.Bytes())
	bc.state.tree.Update(uint64(op.Account1), accountHash)

	account = bc.state.accounts[op.Account2]
	if account == nil {
		panic("empty account")
	}
	_, accountSiblings = bc.state.tree.GetProof(uint64(op.Account2))
	proof = appendSiblings(proof, accountSiblings)
	pubAccountHash = account.GetPubAccountHash()
	proof = append(proof, pubAccountHash.Bytes()...)
	// update balance of token
	token2Amount, token2Siblings = account.tree.GetProof(uint64(op.Token2))
	account.tree.Update(uint64(op.Token2), util.SubAmount(token2Amount, amount2))
	proof = appendTokenProof(proof, token2Amount, token2Siblings)

	token1Amount, token1Siblings = account.tree.GetProof(uint64(op.Token1))
	account.tree.Update(uint64(op.Token1), util.AddAmount(token1Amount, amount1))
	proof = appendTokenProof(proof, token1Amount, token1Siblings)

	token0Amount, token0Siblings = account.tree.GetProof(FeeTokenIndex)
	account.tree.Update(FeeTokenIndex, util.SubAmount(token0Amount, fee2))
	proof = appendTokenProof(proof, token0Amount, token0Siblings)
	// update root to merkle tree
	accountHash = crypto.Keccak256Hash(account.tree.rootHash().Bytes(), pubAccountHash.Bytes())
	bc.state.tree.Update(uint64(op.Account2), accountHash)

	if loo != nil {
		bc.looMax += 1
		_, looSiblings := bc.looState.tree.GetProof(bc.looMax)
		proof = appendSiblings(proof, looSiblings)
		bc.looState.tree.Update(bc.looMax, loo.Hash())
		bc.looState.loos[bc.looMax] = loo
	}
	fee = new(big.Int).Add(fee1, fee2)
	return
}

func (bc *Blockchain) handleTotalFee(fee *big.Int) (proof hexutil.Bytes) {
	account, ok := bc.state.accounts[AdminIndex]
	if !ok {
		panic("no admin account")
	}

	_, accountSiblings := bc.state.tree.GetProof(uint64(AdminIndex))
	proof = appendSiblings(proof, accountSiblings)

	pubAccountHash := account.GetPubAccountHash()
	proof = append(proof, pubAccountHash.Bytes()...)

	feeAmount, feeSiblings := account.tree.GetProof(uint64(FeeTokenIndex))
	proof = appendTokenProof(proof, feeAmount, feeSiblings)

	account.tree.Update(uint64(FeeTokenIndex), util.AddAmount(feeAmount, fee))
	accountHash := crypto.Keccak256Hash(account.tree.rootHash().Bytes(), pubAccountHash.Bytes())

	bc.state.tree.Update(uint64(AdminIndex), accountHash)
	return proof
}

func appendTokenProof(proof hexutil.Bytes, tokenAmount common.Hash, siblings []common.Hash) hexutil.Bytes {
	proof = append(proof, tokenAmount.Bytes()...)
	for _, hash := range siblings {
		proof = append(proof, hash.Bytes()...)
	}
	return proof
}

func appendSiblings(proof hexutil.Bytes, siblings []common.Hash) hexutil.Bytes {
	for _, hash := range siblings {
		proof = append(proof, hash.Bytes()...)
	}
	return proof
}

func (bc *Blockchain) GetStateData() *StateData {
	return &StateData{
		StateRoot:  bc.state.tree.rootHash(),
		LOORoot:    bc.looState.tree.rootHash(),
		AccountMax: bc.accountMax,
		LOOMax:     bc.looMax,
	}
}

type LeftOverOrderList struct {
	loos map[uint64]*types.LeftOverOrder
	tree *MerkleTree
}

func NewLOOList() *LeftOverOrderList {
	return &LeftOverOrderList{
		loos: make(map[uint64]*types.LeftOverOrder),
		tree: NewTree(LOOTreeDeep),
	}
}
