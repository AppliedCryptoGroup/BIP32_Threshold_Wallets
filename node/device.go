package node

import "github.com/coinbase/kryptology/pkg/tecdsa/gg20/dealer"

type SecretKeyShare *dealer.Share

// Device represents a secret shared non-hardened node.
type Device struct {
	state State

	deviceIdx      int // Index of the device with respect to the secret sharing.
	secretKeyShare SecretKeyShare
	publicKey      *PublicKey
}
