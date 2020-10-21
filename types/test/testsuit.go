package test

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"

	"github.com/KyberNetwork/l2-contract-test-suite/types"
	"github.com/KyberNetwork/l2-contract-test-suite/types/blockchain"
)

type Suit struct {
	Msg              string
	GenesisStateHash common.Hash
	Steps            []Step
}
type StepType uint

const (
	NoOp StepType = iota
	SubmitBlock
	AccuseBlockFraudProof
	SubmitDeposit
	CompleteWithdraw
	SubmitExit
	CompleteExit
)

type Step struct {
	Action StepType
	Data   interface{}
}

type SubmitBlockStep struct {
	BlockNumber uint32
	MiniBlocks  []*types.MiniBlock
	Timestamp   uint32
}

type SubmitExitStep struct {
	AccountID   uint32
	BalanceRoot common.Hash
	Timestamp   uint32
	BlockNumber uint32
	Proof       hexutil.Bytes
}

type AccuseBlockFraudProofStep struct {
	BlockNumber        uint
	MiniBlockNumber    uint
	MiniBlock          *types.MiniBlock
	PrevStateData      *blockchain.StateData
	MiniBlockProof     hexutil.Bytes
	PrevStateHashProof hexutil.Bytes
	ExecutionProof     []hexutil.Bytes
}

type CompleteExitStep struct {
	AccountID    uint32
	TokenIDs     []uint16
	TokenAmounts []*big.Int
	Siblings     []common.Hash
}
