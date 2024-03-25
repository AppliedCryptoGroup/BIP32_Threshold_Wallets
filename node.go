package main

// Node represents a hardened node.
type Node struct {
	state State

	secretKey []byte
	publicKey []byte
}
