package node

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"

	"github.com/coinbase/kryptology/pkg/core/curves"
	"github.com/coinbase/kryptology/pkg/tecdsa/gg20/dealer"
)

const seed = "fffcf9f6f3f0edeae7e4e1dedbd8d5d2cfccc7a1c3c0bdbab7b4b1aeaba8a5a29f9c999693908d8a8784817e7b7875726f6c696663605d5a5754514e4b484542"

// GenSharedKey generates a threshold shared secret key.
// It outputs n public key shares, n secret key shares, and the global public key.
func GenSharedKey(t uint32, n uint32) (map[uint32]*dealer.PublicShare, map[uint32]*dealer.Share, *curves.Point) {
	k256, err := curves.K256().ToEllipticCurve()
	if err != nil {
		panic(err)
	}

	secret, _ := dealer.NewSecret(k256)
	pkEc, sharesMap, _ := dealer.NewDealerShares(k256, t, n, secret)
	pk, err := curves.K256().Point.Set(pkEc.X, pkEc.Y)
	if err != nil {
		panic(err)
	}
	pubSharesMap, _ := dealer.PreparePublicShares(sharesMap)

	// TODO: define logger
	fmt.Printf("Sharing scheme: Any %d from %d\n", t, n)
	fmt.Printf("Random secret: (%x)\n\n", secret)

	for i := range sharesMap {
		fmt.Printf("Share: %x\n", sharesMap[i].Bytes())
	}

	return pubSharesMap, sharesMap, &pk
}

func NewMasterChainCode() ([]byte, error) {
	// Generate key and chaincode
	hmac := hmac.New(sha512.New, []byte("Bitcoin seed"))
	seedBytes, _ := hex.DecodeString(seed)
	_, err := hmac.Write(seedBytes)
	if err != nil {
		return nil, err
	}
	intermediary := hmac.Sum(nil)

	// Split it into our key and chain code
	// keyBytes := intermediary[:32]
	chainCode := intermediary[32:]

	return chainCode, nil
}
