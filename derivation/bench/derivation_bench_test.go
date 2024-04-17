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

	threshold    = uint32(99)
	numParties   = uint32(200)
	reuseKeyPair = true

	// Number of children to derive per benchmark evaluation.
	numChildren = uint32(1)
)

type thresholdParam struct {
	t uint32
	n uint32
}

var sharingParams = []thresholdParam{
	{t: 2, n: 3},
	{t: 4, n: 10},
	{t: 40, n: 100},
	{t: 99, n: 200},
}

//func init() {
//	log.Info("------------------- BENCHMARK TVRF HARDENED NODE DERIVATION --------------------")
//	log.Infof("t: %d, n: %d, num children: %d, reuse key-pair: %t", threshold, numParties, numChildren, reuseKeyPair)
//	numCPU := runtime.NumCPU()
//	if int(numParties) > numCPU {
//		log.Warnf("Number of devices (%d) is greater than number of CPUs (%d)", numParties, numCPU)
//	}
//}

func BenchmarkMultipleTVRFDerivations(b *testing.B) {
	for _, param := range sharingParams {
		runName := fmt.Sprintf("TVRF Derivation for t=%d,n=%d", param.t, param.n)
		b.Run(runName, func(b *testing.B) {
			benchmarkTVRFDerivation(b, param.t, param.n)
		})
	}
}

func BenchmarkTVRFDerivation(b *testing.B) {
	log.Info("------------------- BENCHMARK TVRF HARDENED NODE DERIVATION --------------------")
	log.Infof("t: %d, n: %d, num children: %d, reuse key-pair: %t", threshold, numParties, numChildren, reuseKeyPair)
	numCPU := runtime.NumCPU()
	if int(numParties) > numCPU {
		log.Warnf("Number of devices (%d) is greater than number of CPUs (%d)", numParties, numCPU)
	}
	benchmarkTVRFDerivation(b, threshold, numParties)
}

func benchmarkTVRFDerivation(b *testing.B, t, n uint32) {
	devices := utils.CreateDevices(t, n)
	ddhTvrf := tvrf.NewDDHTVRF(t, n, curve, sha256, true)
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
