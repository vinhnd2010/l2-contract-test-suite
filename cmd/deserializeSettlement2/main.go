package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/KyberNetwork/l2-contract-test-suite/types"
)

const output = "testdata/deserializeSettlement2.json"

type DeserializeTestSuit struct {
	Data hexutil.Bytes
	Op   types.Settlement2
}

func main() {
	var testSuits []DeserializeTestSuit
	for _, settlement := range []types.Settlement2{
		{
			OpType:     types.SettlementOp21,
			LooID1:     342462,
			AccountID2: 14,
			Amount2: types.PackedAmount{
				Mantisa: 2,
				Exp:     14,
			},
			Rate2: types.PackedAmount{
				Mantisa: 1,
				Exp:     18,
			},
			Fee2: types.PackedFee{
				Mantisa: 2,
				Exp:     7,
			},
			ValidSince2:  1600331441,
			ValidPeriod2: 86400,
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
