package derivation_test

import (
	"testing"

	"github.com/coinbase/kryptology/pkg/core/curves"
	"golang.org/x/crypto/sha3"

	"bip32_threshold_wallet/node"
)

var (
	p256       = curves.P256()
	sha256     = sha3.New256()
	threshold  = uint32(3)
	numParties = uint32(5)
)

func TestNewTVRFDerivation(t *testing.T) {

}

func createDevices() []node.Device {
	pkShares, skShares, pk := node.GenSharedKey(threshold, numParties)
	chaincode, _ := node.NewMasterChainCode()
	index := uint32(0x0)

	devices := make([]node.Device, numParties)
	for i := uint32(0); i < numParties; i++ {
		device, pubkey := node.NewDevice(
			int(i),
			pkShares[i+1].Point,
			skShares[i+1].ShamirShare,
			pk,
			index,
			chaincode,
			nil,
		)
		devices = append(devices, node.Device{})
	}
	return devices
}
