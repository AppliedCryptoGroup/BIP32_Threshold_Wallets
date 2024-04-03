package tvrf

import (
	"crypto/rand"
	"fmt"

	"github.com/coinbase/kryptology/pkg/core/curves"
)

const scalarBytes = 32

type Proof struct {
	Res curves.Scalar // response value
	Ch  curves.Scalar // hash of intermediate proof values to streamline equality checks
	g   curves.Point  // h(m)
}

// dleq: log_{g}(g^x) == log_{h}(h^x)
// g = hash(m), x = priKeyShare, h = base point
// g^x = phi, h^x = pk,
func (t *DDHTVRF) ProveEq(phi curves.Point, m Message, sk SecretKeyShare, pk PublicKeyShare) *Proof {
	g := t.curve.Point.Hash(m)
	r := t.curve.Scalar.Random(rand.Reader)
	com1 := g.Mul(r)
	com2 := t.curve.ScalarBaseMult(r)
	var marshaledValue []byte
	phiMar, _ := pointMarshalBinary(phi)
	pkMar, _ := pointMarshalBinary(*pk.value)
	com1Mar, _ := pointMarshalBinary(com1)
	com2Mar, _ := pointMarshalBinary(com2)
	marshaledValue = append(marshaledValue, phiMar...)
	marshaledValue = append(marshaledValue, pkMar...)
	marshaledValue = append(marshaledValue, com1Mar...)
	marshaledValue = append(marshaledValue, com2Mar...)
	ch := t.curve.Scalar.Hash(marshaledValue)

	res := t.curve.Scalar.One()

	res = res.Mul(r).Sub(ch.Mul(*sk))

	return &Proof{res, ch, g}
}

// func (t *DDHTVRF) VerifyEq(phi curves.Point, sk SecretKeyShare, pk PublicKeyShare, proof *Proof){
// 	res := proof.Res
// 	ch := proof.Ch
// 	g := proof.g
// 	com1 := g.Mul(res)
// }

func pointMarshalBinary(point curves.Point) ([]byte, error) {
	// Always stores points in compressed form
	// The first bytes are the curve name
	// separated by a colon followed by the compressed point
	// bytes
	t := point.ToAffineCompressed()
	name := []byte(point.CurveName())
	output := make([]byte, len(name)+1+len(t))
	copy(output[:len(name)], name)
	output[len(name)] = byte(':')
	copy(output[len(output)-len(t):], t)
	return output, nil
}

func pointUnmarshalBinary(input []byte) (curves.Point, error) {
	if len(input) < scalarBytes+1+len("secp256k1") {
		return nil, fmt.Errorf("invalid byte sequence")
	}
	sep := byte(':')
	i := 0
	for ; i < len(input); i++ {
		if input[i] == sep {
			break
		}
	}
	curve := curves.K256()
	if curve == nil {
		return nil, fmt.Errorf("unrecognized curve")
	}
	return curve.Point.FromAffineCompressed(input[i+1:])
}
