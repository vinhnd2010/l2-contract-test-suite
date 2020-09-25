package types

import (
	util "github.com/KyberNetwork/l2-contract-test-suite/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

type LeftOverOrder struct {
	AccountID   uint32
	SrcToken    uint16
	DestToken   uint16
	Amount      *big.Int
	Rate        *big.Int
	Fee         *big.Int
	ValidSince  uint32
	ValidPeriod uint32
}

func (loo *LeftOverOrder) Hash() common.Hash {
	return crypto.Keccak256Hash(
		util.Uint32ToBytes(loo.AccountID), util.Uint16ToByte(loo.SrcToken), util.Uint16ToByte(loo.DestToken),
		common.BigToHash(loo.Amount).Bytes(), common.BigToHash(loo.Fee).Bytes(), common.BigToHash(loo.Rate).Bytes(),
		util.Uint32ToBytes(loo.ValidSince), util.Uint32ToBytes(loo.ValidPeriod),
	)
}
