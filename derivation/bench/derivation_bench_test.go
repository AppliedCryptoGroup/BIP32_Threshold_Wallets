package bench

import (
	"fmt"
	"runtime"
	"testing"
	"time"

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
	numChildren = 1

	// Simulated network latency
	netLatency = 10 * time.Millisecond
)

type thresholdParam struct {
	t uint32
	n uint32
}

var benchmarkParams = []thresholdParam{
	{t: 2, n: 5},
	{t: 4, n: 10},
	{t: 7, n: 16},
	{t: 15, n: 32},
	{t: 31, n: 64},
	{t: 63, n: 128},
}

func BenchmarkMultipleTVRFDerivations(b *testing.B) {
	log.Info("------------------- BENCHMARK TVRF HARDENED NODE DERIVATION --------------------")
	log.Infof("Reuse key-pair: %t, optimized TVRF: %t", reuseKeyPair, optimizedTvrfCombination)
	log.Infof("Number of CPUs available: %d", runtime.NumCPU())

	if netLatency.Milliseconds() > 0 {
		log.Infof("Simulated network latency: %s", netLatency)
	}

	for _, param := range benchmarkParams {
		runName := fmt.Sprintf("Run t=%d, n=%d", param.t, param.n)
		b.Run(runName, func(b *testing.B) {
			benchmarkTVRFDerivation(b, param.t, param.n)
		})

		// Calculate the total bandwidth used for the derivation which results from all parties sending their evaluation
		// (1 EC point = 64 Bytes) and the proof (2 scalars = 32 Bytes) to the child node.
		bandwidthUsedBits := numChildren * int(param.n) * (64 + 2*32)
		log.Infof("Total bandwidth used: %d Bytes", bandwidthUsedBits)
	}
}

func benchmarkTVRFDerivation(b *testing.B, t, n uint32) {
	devices := utils.CreateDevices(t, n)
	ddhTvrf := tvrf.NewDDHTVRF(t, n, curve, sha256, optimizedTvrfCombination)
	deriv := derivation.NewTVRFDerivation(curve, devices, ddhTvrf, reuseKeyPair)
	deriv.SetNetworkLatency(netLatency)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		deriveChildren(b, deriv, numChildren)
	}
}

func deriveChildren(b *testing.B, deriv derivation.TVRFDerivation, numChildren int) {
	for i := 0; i < numChildren; i++ {
		err, _ := deriv.DeriveHardenedChild(uint32(i))
		if err != nil {
			b.Fatal(err)
		}
	}
}
