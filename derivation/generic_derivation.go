package derivation

import (
	"errors"

	"bip32_threshold_wallet/node"
)

type GenericDerivation struct {
	devices []node.Device
}

func (gd GenericDerivation) DeriveNonHardenedChild(childIdx uint32) (error, []node.Device) {
	nonHardDerivation := NonHardDerivation{devices: gd.devices}
	return nonHardDerivation.DeriveNonHardenedChild(childIdx)
}

func (gd GenericDerivation) DeriveHardenedChild(childIdx uint32) (error, *node.Node) {
	return errors.New("not implemented"), nil
}
