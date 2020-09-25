package proof

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"testing"
)

func TestBuildMiniBlockProof(t *testing.T) {
	var (
		miniBlockHashes = []common.Hash{
			common.HexToHash("0x1"),
			common.HexToHash("0x2"),
			common.HexToHash("0x3"),
			common.HexToHash("0x4"),
			common.HexToHash("0x5"),
		}
	)
	out :=BuildBlockInfoProof(miniBlockHashes, 1)
	fmt.Println(out)
}
