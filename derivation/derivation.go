package derivation

import (
	"bip32_threshold_wallet/node"
)

type Derivation interface {
	// DeriveNonHardenedChild derives a non-hardened child from the current node, which is shared among the devices.
	DeriveNonHardenedChild(childIdx uint32) (error, []node.Device)

	// DeriveHardenedChild derives a hardened child from the current node, which is shared among the devices.
	DeriveHardenedChild(childIdx uint32) (error, *node.Node)
}

/*
Derivation methods:
GenericDerivation = DeriveNonHardenedChild with DPSS and DeriveHardenedChild with MPC-Hash
TVRFDerivation = DeriveNonHardenedChild with DPSS and DeriveHardenedChild with a TVRF
TVRFDerivationOpt = DeriveNonHardenedChild with DPSS while reusing the key for DeriveHardenedChild with the TVRF
*/
