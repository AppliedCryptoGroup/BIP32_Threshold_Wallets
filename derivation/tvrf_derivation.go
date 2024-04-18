package derivation

import (
	"crypto/sha256"
	"encoding/binary"
	"math/rand"
	"runtime"
	"time"

	"github.com/coinbase/kryptology/pkg/core/curves"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"bip32_threshold_wallet/node"
	"bip32_threshold_wallet/tvrf"
)

type TVRFDerivation struct {
	curve   *curves.Curve
	devices []node.Device
	tvrf    tvrf.TVRF

	reuseKeyPair bool

	// netLatency is used to simulate network latency in the derivation process when parties send their evaluations
	// to the child node.
	netLatency time.Duration
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

func (td *TVRFDerivation) SetNetworkLatency(netLatency time.Duration) {
	td.netLatency = netLatency
}

func (td *TVRFDerivation) DeriveNonHardenedChild(childIdx uint32) (error, []node.Device) {
	nonHardDerivation := NonHardDerivation{devices: td.devices}
	return nonHardDerivation.DeriveNonHardenedChild(childIdx)
}

func (td *TVRFDerivation) DeriveHardenedChild(childIdx uint32) (error, *node.Node) {
	childIdxBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(childIdxBytes, childIdx)

	log.Trace("evaluating TVRF for all devices")
	evals, err := td.parallelTVRFEval(childIdxBytes)

	// Simulate network latency, where all parties would send their evaluations in parallel to the child node.
	time.Sleep(td.netLatency)

	log.Trace("combining evaluations")
	combinedEval, err := td.tvrf.Combine(evals)
	if err != nil {
		return errors.Wrap(err, "combining evaluations"), nil
	}
	log.Tracef("combined evaluation: %x", combinedEval.Eval.ToAffineCompressed())

	log.Trace("verifying combined evaluation")
	valid := td.tvrf.Verify(*combinedEval)
	if !valid {
		return errors.New("verification of combined evaluation failed"), nil
	}

	log.Trace("generating ECDSA key pair for child node")
	sk, pk := td.genECDSAKeyPair(combinedEval)
	child := node.NewNode(childIdx, nil, sk, pk)

	return nil, &child
}

func (td *TVRFDerivation) sequentialTVRFEval(childIdxBytes []byte) ([]*tvrf.PartialEvaluation, error) {
	evals := make([]*tvrf.PartialEvaluation, len(td.devices))

	for i, d := range td.devices {
		dSk, dPk := d.KeyPair()
		err, sk, pk := tvrf.ShamirShareToKeyPair(td.curve, dSk, dPk)
		if err != nil {
			return nil, errors.Wrap(err, "converting key pairs")
		}

		eval, err := td.tvrf.PEval(childIdxBytes, sk, *pk)
		if err != nil {
			return nil, errors.Wrap(err, "evaluation failed")
		}
		evals[i] = eval
	}

	return evals, nil
}

func (td *TVRFDerivation) parallelTVRFEval(childIdxBytes []byte) ([]*tvrf.PartialEvaluation, error) {
	devicesChan := make(chan node.Device, len(td.devices))
	errorsChan := make(chan error, len(td.devices))
	evalsChan := make(chan *tvrf.PartialEvaluation, len(td.devices))

	numCPU := runtime.NumCPU()

	for i := 0; i < numCPU; i++ {
		go func() {
			for d := range devicesChan {
				dSk, dPk := d.KeyPair()
				err, sk, pk := tvrf.ShamirShareToKeyPair(td.curve, dSk, dPk)
				if err != nil {
					errorsChan <- errors.Wrap(err, "converting key pairs")
					return
				}

				eval, err := td.tvrf.PEval(childIdxBytes, sk, *pk)
				if err != nil {
					errorsChan <- errors.Wrap(err, "evaluation failed")
					return
				}
				evalsChan <- eval
			}
		}()
	}

	// Send all devices to the goroutines and close the channel to break the loop.
	for _, d := range td.devices {
		devicesChan <- d
	}
	close(devicesChan)

	evals := make([]*tvrf.PartialEvaluation, len(td.devices))
	for i := range td.devices {
		select {
		case err := <-errorsChan:
			return nil, err
		case eval := <-evalsChan:
			evals[i] = eval
		}
	}
	return evals, nil

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
