package bench

import (
	"log"
	"testing"

	"github.com/coinbase/kryptology/pkg/core/curves"
	"golang.org/x/crypto/sha3"

	"bip32_threshold_wallet/derivation"
	"bip32_threshold_wallet/tvrf"
	"bip32_threshold_wallet/utils"
)

var (
	curve  = curves.K256()
	sha256 = sha3.New256()

	threshold    = uint32(5)
	numParties   = uint32(20)
	reuseKeyPair = true

	// Number of children to derive per benchmark evaluation.
	numChildren = uint32(10)
)

func init() {
	log.Printf("------------------- BENCHMARK TVRF HARDENED NODE DERIVATION --------------------")
	log.Printf("t: %d, n: %d, num children: %d, reuse keypair: %t", threshold, numParties, numChildren, reuseKeyPair)
}

func BenchmarkTVRFDerivation(b *testing.B) {
	devices := utils.CreateDevices(threshold, numParties)
	ddhTvrf := tvrf.NewDDHTVRF(threshold, numParties, curve, sha256)
	deriv := derivation.NewTVRFDerivation(curve, devices, ddhTvrf, reuseKeyPair)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		deriveChildren(b, deriv, numChildren)
	}
}

func deriveChildren(b *testing.B, deriv derivation.TVRFDerivation, numChildren uint32) {
	for i := uint32(0); i < numChildren; i++ {
		err, _ := deriv.DeriveHardenedChild(i)
		if err != nil {
			b.Fatal(err)
		}
	}
}
