package types

import (
	"math/big"

	"github.com/KyberNetwork/l2-contract-test-suite/common"
)

type Transaction interface {
	ToBytes() []byte
}

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

var (
	precision = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
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

func (f *Fee) Big() *big.Int {
	tmp := big.NewInt(int64(f.Exp))
	tmp = new(big.Int).Exp(big.NewInt(10), tmp, nil)
	tmp = new(big.Int).Mul(tmp, big.NewInt(int64(f.Mantisa)))
	return tmp
}

func (f *Fee) MarshalText() ([]byte, error) {
	return []byte("0x" + f.Big().Text(16)), nil
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

func (a *Amount) Big() *big.Int {
	tmp := big.NewInt(int64(a.Exp))
	tmp = new(big.Int).Exp(big.NewInt(10), tmp, nil)
	tmp = new(big.Int).Mul(tmp, big.NewInt(int64(a.Mantisa)))
	return tmp
}

func (a *Amount) MarshalText() ([]byte, error) {
	return []byte("0x" + a.Big().Text(16)), nil
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

func calAmountOut(amount *big.Int, rate *big.Int) *big.Int {
	return new(big.Int).Div(new(big.Int).Mul(amount, rate), precision)
}

func calAmountIn(amount *big.Int, rate *big.Int) *big.Int {
	tmp := new(big.Int).Mul(amount, precision)
	tmp = new(big.Int).Add(tmp, big.NewInt(1))
	tmp = new(big.Int).Sub(tmp, rate)
	return tmp.Div(tmp, rate)
}

func (s *Settlement1) GetSettlementValue() (amount1 *big.Int, amount2 *big.Int, fee1 *big.Int, fee2 *big.Int, loo *LeftOverOrder) {
	var (
		orderAmount1 = s.Amount1.Big()
		orderRate1   = s.Rate1.Big()
		orderAmount2 = s.Amount2.Big()
		orderRate2   = s.Rate2.Big()
	)

	if s.ValidSince1 <= s.ValidSince2 {
		amount2 = calAmountOut(orderAmount1, orderRate1)
		if amount2.Cmp(orderAmount2) == 1 {
			amount2.Set(orderAmount2)
			amount1 = calAmountIn(orderAmount2, orderRate1)
		} else {
			amount1 = new(big.Int).Set(orderAmount1)
		}
	} else {
		amount1 = calAmountOut(orderAmount2, orderRate2)
		if amount1.Cmp(orderAmount1) == 1 {
			amount1.Set(orderAmount1)
			amount2 = calAmountIn(orderAmount1, orderAmount2)
		} else {
			amount2 = new(big.Int).Set(orderAmount2)
		}
	}

	fee1 = s.Fee1.Big()
	fee2 = s.Fee2.Big()

	if amount1.Cmp(orderAmount1) < 0 { //left-over order at order1
		fee1 = new(big.Int).Div(new(big.Int).Mul(s.Fee1.Big(), amount1), orderAmount1)
		loo = &LeftOverOrder{
			AccountID:   s.Account1,
			SrcToken:    s.Token1,
			DestToken:   s.Token2,
			Amount:      new(big.Int).Sub(orderAmount1, amount1),
			Rate:        orderRate1,
			Fee:         new(big.Int).Sub(s.Fee1.Big(), fee1),
			ValidSince:  s.ValidSince1,
			ValidPeriod: s.ValidPeriod1,
		}
	}

	if amount2.Cmp(orderAmount2) < 0 { //left-over order at order1
		fee2 = new(big.Int).Div(new(big.Int).Mul(s.Fee2.Big(), amount2), orderAmount2)
		loo = &LeftOverOrder{
			AccountID:   s.Account2,
			SrcToken:    s.Token2,
			DestToken:   s.Token1,
			Amount:      new(big.Int).Sub(orderAmount2, amount2),
			Rate:        orderRate2,
			Fee:         new(big.Int).Sub(s.Fee2.Big(), fee2),
			ValidSince:  s.ValidSince2,
			ValidPeriod: s.ValidPeriod2,
		}
	}
	return
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
	var tmp byte = data1[3] | data2[0]
	out = append(out, data1[:3]...)
	out = append(out, tmp)
	out = append(out, data2[1:]...)
	return out
}

type Settlement2 struct {
	OpType      OpType
	LeftoverID  uint64
	AccountID   uint32
	Amount1     Amount
	Rate1       Amount
	Fee1        Fee
	ValidSince  uint32
	ValidPeriod uint32
}

func (s *Settlement2) ToBytes() []byte {
	var out []byte
	// the first 6 bytes, 4 bit opType, 44 bits LeftOverID
	head := uint64(0)
	head = head | (uint64(s.OpType) << 44)
	head = head | (uint64(s.LeftoverID))
	out = append(out, common.Uint64ToBytes(head)[1:]...)
	out = append(out, common.Uint32ToBytes(s.AccountID)...)
	out = append(out, s.Amount1.toBytes()...)
	out = append(out, s.Rate1.toBytes()...)
	out = append(out, s.Fee1.toBytes()...)
	out = append(out, common.Uint32ToBytes(s.ValidSince)...)

	data := common.Uint32ToBytes(s.ValidPeriod)
	out = append(out, data[:3]...)
	return out
}

type Settlement3 struct {
	OpType      OpType
	LeftoverID1 uint64
	LeftoverID2 uint64
}

func (s *Settlement3) ToBytes() []byte {
	var out []byte
	// the first 6 bytes, 4 bit opType, 44 bits LeftOverID1
	head := uint64(0)
	head = head | (uint64(s.OpType) << 44)
	head = head | (uint64(s.LeftoverID1))
	out = append(out, common.Uint64ToBytes(head)[1:]...)
	out = append(out, common.Uint64ToBytes(s.LeftoverID2)...)
	return out
}
