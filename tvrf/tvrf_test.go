package tvrf_test

import (
	"testing"

	"github.com/coinbase/kryptology/pkg/core/curves"
	"github.com/coinbase/kryptology/pkg/tecdsa/gg20/dealer"
	"github.com/stretchr/testify/assert"
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

	pk, sharesMap, _ := dealer.NewDealerShares(p256ec, threshold, numParties, secret)
	_ = pk

	ddhTvrf := tvrf.NewDDHTVRF(threshold, numParties, p256, sha256)

	message := []byte("Hello, World!")
	err, ski, pki := tvrf.ShamirShareToKeyPair(p256, sharesMap[1].ShamirShare, sharesMap[1].Point)
	assert.NoError(t, err)

	ddhTvrf.PEval(message, ski, *pki)
}
