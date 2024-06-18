package derivation

import (
	"crypto/rand"
	"fmt"

	"github.com/pkg/errors"
	"github.com/tyler-smith/go-bip32"
)

type StandardBIP32Derivation struct {
	masterKey *bip32.Key
}

func NewStandardBIP32Derivation() (error, *StandardBIP32Derivation) {
	seed, err := genSeed()
	if err != nil {
		return err, nil
	}

	key, err := bip32.NewMasterKey(seed)
	if err != nil {
		return errors.Wrap(err, "failed to generate master key"), nil
	}

	return nil, &StandardBIP32Derivation{masterKey: key}
}

func (s *StandardBIP32Derivation) DeriveNonHardenedChild(childIdx uint32) (*bip32.Key, error) {
	// TODO: Check if the index is hardened.

	key, err := s.masterKey.NewChildKey(childIdx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to derive child key")
	}

	return key, nil
}

func (s *StandardBIP32Derivation) DeriveHardenedChild(childIdx uint32) (*bip32.Key, error) {
	// TODO: Check if the index is non-hardened.

	key, err := s.masterKey.NewChildKey(childIdx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to derive child key")
	}

	return key, nil
}

// Generates 32 byte seed for BIP32 derivation
func genSeed() ([]byte, error) {
	seed := make([]byte, 32)
	_, err := rand.Read(seed)
	if err != nil {
		return nil, fmt.Errorf("failed to generate seed: %v", err)
	}
	return seed, nil
}
