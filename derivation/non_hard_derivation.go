package derivation

import (
	"errors"

	"bip32_threshold_wallet/node"
)

type NonHardDerivation struct {
	devices []node.Device
}

func (nhd NonHardDerivation) DeriveNonHardenedChild(childIdx uint32) ([]node.Device, error) {
	// Execute DPSS between nd.devices and the new devices
	return nil, errors.New("not implemented")
}
