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
			LeftoverID: 1,
			AccountID:  14,
			Amount1: types.Amount{
				Mantisa: 2,
				Exp:     14,
			},
			Rate1: types.Amount{
				Mantisa: 1,
				Exp:     18,
			},
			Fee1: types.Fee{
				Mantisa: 2,
				Exp:     7,
			},
			ValidSince:  1600331441,
			ValidPeriod: 86400,
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
