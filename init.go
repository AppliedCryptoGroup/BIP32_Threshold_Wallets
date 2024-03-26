package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strconv"

	"github.com/btcsuite/btcd/btcec"
	"github.com/coinbase/kryptology/pkg/core/curves"
	"github.com/coinbase/kryptology/pkg/tecdsa/gg20/dealer"

	"bip32_threshold_wallet/node"

	"go.dedis.ch/dela/crypto"
	"go.dedis.ch/dela/mino"
	"go.dedis.ch/dela/mino/minogrpc"
	"go.dedis.ch/dela/mino/router/tree"
	"go.dedis.ch/kyber/v3"
)

// CollectiveAuthority is a fake implementation of the cosi.CollectiveAuthority
// interface.
type CollectiveAuthority struct {
	crypto.CollectiveAuthority
	addrs   []mino.Address
	pubkeys []kyber.Point
}

// NewAuthority returns a new collective authority of n members with new signers
// generated by g.
func NewAuthority(addrs []mino.Address, pubkeys []kyber.Point) CollectiveAuthority {

	return CollectiveAuthority{
		pubkeys: pubkeys,
		addrs:   addrs,
	}
}

// leave the parameters here for future use
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

// the function generates big ints which might be needed for threshold ecdsa generation.
func B10(s string) *big.Int {
	x, ok := new(big.Int).SetString(s, 10)
	if !ok {
		panic("Couldn't derive big.Int from string")
	}
	return x
}

func InitDevices(t int, n int) (CollectiveAuthority, []node.Device) {
	minos := make([]mino.Mino, n)
	devices := make([]node.Device, n)
	addrs := make([]mino.Address, n)
	for i := 0; i < n; i++ {
		addr := minogrpc.ParseAddress("127.0.0.1", 0)

		minogrpc, _ := minogrpc.NewMinogrpc(addr, nil, tree.NewRouter(minogrpc.NewAddressFactory()))
		// fmt.Println(err)

		defer minogrpc.GracefulStop()

		minos[i] = minogrpc
		addrs[i] = minogrpc.GetAddress()
	}

	pubkeys := make([]kyber.Point, len(minos))

	privShares, pubKey := GenSk(uint32(t), uint32(n))

	for i, mino := range minos {
		for _, m := range minos {
			mino.(*minogrpc.Minogrpc).GetCertificateStore().Store(m.GetAddress(),
				m.(*minogrpc.Minogrpc).GetCertificateChain())
		}

		// privkey := suite.Scalar().Pick(suite.RandomStream())
		// pubkey := suite.Point().Mul(privkey, nil)

		device, pubkey := node.NewDevice(i, privShares[uint32(i)+1], pubKey, mino.(*minogrpc.Minogrpc))
		pubkeys[i] = pubkey
		devices[i] = device
	}

	Authority := NewAuthority(addrs, pubkeys)

	field := curves.NewField(btcec.S256().Params().N)
	rnd := rand.Reader
	rho, _ := field.RandomElement(rnd)
	devices[0].RandSk(*rho)

	return Authority, devices
}

func GenSk(tshare uint32, nshare uint32) (map[uint32]*dealer.Share, *curves.EcPoint) {

	// tshare := uint32(2)
	// nshare := uint32(3)u
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

	return sharesMap, pk
}
