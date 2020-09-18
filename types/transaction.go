package types

import (
	"math/big"

	"github.com/KyberNetwork/l2-contract-test-suite/common"
)

type OpType uint8

const (
	NoOp OpType = iota
	SettlementOp11
	SettlementOp12
	SettlementOp13
	SettlementOp21
	SettlementOp22
	SettlementOp3
	DepositToNew
	Deposit
	Withdraw
	Exit
)

/// 10 bit for mantisa, 6 bit for
type Fee struct {
	Mantisa uint16
	Exp     uint8
}

// 10 bit for mantisa, 6 bit for exp
func (f *Fee) toBytes() []byte {
	out := uint16(0)
	out = out | (f.Mantisa << 6)
	out = out | (uint16(f.Exp))
	return common.Uint16ToByte(out)
}

func (f *Fee) MarshalText() ([]byte, error) {
	tmp := big.NewInt(int64(f.Exp))
	tmp = new(big.Int).Exp(big.NewInt(10), tmp, nil)
	tmp = new(big.Int).Mul(tmp, big.NewInt(int64(f.Mantisa)))
	return []byte("0x" + tmp.Text(16)), nil
}

/// @dev 32 bits for mantisa, 8 bits for exp
type Amount struct {
	Mantisa uint32
	Exp     uint8
}

func (a *Amount) toBytes() []byte {
	var out []byte
	out = append(out, common.Uint32ToBytes(a.Mantisa)...)
	out = append(out, common.Uint8ToByte(a.Exp))
	return out
}

func (a *Amount) MarshalText() ([]byte, error) {
	tmp := big.NewInt(int64(a.Exp))
	tmp = new(big.Int).Exp(big.NewInt(10), tmp, nil)
	tmp = new(big.Int).Mul(tmp, big.NewInt(int64(a.Mantisa)))
	return []byte("0x" + tmp.Text(16)), nil
}

type Settlement1 struct {
	OpType       OpType
	Token1       uint16
	Token2       uint16
	Account1     uint32
	Account2     uint32
	Rate1        Amount
	Rate2        Amount
	Amount1      Amount
	Amount2      Amount
	Fee1         Fee
	Fee2         Fee
	ValidSince1  uint32
	ValidSince2  uint32
	ValidPeriod1 uint32
	ValidPeriod2 uint32
}

func (s *Settlement1) ToBytes() []byte {
	var out []byte
	// the first 3 bytes, 4 bit opType, 10 bits token1, 10 bits token2
	head := uint32(0)
	head = head | (uint32(s.OpType) << 20)
	head = head | (uint32(s.Token1) << 10)
	head = head | (uint32(s.Token2))
	out = append(out, common.Uint32ToBytes(head)[1:]...)
	out = append(out, common.Uint32ToBytes(s.Account1)...)
	out = append(out, common.Uint32ToBytes(s.Account2)...)
	out = append(out, s.Amount1.toBytes()...)
	out = append(out, s.Amount2.toBytes()...)
	out = append(out, s.Rate1.toBytes()...)
	out = append(out, s.Rate2.toBytes()...)
	out = append(out, s.Fee1.toBytes()...)
	out = append(out, s.Fee2.toBytes()...)
	out = append(out, common.Uint32ToBytes(s.ValidSince1)...)
	out = append(out, common.Uint32ToBytes(s.ValidSince2)...)

	data1 := common.Uint32ToBytes(s.ValidPeriod1 << 4)
	data2 := common.Uint32ToBytes(s.ValidPeriod2)
	var tmp  byte = data1[3] | data2[0]
	out = append(out, data1[:3]...)
	out = append(out, tmp)
	out = append(out, data2[1:]...)
	return out
}
