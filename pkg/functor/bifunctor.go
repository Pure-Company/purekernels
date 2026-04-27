package functor

import (
	"github.com/Pure-Company/purekernels/pkg/either"
	"github.com/Pure-Company/purekernels/pkg/monoid"
	"github.com/Pure-Company/purekernels/pkg/pair"
)

// BimapEither applies functions to both sides of Either
func BimapEither[L1, L2, R1, R2 any](
	e either.Either[L1, R1],
	fl func(L1) L2,
	fr func(R1) R2,
) either.Either[L2, R2] {
	if e.IsLeft() {
		return either.Left[L2, R2](fl(e.GetLeft()))
	}
	return either.Right[L2, R2](fr(e.GetRight()))
}

// BimapPair applies functions to both components of Pair
func BimapPair[A1, A2, B1, B2 any](
	p pair.Pair[A1, B1],
	fa func(A1) A2,
	fb func(B1) B2,
) pair.Pair[A2, B2] {
	return pair.NewPair(fa(p.First), fb(p.Second))
}

// BimapValidation applies functions to both sides of Validation
func BimapValidation[E1, E2, A1, A2 any](
	v either.Validation[E1, A1],
	fe func(E1) E2,
	fa func(A1) A2,
	m2 monoid.Monoid[E2],
) either.Validation[E2, A2] {
	val, errs, isValid := v.Get()
	if isValid {
		return either.Valid[E2, A2](m2, fa(val))
	}
	return either.Invalid[E2, A2](m2, fe(errs))
}
