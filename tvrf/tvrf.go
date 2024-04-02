package tvrf

import (
	"errors"
	"hash"
	"math/big"

	"github.com/coinbase/kryptology/pkg/core/curves"
)

// The implementation of the DDH-based TVRF as proposed in https://eprint.iacr.org/2020/096.

type TVRF interface {
	// PEval computes the partial evaluation of the TVRF.
	PEval(m Message, sk SecretKeyShare, pk PublicKeyShare) (*PartialEvaluation, error)
	// Verify verifies the evaluation of the TVRF on the given message.
	Verify(pk PublicKey, m Message, eval Evaluation) (bool, error)
	// Combine combines at least t partial evaluations to compute the final evaluation of the TVRF.
	Combine(evals []PartialEvaluation) (*Evaluation, error)
}

type DDHTVRF struct {
	t uint32
	n uint32
	//publicKeys PublicKeys

	curve *curves.Curve
	hash  hash.Hash
}

type Evaluation struct {
	Eval  curves.Point
	Proof []PartialEvaluation
}

type PartialEvaluation struct {
	pk PublicKeyShare
	// TODO: Replace with suitable types.
	Eval  curves.Point
	Proof []byte
}

type PublicKey *curves.Point

// PublicKeyShare represents a public key share, also containing an index.
type PublicKeyShare struct {
	idx   uint32
	value *curves.Point
}

// SecretKeyShare represents a secret key share.
type SecretKeyShare *curves.Scalar

// TODO: Replace with suitable type.
type Message []byte

func NewDDHTVRF(t uint32, n uint32, curve *curves.Curve, hash hash.Hash) (*DDHTVRF, error) {
	return &DDHTVRF{
		t:     t,
		n:     n,
		curve: curve,
		hash:  hash,
	}, nil
}

func (t *DDHTVRF) PEval(m Message, sk SecretKeyShare, pk PublicKeyShare) (*PartialEvaluation, error) {
	h := t.curve.Point.Hash(m)
	phi := h.Mul(*sk)
	// TODO: proof := t.ZKP.Prove(phi, sk, pk)
	eval := PartialEvaluation{
		pk:    pk,
		Eval:  phi,
		Proof: nil,
	}

	return &eval, nil
}

func (t *DDHTVRF) Verify(pk PublicKey, m Message, eval Evaluation) (bool, error) {
	return false, errors.New("not implemented")
}

func (t *DDHTVRF) Combine(evals []PartialEvaluation) (*Evaluation, error) {
	if len(evals) < int(t.t) {
		return nil, errors.New("not enough partial evaluations, need at least t evaluations to combine")
	}

	correctEvals := make([]PartialEvaluation, 0)
	for _, eval := range evals {
		// TODO: Check if t.ZKP.Verify(eval.Proof, eval.Eval, eval.pk)
		correctEvals = append(correctEvals, eval)
	}
	if len(correctEvals) < int(t.t) {
		return nil, errors.New("not enough correct partial evaluations")
	}

	combinedEval := t.combineEvaluations(correctEvals)

	return &Evaluation{
		Eval:  combinedEval,
		Proof: nil,
	}, nil
}

func (t *DDHTVRF) combineEvaluations(evals []PartialEvaluation) curves.Point {
	indicesSet := make([]int, 0)
	for _, eval := range evals {
		indicesSet = append(indicesSet, int(eval.pk.idx))
	}

	combinedEval := t.curve.Point.Identity() // TODO does this correspond to 1?
	// Compute combinedEval = \prod eval_i^{\lambda_i}
	for _, eval := range evals {
		lambda := t.lagrangeCoefficient(int(eval.pk.idx), indicesSet)
		combinedEval = combinedEval.Add(eval.Eval.Mul(lambda))
	}

	return combinedEval
}

// lagrangeCoefficient computes the Lagrange coefficient for the given index at the 0 evaluation.
func (t *DDHTVRF) lagrangeCoefficient(idx int, indicesSet []int) curves.Scalar {
	// \prod_{k\in indicesSet\setminus 0} (idx-k) / (0-k)
	lambda := t.curve.Scalar.One() // TODO does this correspond to 1?

	for _, k := range indicesSet {
		// Should never happen as 0 is not in the set.
		if k == 0 {
			continue
		}

		// lambda = lambda * (idx-k) / (0-k)
		numerator := big.NewInt(int64(idx - k))
		numeratorScalar, err := t.curve.Scalar.SetBigInt(numerator)
		if err != nil {
			panic(err)
		}
		denominator := big.NewInt(int64(0 - k))
		denominatorScalar, err := t.curve.Scalar.SetBigInt(denominator)
		if err != nil {
			panic(err)
		}
		lambda = lambda.Mul(numeratorScalar.Div(denominatorScalar))
	}

	return lambda
}
