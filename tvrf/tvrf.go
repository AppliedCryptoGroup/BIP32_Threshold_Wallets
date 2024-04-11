package tvrf

import (
	"hash"
	"math/big"

	"github.com/coinbase/kryptology/pkg/core/curves"
	v1 "github.com/coinbase/kryptology/pkg/sharing/v1"
	"github.com/pkg/errors"
)

// The implementation of the DDH-based TVRF as proposed in https://eprint.iacr.org/2020/096.

type TVRF interface {
	// PEval computes the partial evaluation of the TVRF.
	PEval(m Message, sk SecretKeyShare, pk PublicKeyShare) (*PartialEvaluation, error)
	// Verify verifies the evaluation of the TVRF on the given message.
	Verify(eval Evaluation) bool
	// Combine combines at least t partial evaluations to compute the final evaluation of the TVRF.
	Combine(evals []*PartialEvaluation) (*Evaluation, error)
}

type DDHTVRF struct {
	t uint32
	n uint32

	curve *curves.Curve
	hash  hash.Hash
}

type Evaluation struct {
	Eval  curves.Point
	Proof []*PartialEvaluation
}

type PartialEvaluation struct {
	PubKeyShare PublicKeyShare
	// TODO: Replace with suitable types.
	Eval  curves.Point
	Proof *Proof
}

type PublicKey curves.Point

// PublicKeyShare represents a public key share, also containing an index.
type PublicKeyShare struct {
	Idx   uint32
	Value *curves.Point
}

// SecretKeyShare represents a secret key share.
type SecretKeyShare curves.Scalar

// TODO: Replace with suitable type.
type Message []byte

func NewDDHTVRF(t uint32, n uint32, curve *curves.Curve, hash hash.Hash) *DDHTVRF {
	return &DDHTVRF{
		t:     t,
		n:     n,
		curve: curve,
		hash:  hash,
	}
}

func (t *DDHTVRF) PEval(m Message, sk SecretKeyShare, pk PublicKeyShare) (*PartialEvaluation, error) {
	h := t.curve.Point.Hash(m)
	phi := h.Mul(sk)
	proof := t.proveEq(phi, m, sk, pk)
	eval := PartialEvaluation{
		PubKeyShare: pk,
		Eval:        phi,
		Proof:       proof,
	}

	return &eval, nil
}

func (t *DDHTVRF) Verify(eval Evaluation) bool {
	correctEvals := make([]*PartialEvaluation, 0)
	for _, e := range eval.Proof {
		if t.verifyEq(e.Eval, e.PubKeyShare, e.Proof) {
			correctEvals = append(correctEvals, e)
		}
	}

	if eval.Eval.Equal(t.combineEvaluations(correctEvals)) {
		return true
	} else {
		return false
	}
}

func (t *DDHTVRF) Combine(evals []*PartialEvaluation) (*Evaluation, error) {
	if len(evals) < int(t.t) {
		return nil, errors.New("not enough partial evaluations, need at least t evaluations to combine")
	}

	correctEvals := make([]*PartialEvaluation, 0)
	for _, e := range evals {
		if t.verifyEq(e.Eval, e.PubKeyShare, e.Proof) {
			correctEvals = append(correctEvals, e)
		}
	}
	if len(correctEvals) < int(t.t) {
		return nil, errors.New("not enough correct partial evaluations")
	}

	combinedEval := t.combineEvaluations(correctEvals)

	return &Evaluation{
		Eval:  combinedEval,
		Proof: correctEvals,
	}, nil
}

func (t *DDHTVRF) VerifyPartialEval(eval *PartialEvaluation) bool {
	return t.verifyEq(eval.Eval, eval.PubKeyShare, eval.Proof)
}

func (t *DDHTVRF) combineEvaluations(evals []*PartialEvaluation) curves.Point {
	indicesSet := make([]int, 0)
	for _, eval := range evals {
		indicesSet = append(indicesSet, int(eval.PubKeyShare.Idx))
	}

	combinedEval := t.curve.Point.Identity() // TODO does this correspond to 1?
	// Compute combinedEval = \prod eval_i^{\lambda_i}
	for _, eval := range evals {
		lambda := t.lagrangeCoefficient(int(eval.PubKeyShare.Idx), indicesSet)
		combinedEval = combinedEval.Add(eval.Eval.Mul(lambda))
	}

	return combinedEval
}

// lagrangeCoefficient computes the Lagrange coefficient for the given index at the 0 evaluation.
func (t *DDHTVRF) lagrangeCoefficient(idx int, indicesSet []int) curves.Scalar {
	// \prod_{k\in indicesSet\setminus 0} (Idx-k) / (0-k)
	lambda := t.curve.Scalar.One() // TODO does this correspond to 1?

	for _, k := range indicesSet {
		if k == 0 || k == idx {
			continue
		}

		// lambda = lambda * (Idx-k) / (0-k)
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

func ShamirShareToKeyPair(curve *curves.Curve, secretShare *v1.ShamirShare, pubShare *curves.EcPoint) (error, SecretKeyShare, *PublicKeyShare) {
	sk, err := curve.Scalar.SetBytes(secretShare.Value.Bytes())
	// FIXME: This seems to fail sometimes with error "invalid length"
	if err != nil {
		return errors.Wrap(err, "setting secret key scalar"), nil, nil
	}

	pkPoint, err := curve.Point.Set(pubShare.X, pubShare.Y)
	if err != nil {
		return errors.Wrap(err, "setting public key coords"), nil, nil

	}
	pk := &PublicKeyShare{
		Idx:   secretShare.Identifier,
		Value: &pkPoint,
	}

	return nil, sk, pk
}
