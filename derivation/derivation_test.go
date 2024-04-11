package derivation_test

import (
	"testing"

	"github.com/coinbase/kryptology/pkg/core/curves"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/sha3"

	"bip32_threshold_wallet/derivation"
	"bip32_threshold_wallet/node"
	"bip32_threshold_wallet/tvrf"
)

var (
	curve      = curves.K256()
	sha256     = sha3.New256()
	threshold  = uint32(3)
	numParties = uint32(5)
)

func TestNewTVRFDerivation(t *testing.T) {
	devices := createDevices()
	ddhTvrf := tvrf.NewDDHTVRF(threshold, numParties, curve, sha256)
	deriv := derivation.NewTVRFDerivation(curve, devices, ddhTvrf, true)

	err, childNode1 := deriv.DeriveHardenedChild(1)
	assert.NoError(t, err)
	assert.NotNil(t, childNode1)

	err, childNode1Clone := deriv.DeriveHardenedChild(1)
	assert.NoError(t, err)
	assert.NotNil(t, childNode1Clone)

	assert.Truef(t, (*childNode1.PublicKey).Equal(*childNode1Clone.PublicKey), "Public keys should be the same")

	err, childNode2 := deriv.DeriveHardenedChild(2)
	assert.NoError(t, err)
	assert.NotNil(t, childNode2)

	assert.Falsef(t, (*childNode1.PublicKey).Equal(*childNode2.PublicKey), "Public keys should be different")
}

func createDevices() []node.Device {
	pkShares, skShares, pk := node.GenSharedKey(threshold, numParties)
	chaincode, _ := node.NewMasterChainCode()
	index := uint32(0x0)

	devices := make([]node.Device, numParties)
	for i := uint32(0); i < numParties; i++ {
		device, _ := node.NewDevice(
			int(i),
			pkShares[i+1].Point,
			skShares[i+1].ShamirShare,
			pk,
			index,
			chaincode,
			nil,
		)
		devices[i] = device
	}
	return devices
}
