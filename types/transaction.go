package types

import (
	"math/big"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/KyberNetwork/l2-contract-test-suite/common"
)

type Transaction interface {
	// ToBytes returns pubData, submit to blockchain
	ToBytes() []byte
}

type OpType uint8

const (
	NoOp           OpType = iota // 0
	SettlementOp11               // 1
	SettlementOp12               // 2
	SettlementOp13               // 3
	SettlementOp21               // 4
	SettlementOp22               // 5
	SettlementOp3                // 6
	DepositToNew                 //7
	Deposit                      //8
	Withdraw                     // 9
	Exit                         // 10
)

var (
	precision = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
)

/// 10 bit for mantisa, 6 bit for
type PackedFee struct {
	Mantisa uint16 `json:"mantisa,string"`
	Exp     uint8  `json:"fee,string"`
}

// 10 bit for mantisa, 6 bit for exp
func (f *PackedFee) toBytes() []byte {
	out := uint16(0)
	out = out | (f.Mantisa << 6)
	out = out | (uint16(f.Exp))
	return common.Uint16ToByte(out)
}

func (f *PackedFee) Big() *big.Int {
	tmp := big.NewInt(int64(f.Exp))
	tmp = new(big.Int).Exp(big.NewInt(10), tmp, nil)
	tmp = new(big.Int).Mul(tmp, big.NewInt(int64(f.Mantisa)))
	return tmp
}

func (f *PackedFee) MarshalText() ([]byte, error) {
	return []byte("0x" + f.Big().Text(16)), nil
}

/// @dev 32 bits for mantisa, 8 bits for exp
type PackedAmount struct {
	Mantisa uint32 `json:"mantissa,string"`
	Exp     uint8  `json:"exp,string"`
}

func (a *PackedAmount) ToBytes() []byte {
	var out []byte
	out = append(out, common.Uint32ToBytes(a.Mantisa)...)
	out = append(out, common.Uint8ToByte(a.Exp))
	return out
}

func (a PackedAmount) Big() *big.Int {
	tmp := big.NewInt(int64(a.Exp))
	tmp = new(big.Int).Exp(big.NewInt(10), tmp, nil)
	tmp = new(big.Int).Mul(tmp, big.NewInt(int64(a.Mantisa)))
	return tmp
}

func (a *PackedAmount) MarshalText() ([]byte, error) {
	return []byte("0x" + a.Big().Text(16)), nil
}

type Settlement1 struct {
	OpType       OpType
	Token1       uint16
	Token2       uint16
	Account1     uint32
	Account2     uint32
	Rate1        PackedAmount
	Rate2        PackedAmount
	Amount1      PackedAmount
	Amount2      PackedAmount
	Fee1         PackedFee
	Fee2         PackedFee
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
	tmp = new(big.Int).Add(tmp, rate)
	tmp = new(big.Int).Sub(tmp, big.NewInt(1))
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
	out = append(out, s.Amount1.ToBytes()...)
	out = append(out, s.Amount2.ToBytes()...)
	out = append(out, s.Rate1.ToBytes()...)
	out = append(out, s.Rate2.ToBytes()...)
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
	OpType       OpType
	LooID1       uint64
	AccountID2   uint32
	Amount2      PackedAmount
	Rate2        PackedAmount
	Fee2         PackedFee
	ValidSince2  uint32
	ValidPeriod2 uint32
}

func (s *Settlement2) ToBytes() []byte {
	var out []byte
	// the first 6 bytes, 4 bit opType, 44 bits LeftOverID
	head := uint64(0)
	head = head | (uint64(s.OpType) << 44)
	head = head | (s.LooID1)
	out = append(out, common.Uint48ToBytes(head)...)
	out = append(out, common.Uint32ToBytes(s.AccountID2)...)
	out = append(out, s.Amount2.ToBytes()...)
	out = append(out, s.Rate2.ToBytes()...)
	out = append(out, s.Fee2.toBytes()...)
	out = append(out, common.Uint32ToBytes(s.ValidSince2)...)
	out = append(out, common.Uint32ToBytes(s.ValidPeriod2<<4)...)
	return out
}

func (s *Settlement2) GetSettlementValue(loo1 *LeftOverOrder) (amount1 *big.Int, amount2 *big.Int, fee1 *big.Int, fee2 *big.Int, loo2 *LeftOverOrder) {
	var (
		orderAmount1 = new(big.Int).Set(loo1.Amount)
		orderRate1   = new(big.Int).Set(loo1.Rate)
		orderAmount2 = s.Amount2.Big()
		orderRate2   = s.Rate2.Big()
	)

	if loo1.ValidSince <= s.ValidSince2 {
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
			amount2 = calAmountIn(orderAmount1, orderRate2)
		} else {
			amount2 = new(big.Int).Set(orderAmount2)
		}
	}

	if amount1.Cmp(orderAmount1) < 0 {
		fee1 = new(big.Int).Div(new(big.Int).Mul(loo1.Fee, amount1), orderAmount1)
	} else {
		fee1 = new(big.Int).Set(loo1.Fee)
	}

	if amount2.Cmp(orderAmount2) < 0 { //left-over loo1 at order 2
		fee2 = new(big.Int).Div(new(big.Int).Mul(s.Fee2.Big(), amount2), orderAmount2)
		loo2 = &LeftOverOrder{
			AccountID:   s.AccountID2,
			SrcToken:    loo1.DestToken,
			DestToken:   loo1.SrcToken,
			Amount:      new(big.Int).Sub(orderAmount2, amount2),
			Rate:        orderRate2,
			Fee:         new(big.Int).Sub(s.Fee2.Big(), fee2),
			ValidSince:  s.ValidSince2,
			ValidPeriod: s.ValidPeriod2,
		}
	} else {
		fee2 = s.Fee2.Big()
	}

	loo1.Amount = orderAmount1.Sub(orderAmount1, amount1)
	loo1.Fee = loo1.Fee.Sub(loo1.Fee, fee1)

	return
}

type Settlement3 struct {
	LooID1 uint64
	LooID2 uint64
}

func (s *Settlement3) ToBytes() []byte {
	var out []byte
	// the first 6 bytes, 4 bit opType, 44 bits LeftOverID1
	head := uint64(0)
	head = head | (uint64(SettlementOp3) << 44)
	head = head | (s.LooID1)
	out = append(out, common.Uint48ToBytes(head)...)
	// next 6 bytes: 44 looID2x
	out = append(out, common.Uint48ToBytes(s.LooID2<<4)...)
	return out
}

type DepositOp struct {
	DepositID uint64
	AccountID uint32
	TokenID   uint16
	Amount    *big.Int
}

func (d *DepositOp) ToBytes() []byte {
	head := uint64(0)
	head = head | (uint64(Deposit) << 44)
	head = head | (d.DepositID)
	return common.Uint48ToBytes(head)
}

type DepositToNewOp struct {
	DepositID  uint64
	PubKey     hexutil.Bytes
	WithdrawTo ethCommon.Address
	TokenID    uint16
	Amount     *big.Int
}

func (d *DepositToNewOp) ToBytes() []byte {
	head := uint64(0)
	head = head | (uint64(DepositToNew) << 44)
	head = head | (d.DepositID)
	return common.Uint48ToBytes(head)
}

type WithdrawOp struct {
	TokenID    uint16
	Amount     PackedAmount
	DestAddr   ethCommon.Address
	AccountID  uint32
	ValidSince uint32
	Fee        PackedFee
	WithdrawID uint
}

func (w *WithdrawOp) ToBytes() []byte {
	var (
		data = uint16(0)
		out  []byte
	)
	data |= uint16(Withdraw) << 12
	data |= w.TokenID << 2
	out = append(out, common.Uint16ToByte(data)...)
	out = append(out, w.Amount.ToBytes()...)
	out = append(out, w.DestAddr.Bytes()...)
	out = append(out, common.Uint32ToBytes(w.AccountID)...)
	out = append(out, common.Uint32ToBytes(w.ValidSince)...)
	out = append(out, w.Fee.toBytes()...)
	return out
}

type ExitOp struct {
	AccountID   uint32
	AccountRoot ethCommon.Hash
}

func (exit *ExitOp) ToBytes() []byte {
	var out []byte
	var data = uint8(Exit) << 4
	out = append(out, data)
	out = append(out, common.Uint32ToBytes(exit.AccountID)...)
	out = append(out, exit.AccountRoot.Bytes()...)
	return out
}
