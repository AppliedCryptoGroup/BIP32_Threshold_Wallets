package derivation_test

import (
	sha2562 "crypto/sha256"
	"testing"

	"github.com/coinbase/kryptology/pkg/core/curves"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"bip32_threshold_wallet/derivation"
	"bip32_threshold_wallet/tvrf"
	"bip32_threshold_wallet/utils"
)

var (
	curve      = curves.K256()
	sha256     = sha2562.New()
	threshold  = uint32(3)
	numParties = uint32(5)
)

func TestNewTVRFDerivation(t *testing.T) {
	devices := utils.CreateDevices(threshold, numParties)
	ddhTvrf := tvrf.NewDDHTVRF(threshold, numParties, curve, sha256, true)
	deriv := derivation.NewTVRFDerivation(curve, devices, ddhTvrf, true)

	err, childNode1 := deriv.DeriveHardenedChild(1)
	require.NoError(t, err)
	require.NotNil(t, childNode1)

	err, childNode1Clone := deriv.DeriveHardenedChild(1)
	require.NoError(t, err)
	require.NotNil(t, childNode1Clone)

	assert.Truef(t, (*childNode1.PublicKey).Equal(*childNode1Clone.PublicKey), "Public keys should be the same")

	err, childNode2 := deriv.DeriveHardenedChild(2)
	require.NoError(t, err)
	require.NotNil(t, childNode2)

	assert.Falsef(t, (*childNode1.PublicKey).Equal(*childNode2.PublicKey), "Public keys should be different")
}
