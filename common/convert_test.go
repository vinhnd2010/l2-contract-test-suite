package common

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUint48ToBytes(t *testing.T) {
	var x uint64 = 1
	data := hexutil.Encode(Uint48ToBytes(x))
	require.Equal(t, "0x000000000001", data)
}
