package node

import (
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/coinbase/kryptology/pkg/core/curves"
	v1 "github.com/coinbase/kryptology/pkg/sharing/v1"
	"go.dedis.ch/dela/dkg/pedersen/types"
	"go.dedis.ch/dela/mino"
	"go.dedis.ch/dela/serde"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/suites"
	"golang.org/x/crypto/sha3"
)

var suite = suites.MustFind("Ed25519")
var curve = btcec.S256()

type SecretKeyShare *v1.ShamirShare
type PublicKeyShare *curves.EcPoint

// Device represents a secret shared non-hardened node.
type Device struct {
	state State

	deviceIdx       int // Index of the device with respect to the secret sharing.
	t               uint32
	n               uint32
	secretKeyShare  SecretKeyShare
	publicKeyShare  PublicKeyShare
	publicKeyGlobal PublicKey // Global public key

	mino    mino.Mino
	factory serde.Factory
	privkey kyber.Scalar
}

func NewDevice(idx int, pk PublicKeyShare, sk SecretKeyShare, pkG PublicKey, index uint32, ch []byte, m mino.Mino) (Device, kyber.Point) {
	factory := types.NewMessageFactory(m.GetAddressFactory())

	privkey := suite.Scalar().Pick(suite.RandomStream())
	pubkey := suite.Point().Mul(privkey, nil)

	state := State{
		nodeIdx:   index,
		chainCode: ch,
	}

	return Device{
		state:           state,
		deviceIdx:       idx,
		secretKeyShare:  sk,
		publicKeyShare:  pk,
		publicKeyGlobal: pkG,
		privkey:         privkey,
		mino:            m,
		factory:         factory,
	}, pubkey
}

// RandSk randomizes the device's secret key share using the given randomness rho.
// TODO: Need to check whether rho and the a_i values should be curves.Element or curves.Scalar.
// Shamir shares can be safely transformed to scalars using: curve.Scalar.SetBytes(secretKeyShare.Bytes())
func (d *Device) RandSk(rho curves.Element) *v1.ShamirShare {
	field := curves.NewField(curve.Params().N)

	indexElement := field.NewElement(big.NewInt(int64(d.deviceIdx)))
	rhoPrime := computeCoefficient(rho, int(d.t), field) // The last coefficient.
	for i := int(d.t) - 1; i >= 0; i-- {
		var aki *curves.Element
		// Since F(x) = rho + a_1*x + ... + a_{t-1}*x^{t-1} + a_t*x^t
		if i == 0 {
			aki = rho.Clone()
		} else {
			aki = computeCoefficient(rho, i, field)
		}
		// rho' = rho + ak_i^i * deviceIdx^iExp
		rhoPrime = rhoPrime.Mul(indexElement).Add(aki)
	}

	// sk_i' = sk_i + rho_i'
	skPrime := d.secretKeyShare.Value.Clone().Add(rhoPrime)

	return &v1.ShamirShare{
		Identifier: uint32(d.deviceIdx),
		Value:      skPrime,
	}
}

// KeyPair returns the device's key pair.
// Should only be used for reusing this key-pair for the TVFR.
func (d *Device) KeyPair() (SecretKeyShare, PublicKeyShare) {
	return d.secretKeyShare, d.publicKeyShare
}

// Computes ak_i = H(rho || i)
func computeCoefficient(rho curves.Element, index int, field *curves.Field) *curves.Element {
	akiBytes := sha3.Sum256(append(rho.Bytes(), []byte{byte(index)}...))
	curve := curves.K256()
	scalar, err := curve.Scalar.SetBytes(akiBytes[:])
	println(scalar)
	println(err)

	return field.ElementFromBytes(akiBytes[:])
}
