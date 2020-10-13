package blockchain

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const zeroHash = "0x0"

type MerkleTree struct {
	deep uint
	root *node
}

func NewTree(deep uint) *MerkleTree {
	return &MerkleTree{
		deep: deep,
		root: newNode(common.HexToHash(zeroHash), deep-1),
	}
}

func (tr *MerkleTree) GetProof(k uint64) (common.Hash, []common.Hash) {
	return tr.root.getProof(k)
}

func (tr *MerkleTree) Update(k uint64, v common.Hash) {
	tr.root.update(k, v)
}

func (tr *MerkleTree) RootHash() common.Hash {
	return tr.root.value
}

type node struct {
	deep  uint
	value common.Hash
	left  *node
	right *node
}

func newNode(value common.Hash, deep uint) *node {
	return &node{
		deep:  deep,
		value: value,
	}
}

func (n *node) leftHash() common.Hash {
	if n.left == nil {
		return common.HexToHash(zeroHash)
	}
	return n.left.value
}

func (n *node) rightHash() common.Hash {
	if n.right == nil {
		return common.HexToHash(zeroHash)
	}
	return n.right.value
}

func (n *node) getProof(key uint64) (value common.Hash, siblings []common.Hash) {
	if n.deep == 0 {
		return n.value, nil
	}

	isLeft := ((key >> (n.deep - 1)) & 1) == 0
	if isLeft {
		if n.left == nil { // key in null branch
			for i := 0; i < int(n.deep)-1; i++ {
				siblings = append(siblings, common.HexToHash(zeroHash))
			}
			value = common.HexToHash(zeroHash)
		} else {
			value, siblings = n.left.getProof(key)
		}
		siblings = append(siblings, n.rightHash())
	} else {
		if n.right == nil {
			for i := 0; i < int(n.deep)-1; i++ {
				siblings = append(siblings, common.HexToHash(zeroHash))
			}
			value = common.HexToHash(zeroHash)
		} else {
			value, siblings = n.right.getProof(key)
		}
		siblings = append(siblings, n.leftHash())
	}
	return
}

func (n *node) update(key uint64, value common.Hash) {
	if n.deep == 0 {
		n.value = value
		return
	}

	isLeft := ((key >> (n.deep - 1)) & 1) == 0
	if isLeft {
		if n.left == nil {
			n.left = newNode(common.HexToHash(zeroHash), n.deep-1)
		}
		n.left.update(key, value)
	} else {
		if n.right == nil {
			n.right = newNode(common.HexToHash(zeroHash), n.deep-1)
		}
		n.right.update(key, value)
	}
	// update the root value
	n.value = GetRoot(n.leftHash(), n.rightHash())
}

//func GetRootTree(key uint64 , child common.Hash, siblings []common.Hash) common.Hash {
//	for _, siblings.
//}

func GetRoot(left common.Hash, right common.Hash) common.Hash {
	if left == common.HexToHash(zeroHash) && right == common.HexToHash(zeroHash) {
		return common.HexToHash(zeroHash)
	}
	return crypto.Keccak256Hash(left.Bytes(), right.Bytes())
}
