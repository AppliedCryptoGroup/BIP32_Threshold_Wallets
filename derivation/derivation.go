package derivation

import (
	"github.com/tyler-smith/go-bip32"

	"bip32_threshold_wallet/node"
)

/*
Derivation methods:
GenericDerivation = DeriveNonHardenedChild with DPSS and DeriveHardenedChild with MPC-Hash
TVRFDerivation = DeriveNonHardenedChild with DPSS and DeriveHardenedChild with a TVRF
TVRFDerivationOpt = DeriveNonHardenedChild with DPSS while reusing the key for DeriveHardenedChild with the TVRF
*/

type ThresholdDerivation interface {
	// DeriveNonHardenedChild derives a non-hardened child from the current node, which is shared among the devices.
	DeriveNonHardenedChild(childIdx uint32) ([]node.Device, error)

	// DeriveHardenedChild derives a hardened child from the current node, which is shared among the devices.
	DeriveHardenedChild(childIdx uint32) (*node.Node, error)
}

type StandardDerivation interface {
	// DeriveNonHardenedChild derives a non-hardened child from the current key.
	DeriveNonHardenedChild(childIdx uint32) (*bip32.Key, error)

	// DeriveHardenedChild derives a hardened child from the current key.
	DeriveHardenedChild(childIdx uint32) (*bip32.Key, error)
}
