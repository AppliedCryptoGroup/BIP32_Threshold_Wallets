package tvrf_test

import (
	"testing"

	"github.com/coinbase/kryptology/pkg/core/curves"
	"github.com/coinbase/kryptology/pkg/tecdsa/gg20/dealer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/sha3"

	"bip32_threshold_wallet/tvrf"
)

var (
	p256       = curves.P256()
	sha256     = sha3.New256()
	threshold  = uint32(3)
	numParties = uint32(5)
)

func TestTVRF(t *testing.T) {
	p256ec, _ := p256.ToEllipticCurve()
	secret, _ := dealer.NewSecret(p256ec)

	_, sharesMap, _ := dealer.NewDealerShares(p256ec, threshold, numParties, secret)

	ddhTvrf := tvrf.NewDDHTVRF(threshold, numParties, p256, sha256)

	message := []byte("Hello, World!")

	var secretKeys []tvrf.SecretKeyShare
	var publicKeys []tvrf.PublicKeyShare
	for i := uint32(1); i <= numParties; i++ {
		err, ski, pki := tvrf.ShamirShareToKeyPair(p256, sharesMap[i].ShamirShare, sharesMap[i].Point)
		require.NoError(t, err)
		secretKeys = append(secretKeys, ski)
		publicKeys = append(publicKeys, *pki)
	}

	t.Run("Verify partial evaluation", func(t *testing.T) {
		peval, err := ddhTvrf.PEval(message, secretKeys[0], publicKeys[0])
		assert.NoError(t, err)

		// Verify the evaluation
		valid := ddhTvrf.VerifyPartialEval(peval)
		assert.Truef(t, valid, "evaluation verification failed")

		fakeEval := &tvrf.PartialEvaluation{
			PubKeyShare: peval.PubKeyShare,
			Eval:        peval.Eval.Mul(p256.Scalar.New(3)),
			Proof:       peval.Proof,
		}
		notValid := ddhTvrf.VerifyPartialEval(fakeEval)
		assert.Falsef(t, notValid, "evaluation should not be valid")
	})

	t.Run("Combine partial evaluations", func(t *testing.T) {
		pevals := make([]*tvrf.PartialEvaluation, 0)
		for i := uint32(0); i < threshold; i++ {
			peval, err := ddhTvrf.PEval(message, secretKeys[i], publicKeys[i])
			require.NoError(t, err)
			pevals = append(pevals, peval)
		}

		eval, err := ddhTvrf.Combine(pevals)
		require.NoError(t, err)

		// Verify the evaluation
		valid := ddhTvrf.Verify(*eval)
		assert.Truef(t, valid, "evaluation verification failed")
	})

}
