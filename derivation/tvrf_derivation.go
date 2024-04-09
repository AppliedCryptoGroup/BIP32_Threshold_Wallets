package derivation

import (
	"encoding/binary"

	"github.com/coinbase/kryptology/pkg/core/curves"
	"github.com/pkg/errors"

	"bip32_threshold_wallet/node"
	"bip32_threshold_wallet/tvrf"
)

type TVRFDerivation struct {
	curve   *curves.Curve
	devices []node.Device
	tvrf    tvrf.TVRF

	reuseKeyPair bool
}

// NewTVRFDerivation creates a new TVRF derivation instance.
// All devices will participate in the derivation but note that only t of them suffice.
func NewTVRFDerivation(curve *curves.Curve, devices []node.Device, tvrf tvrf.TVRF, reuseKeyPair bool) TVRFDerivation {
	if !reuseKeyPair {
		panic("Separate TVRF key pairs are not supported yet")
	}

	return TVRFDerivation{
		curve:        curve,
		devices:      devices,
		tvrf:         tvrf,
		reuseKeyPair: reuseKeyPair,
	}
}

func (td TVRFDerivation) DeriveNonHardenedChild(childIdx uint32) (error, []node.Device) {
	nonHardDerivation := NonHardDerivation{devices: td.devices}
	return nonHardDerivation.DeriveNonHardenedChild(childIdx)
}

func (td TVRFDerivation) DeriveHardenedChild(childIdx uint32) (error, *node.Node) {
	// convert childIdx to byte array
	childIdxBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(childIdxBytes, childIdx)

	evals := make([]*tvrf.PartialEvaluation, len(td.devices))

	for i, d := range td.devices {
		dSk, dPk := d.KeyPair()
		err, sk, pk := tvrf.ShamirShareToKeyPair(td.curve, dSk, dPk)
		if err != nil {
			return errors.Wrap(err, "converting key pairs"), nil
		}

		eval, err := td.tvrf.PEval(childIdxBytes, sk, *pk)
		if err != nil {
			return errors.Wrap(err, "evaluation failed"), nil
		}
		evals[i] = eval
	}

	// TODO: Mock sending evaluations to the child node with some networking delay

	// Combine the evaluations.
	combinedEval, err := td.tvrf.Combine(evals)
	if err != nil {
		return errors.Wrap(err, "combining evaluations"), nil
	}

	// Verify the combined evaluation.
	valid := td.tvrf.Verify(*combinedEval)
	if !valid {
		return errors.New("verification of combined evaluation failed"), nil
	}

	// sk,pk= ECDSAKeyGen(combinedEval)
	// node := node.NewNode(sk, pk)

	return errors.New("not implemented"), nil
}
