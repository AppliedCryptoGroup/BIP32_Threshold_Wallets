package tvrf

import (
	"crypto/hmac"
	"crypto/rand"

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
	pkMar, _ := pointMarshalBinary(*pk.Value)
	com1Mar, _ := pointMarshalBinary(com1)
	com2Mar, _ := pointMarshalBinary(com2)
	marshaledValue = append(marshaledValue, phiMar...)
	marshaledValue = append(marshaledValue, pkMar...)
	marshaledValue = append(marshaledValue, com1Mar...)
	marshaledValue = append(marshaledValue, com2Mar...)
	ch := t.curve.Scalar.Hash(marshaledValue)

	res := t.curve.Scalar.One()

	res = res.Mul(r).Sub(ch.Mul(sk))

	return &Proof{res, ch, g}
}

func (t *DDHTVRF) VerifyEq(phi curves.Point, sk SecretKeyShare, pk curves.Point, proof *Proof) bool {
	res := proof.Res
	ch := proof.Ch
	g := proof.g
	rG := g.Mul(res)
	rH := t.curve.ScalarBaseMult(res)
	cxG := phi.Mul(ch)
	cxH := pk.Mul(ch)
	R := rG.Add(cxG)
	Rp := rH.Add(cxH)

	var marshaledValue []byte
	phiMar, _ := pointMarshalBinary(phi)
	pkMar, _ := pointMarshalBinary(pk)
	com1Mar, _ := pointMarshalBinary(R)
	com2Mar, _ := pointMarshalBinary(Rp)
	marshaledValue = append(marshaledValue, phiMar...)
	marshaledValue = append(marshaledValue, pkMar...)
	marshaledValue = append(marshaledValue, com1Mar...)
	marshaledValue = append(marshaledValue, com2Mar...)
	chp := t.curve.Scalar.Hash(marshaledValue)

	return hmac.Equal(chp.Bytes(), ch.Bytes())
}

// Adopted directly from kryptology/pkg/core/curves
func pointMarshalBinary(point curves.Point) ([]byte, error) {
	// Always stores points in compressed form
	// The first bytes are the p256 name
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
