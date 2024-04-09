package node

import "github.com/coinbase/kryptology/pkg/core/curves"

type PublicKey *curves.EcPoint
type SecretKey *curves.Element

// Node represents a hardened node.
type Node struct {
	state State

	secretKey *SecretKey
	publicKey *PublicKey
}

type State struct {
	nodeIdx   uint32 // Index of the node in the derivation tree.
	chainCode []byte
}

func NewNode(index uint32, chainCode []byte, sk *SecretKey, pk *PublicKey) Node {
	return Node{
		state: State{
			nodeIdx:   index,
			chainCode: chainCode,
		},
		secretKey: sk,
		publicKey: pk,
	}
}
