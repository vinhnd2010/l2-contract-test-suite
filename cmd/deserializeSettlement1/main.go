package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/KyberNetwork/l2-contract-test-suite/types"
)

const output = "testdata/deserializeSettlement1.json"

type DeserializeTestSuit struct {
	Data hexutil.Bytes
	Op   types.Settlement1
}

func main() {
	var testSuits []DeserializeTestSuit
	for _, settlement := range []types.Settlement1{
		{
			OpType:   types.SettlementOp11,
			Token1:   1,
			Token2:   2,
			Account1: 14,
			Account2: 15,
			Rate1: types.PackedAmount{
				Mantisa: 1,
				Exp:     18,
			},
			Rate2: types.PackedAmount{
				Mantisa: 2,
				Exp:     17,
			},
			Amount1: types.PackedAmount{
				Mantisa: 2,
				Exp:     14,
			},
			Amount2: types.PackedAmount{
				Mantisa: 3,
				Exp:     13,
			},
			Fee1: types.PackedFee{
				Mantisa: 1,
				Exp:     6,
			},
			Fee2: types.PackedFee{
				Mantisa: 2,
				Exp:     7,
			},
			ValidSince1:  1600331441,
			ValidSince2:  1600331442,
			ValidPeriod1: 86400,
			ValidPeriod2: 86401,
		},
	} {
		testSuits = append(testSuits, DeserializeTestSuit{
			Data: settlement.ToBytes(),
			Op:   settlement,
		})
	}

	b, err := json.MarshalIndent(testSuits, "", "  ")
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(output, b, 0644); err != nil {
		panic(err)
	}
}
