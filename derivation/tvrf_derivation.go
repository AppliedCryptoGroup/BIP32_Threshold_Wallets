package derivation

import (
	"crypto/sha256"
	"encoding/binary"
	"math/rand"
	"runtime"

	"github.com/coinbase/kryptology/pkg/core/curves"
	"github.com/pkg/errors"
	"github.com/prometheus/common/log"

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

func (td *TVRFDerivation) DeriveNonHardenedChild(childIdx uint32) (error, []node.Device) {
	nonHardDerivation := NonHardDerivation{devices: td.devices}
	return nonHardDerivation.DeriveNonHardenedChild(childIdx)
}

func (td *TVRFDerivation) DeriveHardenedChild(childIdx uint32) (error, *node.Node) {
	// convert childIdx to byte array
	childIdxBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(childIdxBytes, childIdx)

	evals := make([]*tvrf.PartialEvaluation, len(td.devices))

	// TODO: Parallelize this
	numCPU := runtime.NumCPU()
	if len(td.devices) > numCPU {
		log.Warnf("Number of devices (%d) is greater than number of CPUs (%d)", len(td.devices), numCPU)
	}
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

	//log.Printf("Combining evaluations")
	combinedEval, err := td.tvrf.Combine(evals)
	if err != nil {
		return errors.Wrap(err, "combining evaluations"), nil
	}
	//log.Printf("Combined evaluation: %x", combinedEval.Eval.ToAffineCompressed())

	//log.Printf("Verifying combined evaluation")
	valid := td.tvrf.Verify(*combinedEval)
	if !valid {
		return errors.New("verification of combined evaluation failed"), nil
	}

	//log.Printf("Generating ECDSA key pair for child node")
	sk, pk := td.genECDSAKeyPair(combinedEval)
	child := node.NewNode(childIdx, nil, sk, pk)

	return nil, &child
}

func (td *TVRFDerivation) genECDSAKeyPair(combinedEval *tvrf.Evaluation) (*curves.Scalar, *curves.Point) {
	seed := combinedEval.Eval.ToAffineUncompressed()
	// TODO: More secure way of getting the randomness seed from the evaluation than this?
	hash := sha256.Sum256(seed)
	seedInt := binary.BigEndian.Uint64(hash[:8])
	src := rand.NewSource(int64(seedInt))
	rng := rand.New(src)

	sk := td.curve.Scalar.Random(rng)
	pk := td.curve.ScalarBaseMult(sk)

	return &sk, &pk
}
