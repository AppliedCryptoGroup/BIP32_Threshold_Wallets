package node

import (
	"github.com/coinbase/kryptology/pkg/core/curves"
	v1 "github.com/coinbase/kryptology/pkg/sharing/v1"
	"github.com/coinbase/kryptology/pkg/tecdsa/gg20/dealer"
	"golang.org/x/crypto/sha3"
)

type SecretKeyShare *dealer.Share

// Device represents a secret shared non-hardened node.
type Device struct {
	state State

	deviceIdx      int // Index of the device with respect to the secret sharing.
	t              uint32
	n              uint32
	secretKeyShare v1.ShamirShare
	publicKey      *PublicKey
}

var hash = sha3.New256()
var curve = curves.K256()

func (d *Device) RandSk(rho curves.Element) *v1.ShamirShare {
	// For k in [t]: a <- H(rho || k)
	// Let f(x) = rho + a_1*x + ... + a_{t-1}*x^{t-1} + a_t*x^t
	// rho_i = f(i) mod q
	// sk_i' = sk_i + rho_i
	// return sk_i'

	aks := make([]curves.Scalar, d.t)
	var err error
	for i := uint32(0); i < d.t; i++ {
		akiBytes := hash.Sum(append(rho.Bytes(), []byte{byte(i)}...))
		aks[i], err = curve.Scalar.SetBytes(akiBytes)
		if err != nil {
			panic(err)
		}
	}

	return nil
}
