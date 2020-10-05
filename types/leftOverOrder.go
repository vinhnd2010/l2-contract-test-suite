package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	util "github.com/KyberNetwork/l2-contract-test-suite/common"
)

type LeftOverOrder struct {
	AccountID   uint32
	SrcToken    uint16
	DestToken   uint16
	Amount      *big.Int
	Fee         *big.Int
	Rate        *big.Int
	ValidSince  uint32
	ValidPeriod uint32
}

func (loo *LeftOverOrder) Bytes() hexutil.Bytes {
	var out hexutil.Bytes
	out = append(out, util.Uint32ToBytes(loo.AccountID)...)
	out = append(out, util.Uint16ToBytes(loo.SrcToken)...)
	out = append(out, util.Uint16ToBytes(loo.DestToken)...)
	out = append(out, common.BigToHash(loo.Amount).Bytes()...)
	out = append(out, common.BigToHash(loo.Fee).Bytes()...)
	out = append(out, common.BigToHash(loo.Rate).Bytes()...)
	out = append(out, util.Uint32ToBytes(loo.ValidSince)...)
	out = append(out, util.Uint32ToBytes(loo.ValidPeriod)...)
	return out
}

func (loo *LeftOverOrder) Hash() common.Hash {
	return crypto.Keccak256Hash(
		util.Uint32ToBytes(loo.AccountID), util.Uint16ToByte(loo.SrcToken), util.Uint16ToByte(loo.DestToken),
		common.BigToHash(loo.Amount).Bytes(), common.BigToHash(loo.Fee).Bytes(), common.BigToHash(loo.Rate).Bytes(),
		util.Uint32ToBytes(loo.ValidSince), util.Uint32ToBytes(loo.ValidPeriod),
	)
}

func (loo *LeftOverOrder) Clone() *LeftOverOrder {
	return &LeftOverOrder{
		AccountID:   loo.AccountID,
		SrcToken:    loo.SrcToken,
		DestToken:   loo.DestToken,
		Amount:      new(big.Int).Set(loo.Amount),
		Fee:         new(big.Int).Set(loo.Fee),
		Rate:        new(big.Int).Set(loo.Rate),
		ValidSince:  loo.ValidSince,
		ValidPeriod: loo.ValidPeriod,
	}
}
