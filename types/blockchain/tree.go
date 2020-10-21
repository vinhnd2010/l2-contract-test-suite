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
		root: newNode(common.HexToHash(zeroHash), deep-1, nil),
	}
}

func (tr *MerkleTree) GetProof(k uint64) (common.Hash, []common.Hash) {
	return tr.root.getProof(k)
}

func (tr *MerkleTree) GetProofBatch(k []uint64) (values []common.Hash, siblings []common.Hash) {
	var nodes []*node
	var tmpKey []uint64
	tmpKey = append(tmpKey, k...)
	for i := 0; i < len(k); i++ {
		node := tr.getNode(k[i])
		if node == nil {
			panic("invalid node batch proof")
		} else {
			values = append(values, node.value)
		}
		nodes = append(nodes, node)
	}

	for deep := uint(0); deep < tr.deep-1; deep++ {
		var tmpKeys2 []uint64
		var tmpNodes []*node
		for i := 0; i < len(tmpKey); {
			if (i != len(tmpKey)-1) && (tmpKey[i]/2 == tmpKey[i+1]/2) {
				tmpKeys2 = append(tmpKeys2, tmpKey[i]/2)
				tmpNodes = append(tmpNodes, nodes[i].parent)
				i += 2
				continue
			}

			if tmpKey[i]%2 == 0 {
				siblings = append(siblings, nodes[i].parent.rightHash())
			} else {
				siblings = append(siblings, nodes[i].parent.leftHash())
			}
			tmpKeys2 = append(tmpKeys2, tmpKey[i]/2)
			tmpNodes = append(tmpNodes, nodes[i].parent)
			i++
		}
		tmpKey = tmpKeys2
		nodes = tmpNodes
	}
	return values, siblings
}

func (tr *MerkleTree) getNode(k uint64) *node {
	return tr.root.getNode(k)
}

func (tr *MerkleTree) Update(k uint64, v common.Hash) {
	tr.root.update(k, v)
}

func (tr *MerkleTree) RootHash() common.Hash {
	return tr.root.value
}

type node struct {
	deep   uint
	value  common.Hash
	left   *node
	right  *node
	parent *node
}

func newNode(value common.Hash, deep uint, parent *node) *node {
	return &node{
		deep:   deep,
		value:  value,
		parent: parent,
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

func (n *node) getNode(key uint64) *node {
	if n.deep == 0 {
		return n
	}
	isLeft := ((key >> (n.deep - 1)) & 1) == 0
	if isLeft {
		if n.left == nil { // key in null branch
			return nil
		} else {
			return n.left.getNode(key)
		}
	} else {
		if n.right == nil {
			return nil
		} else {
			return n.right.getNode(key)
		}
	}
}

func (n *node) update(key uint64, value common.Hash) {
	if n.deep == 0 {
		n.value = value
		return
	}

	isLeft := ((key >> (n.deep - 1)) & 1) == 0
	if isLeft {
		if n.left == nil {
			n.left = newNode(common.HexToHash(zeroHash), n.deep-1, n)
		}
		n.left.update(key, value)
	} else {
		if n.right == nil {
			n.right = newNode(common.HexToHash(zeroHash), n.deep-1, n)
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
