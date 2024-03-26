package node

import (
	"github.com/coinbase/kryptology/pkg/tecdsa/gg20/dealer"
	"go.dedis.ch/dela/dkg/pedersen/types"
	"go.dedis.ch/dela/mino"
	"go.dedis.ch/dela/serde"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/suites"
)

var suite = suites.MustFind("Ed25519")

type SecretKeyShare *dealer.Share

// Device represents a secret shared non-hardened node.
type Device struct {
	state State

	deviceIdx      int // Index of the device with respect to the secret sharing.
	secretKeyShare SecretKeyShare
	publicKey      PublicKey

	mino    mino.Mino
	factory serde.Factory
	privkey kyber.Scalar
}

func NewDevice(idx int, sk SecretKeyShare, pk PublicKey, index uint32, ch []byte, m mino.Mino) (Device, kyber.Point) {
	factory := types.NewMessageFactory(m.GetAddressFactory())

	privkey := suite.Scalar().Pick(suite.RandomStream())
	pubkey := suite.Point().Mul(privkey, nil)

	state := State{
		nodeIdx:   index,
		chainCode: ch,
	}

	return Device{
		state:          state,
		deviceIdx:      idx,
		secretKeyShare: sk,
		publicKey:      pk,
		privkey:        privkey,
		mino:           m,
		factory:        factory,
	}, pubkey
}
