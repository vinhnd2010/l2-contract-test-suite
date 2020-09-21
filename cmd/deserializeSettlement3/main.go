package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/KyberNetwork/l2-contract-test-suite/types"
)

const output = "testdata/deserializeSettlement3.json"

type DeserializeTestSuit struct {
	Data hexutil.Bytes
	Op   types.Settlement3
}

func main() {
	var testSuits []DeserializeTestSuit
	for _, settlement := range []types.Settlement3{
		{
			OpType:      types.SettlementOp3,
			LeftoverID1: 1,
			LeftoverID2: 2,
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
