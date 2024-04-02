package tvrf

import (
	"errors"
	"hash"

	"github.com/coinbase/kryptology/pkg/core/curves"
)

// The implementation of the DDH-based TVRF as proposed in https://eprint.iacr.org/2020/096.

type TVRF struct {
	t uint32
	n uint32
	//publicKeys PublicKeys

	curve *curves.Curve
	hash  hash.Hash
}

type Evaluation struct {
	// TODO: Replace with suitable types.
	Eval  curves.Point
	Proof []byte
}

type PublicKey *curves.Point // TODO: Should we differentiate between public keys and public key shares?
type SecretKey *curves.Scalar

// TODO: Replace with suitable type.
type Message []byte

func NewTVRF(t uint32, n uint32, curve *curves.Curve, hash hash.Hash) (*TVRF, error) {
	return &TVRF{
		t:     t,
		n:     n,
		curve: curve,
		hash:  hash,
	}, nil
}

func (t *TVRF) PEval(m Message, sk SecretKey, pk PublicKey) (*Evaluation, error) {
	h := t.curve.Point.Hash(m)
	phi := h.Mul(*sk)
	// TODO: proof := t.ZKP.Prove(phi, sk, pk)
	eval := Evaluation{
		Eval:  phi,
		Proof: nil,
	}

	return &eval, nil
}

func (t *TVRF) Verify(pk PublicKey, m Message, eval Evaluation) (bool, error) {
	return false, errors.New("not implemented")
}

func (t *TVRF) Combine(evals []Evaluation) (Evaluation, error) {
	return Evaluation{}, errors.New("not implemented")
}
