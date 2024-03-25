package main

// Device represents a secret shared non-hardened node.
type Device struct {
	state State

	deviceIdx      int // Index of the device with respect to the secret sharing.
	secretKeyShare []byte
	publicKey      []byte
}

type State struct {
	nodeIdx   int // Index of the node in the derivation tree.
	chainCode []byte
}
