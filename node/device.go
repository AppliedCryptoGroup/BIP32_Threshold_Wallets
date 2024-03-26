package node

import (
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/coinbase/kryptology/pkg/core/curves"
	v1 "github.com/coinbase/kryptology/pkg/sharing/v1"
	"github.com/coinbase/kryptology/pkg/tecdsa/gg20/dealer"
	"go.dedis.ch/dela/dkg/pedersen/types"
	"go.dedis.ch/dela/mino"
	"go.dedis.ch/dela/serde"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/suites"
	"golang.org/x/crypto/sha3"
)

var suite = suites.MustFind("Ed25519")
var curve = btcec.S256()

type SecretKeyShare *dealer.Share

// Device represents a secret shared non-hardened node.
type Device struct {
	state State

	deviceIdx      int // Index of the device with respect to the secret sharing.
	t              uint32
	n              uint32
	secretKeyShare *v1.ShamirShare
	publicKey      PublicKey

	mino    mino.Mino
	factory serde.Factory
	privkey kyber.Scalar
}

func NewDevice(idx int, sk SecretKeyShare, pk PublicKey, m mino.Mino) (Device, kyber.Point) {
	factory := types.NewMessageFactory(m.GetAddressFactory())

	privkey := suite.Scalar().Pick(suite.RandomStream())
	pubkey := suite.Point().Mul(privkey, nil)

	return Device{
		deviceIdx:      idx,
		secretKeyShare: sk.ShamirShare,
		publicKey:      pk,
		privkey:        privkey,
		mino:           m,
		factory:        factory,
	}, pubkey
}

func (d *Device) RandSk(rho curves.Element) *v1.ShamirShare {
	field := curves.NewField(curve.Params().N)

	//indexElement := field.NewElement(big.NewInt(int64(d.deviceIdx)))
	//rhoPrime := rho.Clone()
	//for i := uint32(0); i < d.t; i++ {
	//	// ak_i = H(rho || i)
	//	akiBytes := hash.Sum(append(rho.Bytes(), []byte{byte(i)}...))
	//	aki := field.ElementFromBytes(akiBytes)
	//	iExp := field.NewElement(big.NewInt(int64(i)))
	//	// rho' = rho + ak_i^i * deviceIdx^iExp
	//	rhoPrime = rhoPrime.Add(aki.Mul(indexElement.Pow(iExp)))
	//}

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

// Computes ak_i = H(rho || i)
func computeCoefficient(rho curves.Element, index int, field *curves.Field) *curves.Element {
	akiBytes := sha3.Sum256(append(rho.Bytes(), []byte{byte(index)}...))
	curve := curves.K256()
	scalar, err := curve.Scalar.SetBytes(akiBytes[:])
	println(scalar)
	println(err)

	return field.ElementFromBytes(akiBytes[:])
}
