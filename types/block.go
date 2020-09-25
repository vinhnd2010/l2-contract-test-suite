package types

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type MiniBlock struct {
	Txs        []Transaction
	StateHash  common.Hash
	Commitment common.Hash
}

func (blk *MiniBlock) MarshalJSON() ([]byte, error) {
	var data hexutil.Bytes

	data = append(data, blk.Commitment.Bytes()...)
	data = append(data, blk.StateHash.Bytes()...)
	for _, tx := range blk.Txs {
		data = append(data, tx.ToBytes()...)
	}
	return json.Marshal(&data)
}

func (blk *MiniBlock) Hash() common.Hash {
	var txData hexutil.Bytes
	for _, tx := range blk.Txs {
		txData = append(txData, tx.ToBytes()...)
	}
	txRoot := crypto.Keccak256Hash(txData)
	return crypto.Keccak256Hash(blk.Commitment.Bytes(), blk.StateHash.Bytes(), txRoot.Bytes())
}
