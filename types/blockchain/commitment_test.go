package blockchain

import (
	"encoding/json"
	"io/ioutil"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"

	"github.com/KyberNetwork/l2-contract-test-suite/testsample"
	"github.com/KyberNetwork/l2-contract-test-suite/types"
)

type Order struct {
	AccountID              uint32             `json:"account_id"`
	SrcTokenID             uint16             `json:"token_id_1"`
	DstTokenID             uint16             `json:"token_id_2"`
	Amount                 types.PackedAmount `json:"amount"`
	Rate                   types.PackedAmount `json:"rate"`
	ValidSince             uint32             `json:"valid_since"`
	ValidPeriod            uint32             `json:"period"`
	Fee                    types.PackedFee
	CouldBePartiallyFilled bool `json:"could_be_partially_filled"`
}

type TestCommitment struct {
	MiniBlockPubData []byte `json:"miniblock_pubdata"`
	FinalCommitment  string `json:"miniblock_commitmemt"`
	Txs              []struct {
		Order1   Order  `json:"order_1"`
		Order2   Order  `json:"order_2"`
		PubData1 []byte `json:"message_to_be_signed_1"`
		PubData2 []byte `json:"message_to_be_signed_2"`
	} `json:"txs"`
}

func TestBuildZkMsg(t *testing.T) {
	data, err := ioutil.ReadFile("../../testdata/miniblock.json")
	require.NoError(t, err)

	var test TestCommitment
	err = json.Unmarshal(data, &test)
	require.NoError(t, err)

	tx := test.Txs[3]
	order1 := tx.Order1
	order2 := tx.Order2

	msg := buildZkMsg(order1.AccountID, order1.SrcTokenID, order1.DstTokenID, order1.Amount, order1.Rate,
		order1.ValidSince, order1.ValidPeriod, order1.Fee, order1.CouldBePartiallyFilled,
	)
	require.Equal(t, tx.PubData1, msg)

	msg = buildZkMsg(order2.AccountID, order2.SrcTokenID, order2.DstTokenID, order2.Amount, order2.Rate,
		order2.ValidSince, order2.ValidPeriod, order2.Fee, order2.CouldBePartiallyFilled,
	)
	require.Equal(t, tx.PubData2, msg)
}

var genesis = &Genesis{
	AccountAlloc: map[uint32]GenesisAccount{
		0: {
			Tokens: map[uint16]*big.Int{
				0: big.NewInt(30000),
				1: big.NewInt(2000000),
			},
			Pubkey:  testsample.PublicKeys[0],
			Address: common.HexToAddress("0xdC70a72AbF352A0E3f75d737430EB896BA9Bf9Ea"),
		},
	},
	AccountMax: 0,
	LooMax:     0,
}

func Test_Commitment(t *testing.T) {
	data, err := ioutil.ReadFile("../../testdata/miniblock.json")
	require.NoError(t, err)

	var test TestCommitment
	err = json.Unmarshal(data, &test)
	require.NoError(t, err)

	bc := NewBlockchain(genesis)

	amount0, _ := new(big.Int).SetString("12000000000000000000", 10)
	miniBlock1 := &types.MiniBlock{
		Txs: []types.Transaction{
			&types.DepositToNewOp{
				PubKey:     hexutil.MustDecode("0xcabf7cd1c1e7954f9d1faf98d604f8cb4772f7c57c9335ad7d16c75a017fc82b"),
				WithdrawTo: common.HexToAddress("0x227b87530bc03015ad5a405d83f3d2f4c5832d12"),
				TokenID:    6,
				Amount:     big.NewInt(100000000000000000),
			}, &types.DepositToNewOp{
				PubKey:     hexutil.MustDecode("0xb60b26f03aa6f2c129f1e4f713965e6ac485fb1e8a1fe91c1399af5df3ff4e09"),
				WithdrawTo: common.HexToAddress("0x320849ec0adffcd6fb0212b59a2ec936cdef5fca"),
				TokenID:    6,
				Amount:     big.NewInt(98000000000000000),
			}, &types.DepositOp{
				AccountID: 1,
				TokenID:   0,
				Amount:    amount0,
			}, &types.Settlement1{
				OpType:   types.SettlementOp11,
				Token1:   0,
				Token2:   6,
				Account1: 1,
				Account2: 2,
				Amount1:  types.PackedAmount{Mantisa: 3, Exp: 18},
				Rate1:    types.PackedAmount{Mantisa: 2, Exp: 16},

				Amount2:      types.PackedAmount{Mantisa: 2, Exp: 16},
				Rate2:        types.PackedAmount{Mantisa: 5, Exp: 19},
				Fee1:         types.PackedFee{Mantisa: 0, Exp: 0},
				Fee2:         types.PackedFee{Mantisa: 0, Exp: 0},
				ValidSince1:  1605323933,
				ValidSince2:  1605323952,
				ValidPeriod1: 268435455,
				ValidPeriod2: 268435455,
			},
		},
	}
	bc.AddMiniBlock(miniBlock1)

	require.Equal(t, miniBlock1.Commitment.Hex()[2:], test.FinalCommitment)
}
