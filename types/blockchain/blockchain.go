package blockchain

import (
	"fmt"
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
const NumTxPerBLock = 4

type Blockchain struct {
	state       *State
	accountMax  uint32
	looState    *LeftOverOrderList
	looMax      uint64
	numDeposit  uint64
	numWithdraw uint
}

type Genesis struct {
	AccountAlloc map[uint32]GenesisAccount
	AccountMax   uint32

	LooAlloc map[uint64]*types.LeftOverOrder
	LooMax   uint64
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

	looState := NewLOOList()
	state := NewState()
	for accountID, accountAlloc := range genesis.AccountAlloc {
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

	for looID, loo := range genesis.LooAlloc {
		looState.loos[looID] = loo.Clone()
		looState.tree.Update(looID, loo.Hash())
	}

	return &Blockchain{
		state:      state,
		accountMax: genesis.AccountMax,
		looState:   looState,
		looMax:     genesis.LooMax,
	}
}

//func (bc *Blockchain) AddBlock(block *types.Blo)

func (bc *Blockchain) AddMiniBlock(block *types.MiniBlock) []hexutil.Bytes {
	var (
		proofs          []hexutil.Bytes
		totalFee        = big.NewInt(0)
		commitmentInput []byte
	)

	for _, tx := range block.Txs {
		switch obj := tx.(type) {
		case *types.Settlement1:
			commitmentInput = append(commitmentInput, bc.buildSettlement1ZkMsg(obj)...)
		case *types.Settlement2:
			commitmentInput = append(commitmentInput, bc.buildSettlement2ZkMsg(obj)...)
		case *types.WithdrawOp:
			commitmentInput = append(commitmentInput, bc.buildWithdrawZkMsg(obj)...)
		default: // append 128 default bytes
			for i := 0; i < 128; i++ {
				commitmentInput = append(commitmentInput, 0)
			}
		}

		switch obj := tx.(type) {
		case *types.Settlement1:
			proof, fee := bc.handleSettlement1(obj)
			proofs = append(proofs, proof)
			totalFee = totalFee.Add(totalFee, fee)
		case *types.Settlement2:
			proof, fee := bc.handleSettlement2(obj)
			proofs = append(proofs, proof)
			totalFee = totalFee.Add(totalFee, fee)
		case *types.Settlement3:
			proof, fee := bc.handleSettlement3(obj)
			proofs = append(proofs, proof)
			totalFee = totalFee.Add(totalFee, fee)
		case *types.DepositOp:
			proof := bc.handleDeposit(obj)
			proofs = append(proofs, proof)
		case *types.DepositToNewOp:
			proof := bc.handleDepositToNew(obj)
			proofs = append(proofs, proof)
		case *types.WithdrawOp:
			proof, fee := bc.handleWithdraw(obj)
			proofs = append(proofs, proof)
			totalFee = totalFee.Add(totalFee, fee)
		case *types.ExitOp:
			proof := bc.handleExit(obj)
			proofs = append(proofs, proof)
		default:
			panic("unsupported type")
		}
	}
	proofs = append(proofs, bc.handleTotalFee(totalFee))
	block.StateHash = bc.GetStateData().Hash()

	if len(block.Txs) <= NumTxPerBLock {
		for len(commitmentInput) < NumTxPerBLock*128 {
			commitmentInput = append(commitmentInput, 0)
		}
	} else {
		fmt.Println("warning: number of txs >4")
	}
	block.Commitment = util.Sha256ToHash(commitmentInput)
	return proofs
}

func (bc *Blockchain) handleDeposit(op *types.DepositOp) (proof hexutil.Bytes) {
	account := bc.state.accounts[op.AccountID]
	if account == nil {
		panic("empty account")
	}

	_, accountSiblings := bc.state.tree.GetProof(uint64(op.AccountID))
	proof = appendSiblings(proof, accountSiblings)
	pubAccountHash := account.GetPubAccountHash()
	proof = append(proof, pubAccountHash.Bytes()...)
	// update account tree
	tokenAmount, tokenSiblings := account.tree.GetProof(uint64(op.TokenID))
	proof = appendTokenProof(proof, tokenAmount, tokenSiblings)
	account.tree.Update(uint64(op.TokenID), util.AddAmount(tokenAmount, op.Amount))
	// update bc tree
	accountHash := crypto.Keccak256Hash(account.tree.RootHash().Bytes(), pubAccountHash.Bytes())
	bc.state.tree.Update(uint64(op.AccountID), accountHash)

	proof = append(proof, util.Uint32ToBytes(op.AccountID)...)
	proof = append(proof, util.Uint16ToBytes(op.TokenID)...)
	proof = append(proof, common.BigToHash(op.Amount).Bytes()...)

	op.DepositID = bc.numDeposit
	bc.numDeposit++
	return proof
}

func (bc *Blockchain) handleDepositToNew(op *types.DepositToNewOp) (proof hexutil.Bytes) {
	bc.accountMax++
	accountID := bc.accountMax
	account := bc.state.accounts[accountID]
	if account != nil {
		panic("account existed")
	}
	_, siblings := bc.state.tree.GetProof(uint64(accountID))

	account = NewAccount(op.PubKey, op.WithdrawTo)
	account.tree.Update(uint64(op.TokenID), common.BigToHash(op.Amount))
	account.GetPubAccountHash()
	bc.state.accounts[accountID] = account

	accountHash := crypto.Keccak256Hash(account.tree.RootHash().Bytes(), account.GetPubAccountHash().Bytes())
	bc.state.tree.Update(uint64(accountID), accountHash)

	proof = append(proof, op.PubKey...)
	proof = append(proof, op.WithdrawTo.Bytes()...)
	proof = append(proof, util.Uint16ToByte(op.TokenID)...)
	proof = append(proof, common.BigToHash(op.Amount).Bytes()...)
	proof = appendSiblings(proof, siblings)

	op.DepositID = bc.numDeposit
	bc.numDeposit++
	return proof
}

func (bc *Blockchain) updateSettlementBalance(
	accountID1, accountID2 uint32, tokenID1, tokenID2 uint16,
	amount1, amount2, fee1, fee2 *big.Int,
) (proof hexutil.Bytes) {
	account := bc.state.accounts[accountID1]
	if account == nil {
		panic("empty account")
	}

	_, accountSiblings := bc.state.tree.GetProof(uint64(accountID1))
	proof = appendSiblings(proof, accountSiblings)
	pubAccountHash := account.GetPubAccountHash()
	proof = append(proof, pubAccountHash.Bytes()...)
	// update balance of token
	token1Amount, token1Siblings := account.tree.GetProof(uint64(tokenID1))
	account.tree.Update(uint64(tokenID1), util.SubAmount(token1Amount, amount1))
	proof = appendTokenProof(proof, token1Amount, token1Siblings)

	token2Amount, token2Siblings := account.tree.GetProof(uint64(tokenID2))
	account.tree.Update(uint64(tokenID2), util.AddAmount(token2Amount, amount2))
	proof = appendTokenProof(proof, token2Amount, token2Siblings)

	token0Amount, token0Siblings := account.tree.GetProof(FeeTokenIndex)
	account.tree.Update(FeeTokenIndex, util.SubAmount(token0Amount, fee1))
	proof = appendTokenProof(proof, token0Amount, token0Siblings)

	// update root to merkle tree
	accountHash := crypto.Keccak256Hash(account.tree.RootHash().Bytes(), pubAccountHash.Bytes())
	bc.state.tree.Update(uint64(accountID1), accountHash)

	account = bc.state.accounts[accountID2]
	if account == nil {
		panic("empty account")
	}
	_, accountSiblings = bc.state.tree.GetProof(uint64(accountID2))
	proof = appendSiblings(proof, accountSiblings)
	pubAccountHash = account.GetPubAccountHash()
	proof = append(proof, pubAccountHash.Bytes()...)
	// update balance of token
	token2Amount, token2Siblings = account.tree.GetProof(uint64(tokenID2))
	account.tree.Update(uint64(tokenID2), util.SubAmount(token2Amount, amount2))
	proof = appendTokenProof(proof, token2Amount, token2Siblings)

	token1Amount, token1Siblings = account.tree.GetProof(uint64(tokenID1))
	account.tree.Update(uint64(tokenID1), util.AddAmount(token1Amount, amount1))
	proof = appendTokenProof(proof, token1Amount, token1Siblings)

	token0Amount, token0Siblings = account.tree.GetProof(FeeTokenIndex)
	account.tree.Update(FeeTokenIndex, util.SubAmount(token0Amount, fee2))
	proof = appendTokenProof(proof, token0Amount, token0Siblings)
	// update root to merkle tree
	accountHash = crypto.Keccak256Hash(account.tree.RootHash().Bytes(), pubAccountHash.Bytes())
	bc.state.tree.Update(uint64(accountID2), accountHash)

	return proof
}

func (bc *Blockchain) handleSettlement1(op *types.Settlement1) (proof hexutil.Bytes, fee *big.Int) {
	account := bc.state.accounts[op.Account1]
	if account == nil {
		panic("empty account")
	}

	amount1, amount2, fee1, fee2, loo := op.GetSettlementValue()

	proof = append(proof, bc.updateSettlementBalance(op.Account1, op.Account2, op.Token1, op.Token2,
		amount1, amount2, fee1, fee2)...)

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

func (bc *Blockchain) handleSettlement2(op *types.Settlement2) (proof hexutil.Bytes, fee *big.Int) {
	loo, ok := bc.looState.loos[op.LooID1]
	if !ok {
		panic("loo not exist")
	}
	_, looSiblings := bc.looState.tree.GetProof(op.LooID1)
	proof = append(proof, loo.Bytes()...)
	proof = appendSiblings(proof, looSiblings)

	amount1, amount2, fee1, fee2, loo2 := op.GetSettlementValue(loo)

	proof = append(proof,
		bc.updateSettlementBalance(
			loo.AccountID, op.AccountID2, loo.SrcToken, loo.DestToken,
			amount1, amount2, fee1, fee2,
		)...)

	bc.looState.tree.Update(op.LooID1, loo.Hash())
	if loo2 != nil {
		bc.looMax += 1
		_, looSiblings := bc.looState.tree.GetProof(bc.looMax)
		proof = appendSiblings(proof, looSiblings)
		bc.looState.tree.Update(bc.looMax, loo2.Hash())
		bc.looState.loos[bc.looMax] = loo2
	}

	totalFee := new(big.Int).Add(fee1, fee2)
	return proof, totalFee
}

func (bc *Blockchain) handleSettlement3(op *types.Settlement3) (proof hexutil.Bytes, fee *big.Int) {
	loo1, ok := bc.looState.loos[op.LooID1]
	if !ok {
		panic("loo not exist")
	}
	proof = append(proof, loo1.Bytes()...)

	loo2, ok := bc.looState.loos[op.LooID2]
	if !ok {
		panic("loo not exist")
	}
	proof = append(proof, loo2.Bytes()...)

	var (
		orderAmount1                 = new(big.Int).Set(loo1.Amount)
		orderRate1                   = new(big.Int).Set(loo1.Rate)
		orderAmount2                 = new(big.Int).Set(loo2.Amount)
		orderRate2                   = new(big.Int).Set(loo2.Rate)
		amount1, amount2, fee1, fee2 *big.Int
	)

	if loo1.ValidSince <= loo2.ValidSince {
		amount2 = util.CalAmountOut(orderAmount1, orderRate1)
		if amount2.Cmp(orderAmount2) == 1 {
			amount2.Set(orderAmount2)
			amount1 = util.CalAmountIn(orderAmount2, orderRate1)
		} else {
			amount1 = new(big.Int).Set(orderAmount1)
		}
	} else {
		amount1 = util.CalAmountOut(orderAmount2, orderRate2)
		if amount1.Cmp(orderAmount1) == 1 {
			amount1.Set(orderAmount1)
			amount2 = util.CalAmountIn(orderAmount1, orderRate2)
		} else {
			amount2 = new(big.Int).Set(orderAmount2)
		}
	}
	if amount1.Cmp(orderAmount1) < 0 { //left-over loo1 at order 1
		fee1 = new(big.Int).Div(new(big.Int).Mul(loo1.Fee, amount1), orderAmount1)
	} else {
		fee1 = new(big.Int).Set(loo1.Fee)
	}

	_, looSiblings := bc.looState.tree.GetProof(op.LooID1)
	proof = appendSiblings(proof, looSiblings)
	loo1.Amount = orderAmount1.Sub(orderAmount1, amount1)
	loo1.Fee = loo1.Fee.Sub(loo1.Fee, fee1)
	bc.looState.tree.Update(op.LooID1, loo1.Hash())

	if amount2.Cmp(orderAmount2) < 0 { //left-over loo1 at order 2
		fee2 = new(big.Int).Div(new(big.Int).Mul(loo2.Fee, amount2), orderAmount2)
	} else {
		fee2 = new(big.Int).Set(loo2.Fee)
	}
	_, looSiblings = bc.looState.tree.GetProof(op.LooID2)
	proof = appendSiblings(proof, looSiblings)
	loo2.Amount = orderAmount2.Sub(orderAmount2, amount2)
	loo2.Fee = loo2.Fee.Sub(loo2.Fee, fee2)
	bc.looState.tree.Update(op.LooID2, loo2.Hash())

	proof = append(proof,
		bc.updateSettlementBalance(
			loo1.AccountID, loo2.AccountID, loo1.SrcToken, loo2.SrcToken,
			amount1, amount2, fee1, fee2,
		)...)

	totalFee := new(big.Int).Add(fee1, fee2)
	return proof, totalFee
}

func (bc *Blockchain) handleWithdraw(op *types.WithdrawOp) (proof hexutil.Bytes, fee *big.Int) {
	fee = op.Fee.Big()
	amount := op.Amount.Big()
	account := bc.state.accounts[op.AccountID]
	if account == nil {
		panic("empty account")
	}

	_, accountSiblings := bc.state.tree.GetProof(uint64(op.AccountID))
	proof = appendSiblings(proof, accountSiblings)
	pubAccountHash := account.GetPubAccountHash()
	proof = append(proof, account.pubKey...)
	proof = append(proof, account.withdrawTo.Bytes()...)
	// update account tree
	tokenAmount, tokenSiblings := account.tree.GetProof(uint64(op.TokenID))
	proof = appendTokenProof(proof, tokenAmount, tokenSiblings)
	account.tree.Update(uint64(op.TokenID), util.SubAmount(tokenAmount, amount))
	// update token fee
	tokenAmount, tokenSiblings = account.tree.GetProof(uint64(FeeTokenIndex))
	proof = appendTokenProof(proof, tokenAmount, tokenSiblings)
	account.tree.Update(uint64(FeeTokenIndex), util.SubAmount(tokenAmount, fee))
	// update bc tree
	accountHash := crypto.Keccak256Hash(account.tree.RootHash().Bytes(), pubAccountHash.Bytes())
	bc.state.tree.Update(uint64(op.AccountID), accountHash)

	op.WithdrawID = bc.numWithdraw
	bc.numWithdraw++
	return proof, fee
}

func (bc *Blockchain) handleExit(op *types.ExitOp) (proof hexutil.Bytes) {
	account := bc.state.accounts[op.AccountID]
	if account == nil {
		panic("empty account")
	}

	balanceRoot := account.tree.RootHash()
	proof = append(proof, balanceRoot.Bytes()...)

	pubAccountHash := account.GetPubAccountHash()
	proof = append(proof, pubAccountHash.Bytes()...)

	_, accountSiblings := bc.state.tree.GetProof(uint64(op.AccountID))
	proof = appendSiblings(proof, accountSiblings)
	// set balance root of this account to bytes32(0)
	accountHash := crypto.Keccak256Hash(common.HexToHash(zeroHash).Bytes(), pubAccountHash.Bytes())
	bc.state.tree.Update(uint64(op.AccountID), accountHash)
	account.isConfirmedExit = true
	// set balanceRoot to operation
	op.AccountRoot = balanceRoot
	return proof
}

// BuildSubmitExitProof builds a proof for user to submit exit
func (bc *Blockchain) BuildSubmitExitProof(accountID uint32) (balanceRoot common.Hash, proof hexutil.Bytes) {
	account := bc.state.accounts[accountID]
	if account == nil {
		panic("empty account")
	}

	proof = append(proof, account.pubKey...)
	balanceRoot = account.tree.RootHash()

	_, accountSibling := bc.state.tree.GetProof(uint64(accountID))
	proof = appendSiblings(proof, accountSibling)
	proof = append(proof, bc.looState.tree.RootHash().Bytes()...)
	proof = append(proof, util.Uint32ToBytes(bc.accountMax)...)
	proof = append(proof, util.Uint48ToBytes(bc.looMax)...)
	return balanceRoot, proof
}

// BuildCompleteExit builds a proof for user to get tokens and complete exit
func (bc *Blockchain) BuildCompleteExit(accountID uint32, tokenIDs []uint16) (amounts []*big.Int, siblings []common.Hash) {
	var keys []uint64
	for i := 0; i < len(tokenIDs); i++ {
		keys = append(keys, uint64(tokenIDs[i]))
	}
	values, siblings := bc.state.accounts[accountID].tree.GetProofBatch(keys)
	for i := 0; i < len(values); i++ {
		amounts = append(amounts, values[i].Big())
	}
	return amounts, siblings
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
	accountHash := crypto.Keccak256Hash(account.tree.RootHash().Bytes(), pubAccountHash.Bytes())

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
		StateRoot:  bc.state.tree.RootHash(),
		LOORoot:    bc.looState.tree.RootHash(),
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

type BlockData struct {
	MiniBlocks      []*types.MiniBlock
	Timestamp       uint32
	MiniBlockNumber uint
	Proof           *FraudProof
}

type FraudProof struct {
	PrevStateData      *StateData
	PrevStateHashProof hexutil.Bytes
	MiniBlockProof     hexutil.Bytes
	ExecutionProof     []hexutil.Bytes
}
