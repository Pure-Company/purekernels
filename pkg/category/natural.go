// pkg/category/natural.go
package category

import (
	"github.com/vinodhalaharvi/purekernels/pkg/either"
	"github.com/vinodhalaharvi/purekernels/pkg/monoid"
	"github.com/vinodhalaharvi/purekernels/pkg/result"
)

// NaturalTransformation represents a transformation between functors
// A natural transformation F ~> G is a family of morphisms that
// commutes with the functorial structure
type NaturalTransformation[F, G, A any] struct {
	transform func(F) G
}

// NewNat creates a natural transformation
func NewNat[F, G, A any](f func(F) G) NaturalTransformation[F, G, A] {
	return NaturalTransformation[F, G, A]{transform: f}
}

// Apply applies the natural transformation
func (n NaturalTransformation[F, G, A]) Apply(fa F) G {
	return n.transform(fa)
}

// Concrete natural transformations

// OptionToSlice transforms Option to Slice
func OptionToSlice[A any](opt monoid.Option[A]) []A {
	if val, ok := opt.Get(); ok {
		return []A{val}
	}
	return []A{}
}

// ResultToEither transforms Result to Either
func ResultToEither[A any](r result.Result[A]) either.Either[error, A] {
	if val, err := r.Get(); err == nil {
		return either.Right[error, A](val)
	} else {
		return either.Left[error, A](err)
	}
}

// EitherToResult transforms Either[error, A] to Result[A]
func EitherToResult[A any](e either.Either[error, A]) result.Result[A] {
	if e.IsRight() {
		return result.Ok(e.GetRight())
	}
	return result.Err[A](e.GetLeft())
}

// EitherToOption transforms Either to Option (discards Left)
func EitherToOption[L, R any](e either.Either[L, R]) monoid.Option[R] {
	return e.ToOption()
}

// ResultToOption transforms Result to Option (discards error)
func ResultToOption[A any](r result.Result[A]) monoid.Option[A] {
	return r.ToOption()
}

// SliceToOption transforms non-empty slice to Some, empty to None
func SliceToOption[A any](xs []A) monoid.Option[A] {
	if len(xs) > 0 {
		return monoid.Some(xs[0])
	}
	return monoid.None[A]()
}

// NatOptionSlice is a concrete natural transformation Option ~> Slice
type NatOptionSlice[A any] struct{}

func (NatOptionSlice[A]) Transform(opt monoid.Option[A]) []A {
	return OptionToSlice(opt)
}

// NatResultEither is a concrete natural transformation Result ~> Either[error, *]
type NatResultEither[A any] struct{}

func (NatResultEither[A]) Transform(r result.Result[A]) either.Either[error, A] {
	return ResultToEither(r)
}
