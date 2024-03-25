package derivation

import (
	"errors"

	"bip32_threshold_wallet/node"
)

type NonHardDerivation struct {
	devices []node.Device
}

func (nhd NonHardDerivation) DeriveNonHardenedChild(childIdx uint32) (error, []node.Device) {
	// Execute DPSS between nd.devices and the new devices
	return errors.New("not implemented"), nil
}
