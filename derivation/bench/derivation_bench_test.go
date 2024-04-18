package bench

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/coinbase/kryptology/pkg/core/curves"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/sha3"

	"bip32_threshold_wallet/derivation"
	"bip32_threshold_wallet/tvrf"
	"bip32_threshold_wallet/utils"
)

var (
	curve  = curves.K256()
	sha256 = sha3.New256()

	reuseKeyPair             = true
	optimizedTvrfCombination = false

	// Number of children to derive per benchmark evaluation.
	numChildren = uint32(1)
)

type thresholdParam struct {
	t uint32
	n uint32
}

var benchmarkParams = []thresholdParam{
	{t: 1, n: 3},
	{t: 4, n: 10},
	{t: 7, n: 16},
	{t: 15, n: 32},
	{t: 31, n: 64},
	{t: 63, n: 128},
}

func BenchmarkMultipleTVRFDerivations(b *testing.B) {
	log.Info("------------------- BENCHMARK TVRF HARDENED NODE DERIVATION --------------------")
	log.Infof("No. children to derive: %d, reuse key-pair: %t, optimized TVRF: %t", numChildren, reuseKeyPair, optimizedTvrfCombination)
	log.Infof("Number of CPUs available: %d", runtime.NumCPU())

	for _, param := range benchmarkParams {
		runName := fmt.Sprintf("Run t=%d, n=%d", param.t, param.n)
		b.Run(runName, func(b *testing.B) {
			benchmarkTVRFDerivation(b, param.t, param.n)
		})
	}
}

func benchmarkTVRFDerivation(b *testing.B, t, n uint32) {
	devices := utils.CreateDevices(t, n)
	ddhTvrf := tvrf.NewDDHTVRF(t, n, curve, sha256, optimizedTvrfCombination)
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
