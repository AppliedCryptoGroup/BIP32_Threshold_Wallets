package derivation

import (
	"errors"

	"bip32_threshold_wallet/node"
)

type TVRFDerivation struct {
	devices []node.Device
}

func (td TVRFDerivation) DeriveNonHardenedChild(childIdx uint32) (error, []node.Device) {
	nonHardDerivation := NonHardDerivation{devices: td.devices}
	return nonHardDerivation.DeriveNonHardenedChild(childIdx)
}

func (td TVRFDerivation) DeriveHardenedChild(childIdx uint32) (error, *node.Node) {
	/*
		for i := 0; i < n; i++ {
			r_i = evaluateTVRF(device[i])
		}
		seed = combineTVRFShares(r_1, ..., r_n)
		-> mock sending seed to the child node
		sk,pk= ECDSAKeyGen(seed)
		node := Node{sk: sk, pk: pk}
	*/
	return errors.New("not implemented"), nil
}
