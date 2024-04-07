package tvrf_test

import (
	"testing"

	"github.com/coinbase/kryptology/pkg/core/curves"
	"github.com/coinbase/kryptology/pkg/tecdsa/gg20/dealer"
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
	ski, _ := p256.Scalar.SetBytes(sharesMap[1].ShamirShare.Value.Bytes())
	pkiPoint, _ := p256.Point.Set(sharesMap[1].Point.X, sharesMap[1].Point.Y)
	pki := tvrf.PublicKeyShare{
		Idx:   0,
		Value: &pkiPoint,
	}
	ddhTvrf.PEval(message, &ski, pki)
}
