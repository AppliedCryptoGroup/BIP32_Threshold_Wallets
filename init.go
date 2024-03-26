package main

import (
	"fmt"
	"math/big"
	"os"
	"strconv"

	"github.com/btcsuite/btcd/btcec"
	"github.com/coinbase/kryptology/pkg/tecdsa/gg20/dealer"
)

func getParams(msg *string, t, n *uint32) {
	argCount := len(os.Args[1:])

	if argCount > 0 {
		*msg = os.Args[1]

	}
	if argCount > 1 {
		val, _ := strconv.Atoi(os.Args[2])
		*t = uint32(val)
	}
	if argCount > 2 {
		val, _ := strconv.Atoi(os.Args[3])
		*n = uint32(val)
	}

}

func B10(s string) *big.Int {
	x, ok := new(big.Int).SetString(s, 10)
	if !ok {
		panic("Couldn't derive big.Int from string")
	}
	return x
}

func initSecretSharing() {

	tshare := uint32(2)
	nshare := uint32(3)
	msg := "Hello"

	getParams(&msg, &tshare, &nshare)

	k256 := btcec.S256()

	ikm, _ := dealer.NewSecret(k256)

	pk, sharesMap, _ := dealer.NewDealerShares(k256, tshare, nshare, ikm)

	fmt.Printf("Message: %s\n", msg)
	fmt.Printf("Sharing scheme: Any %d from %d\n", tshare, nshare)
	fmt.Printf("Random secret: (%x)\n\n", ikm)
	fmt.Printf("Public key: (%s %s)\n\n", pk.X, pk.Y)

	// for len(sharesMap) > int(tshare) {
	// 	delete(sharesMap, uint32(len(sharesMap)))
	// }
	// pubSharesMap, _ := dealer.PreparePublicShares(sharesMap)

	for i := range sharesMap {
		fmt.Printf("Share: %x\n", sharesMap[i].Bytes())
	}
}
