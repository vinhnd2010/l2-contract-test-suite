package common

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
)

func Uint16ToBytes(i uint16) []byte {
	var bur bytes.Buffer
	err := binary.Write(&bur, binary.BigEndian, &i)
	if err != nil {
		panic("Uint16ToBytes")
	}
	return bur.Bytes()
}

func Uint32ToBytes(i uint32) []byte {
	var bur bytes.Buffer
	err := binary.Write(&bur, binary.BigEndian, &i)
	if err != nil {
		panic("Uint32ToBytes")
	}
	return bur.Bytes()
}

func Uint48ToBytes(i uint64) []byte {
	var bur bytes.Buffer
	err := binary.Write(&bur, binary.BigEndian, &i)
	if err != nil {
		panic("Uint48ToBytes")
	}
	return bur.Bytes()[2:]
}

func Uint8ToByte(i uint8) byte {
	return byte(i)
}

func Uint16ToByte(i uint16) []byte {
	var bur bytes.Buffer
	err := binary.Write(&bur, binary.BigEndian, &i)
	if err != nil {
		panic("Uint16ToByte")
	}
	return bur.Bytes()
}

func Uint64ToBytes(i uint64) []byte {
	var bur bytes.Buffer
	err := binary.Write(&bur, binary.BigEndian, &i)
	if err != nil {
		panic("Uint64ToBytes")
	}
	return bur.Bytes()
}

func Sha256ToHash(data []byte) common.Hash {
	var out [32]byte
	sha256Hash := sha256.New()
	sha256Hash.Write(data)
	tmp := sha256Hash.Sum(nil)
	if len(tmp) != 32 {
		panic("hash length is not 32")
	}

	for i := 0; i < 32; i++ {
		out[i] = tmp[i]
	}
	return out
}
