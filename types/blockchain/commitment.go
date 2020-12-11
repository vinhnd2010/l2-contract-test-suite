package blockchain

import (
	"bytes"
	"encoding/binary"
	"math/bits"

	"github.com/ethereum/go-ethereum/common/hexutil"

	util "github.com/KyberNetwork/l2-contract-test-suite/common"
	"github.com/KyberNetwork/l2-contract-test-suite/types"
)

func BuildSettlement1ZkMsg(op *types.Settlement1,
	account1PubKey hexutil.Bytes,
	account2PubKey hexutil.Bytes) []byte {
	var (
		out []byte
	)
	msg1 := buildZkMsg(op.Account1, op.Token1, op.Token2, op.Amount1, op.Rate1,
		op.ValidSince1, op.ValidPeriod1, op.Fee1, op.OpType == types.SettlementOp11)
	out = append(out, ReverseBitsForEachByte(msg1)...)
	out = append(out, ReverseBytes(account1PubKey)...)

	msg2 := buildZkMsg(op.Account2, op.Token2, op.Token1, op.Amount2, op.Rate2,
		op.ValidSince2, op.ValidPeriod2, op.Fee2, op.OpType != types.SettlementOp13)
	out = append(out, ReverseBitsForEachByte(msg2)...)
	out = append(out, ReverseBytes(account2PubKey)...)
	return out
}

func ReverseBitsForEachByte(data []byte) []byte {
	var out []byte
	for i := 0; i < len(data); i++ {
		out = append(out, bits.Reverse8(data[i]))
	}
	return out
}

func ReverseBytes(data []byte) []byte {
	var out []byte
	for i := len(data) - 1; i >= 0; i-- {
		out = append(out, data[i])
	}
	return out
}

func buildZkMsg(
	accountID uint32,
	srcTokenID uint16,
	dstTokenID uint16,
	amount types.PackedAmount,
	rate types.PackedAmount,
	validSince uint32,
	validPeriod uint32,
	fee types.PackedFee,
	couldBePartiallyFilled bool,
) []byte {
	var out []byte
	out = append(out, 1)
	out = append(out, util.Uint32ToBytes(accountID)...)
	out = append(out, amount.ToBytes()...)
	out = append(out, rate.ToBytes()...)
	out = append(out, util.Uint32ToBytes(validSince)...)
	// 28 bit validPeriod + 16 bit fee + 10 bit + 10 bit = 64
	var tmp uint64
	tmp |= uint64(validPeriod) << 36
	tmp |= (uint64(fee.Mantisa)<<6 + uint64(fee.Exp)) << 20
	tmp |= uint64(srcTokenID) << 10
	tmp |= uint64(dstTokenID)
	var bur bytes.Buffer
	err := binary.Write(&bur, binary.BigEndian, &tmp)
	if err != nil {
		panic("failed to write to buffer")
	}
	out = append(out, bur.Bytes()...)
	if couldBePartiallyFilled {
		out = append(out, 128)
	} else {
		out = append(out, 0)
	}
	// append the rest
	for len(out) < 32 {
		out = append(out, 0)
	}
	return out
}

func (bc *Blockchain) buildSettlement1ZkMsg(op *types.Settlement1) []byte {
	var (
		account1 *Account
		account2 *Account
		ok       bool
	)
	if account1, ok = bc.state.accounts[op.Account1]; !ok {
		panic("account 1 not found")
	}
	if account2, ok = bc.state.accounts[op.Account2]; !ok {
		panic("account 1 not found")
	}
	return BuildSettlement1ZkMsg(op, account1.pubKey, account2.pubKey)
}

func (bc *Blockchain) buildSettlement2ZkMsg(op *types.Settlement2) []byte {
	var (
		out     []byte
		account *Account
		ok      bool
	)

	loo := bc.looState.loos[op.LooID1]
	msg1 := buildZkMsg(op.AccountID2, loo.DestToken, loo.SrcToken, op.Amount2, op.Rate2,
		op.ValidSince2, op.ValidPeriod2, op.Fee2, op.OpType == types.SettlementOp21)
	out = append(out, ReverseBitsForEachByte(msg1)...)
	if account, ok = bc.state.accounts[op.AccountID2]; !ok {
		panic("account 1 not found")
	}
	out = append(out, ReverseBytes(account.pubKey)...)

	// append the rest
	for len(out) < 128 {
		out = append(out, 0)
	}
	return out
}

// only for the last miniblock
func (bc *Blockchain) BuildCommitmentProof(block *types.MiniBlock) hexutil.Bytes {
	var proof hexutil.Bytes

	for _, tx := range block.Txs {
		switch obj := tx.(type) {
		case *types.Settlement1:
			account := bc.state.accounts[obj.Account1]
			if account == nil {
				panic("account not found")
			}
			proof = append(proof, util.Uint32ToBytes(obj.Account1)...)
			proof = append(proof, account.pubKey...)
			proof = append(proof, account.withdrawTo.Bytes()...)
			proof = append(proof, account.tree.RootHash().Bytes()...)
			_, accountSiblings := bc.state.tree.GetProof(uint64(obj.Account1))
			proof = appendSiblings(proof, accountSiblings)

			account = bc.state.accounts[obj.Account2]
			if account == nil {
				panic("account not found")
			}
			proof = append(proof, util.Uint32ToBytes(obj.Account2)...)
			proof = append(proof, account.pubKey...)
			proof = append(proof, account.withdrawTo.Bytes()...)
			proof = append(proof, account.tree.RootHash().Bytes()...)
			_, accountSiblings = bc.state.tree.GetProof(uint64(obj.Account2))
			proof = appendSiblings(proof, accountSiblings)
			//TODO: add more case for commitment2 and withdraw
		}
	}

	return proof
}
