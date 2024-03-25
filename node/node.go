package node

// Node represents a hardened node.
type Node struct {
	state State

	secretKey []byte
	publicKey []byte
}

type State struct {
	nodeIdx   uint32 // Index of the node in the derivation tree.
	chainCode []byte
}
