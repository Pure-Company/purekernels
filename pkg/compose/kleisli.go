// pkg/compose/kleisli.go
package compose

import (
	"github.com/vinodhalaharvi/purekernels/pkg/effect"
	"github.com/vinodhalaharvi/purekernels/pkg/either"
	"github.com/vinodhalaharvi/purekernels/pkg/monoid"
	"github.com/vinodhalaharvi/purekernels/pkg/result"
)

// KleisliEither composes A -> Either[L,B] with B -> Either[L,C]
func KleisliEither[L, A, B, C any](
	f func(A) either.Either[L, B],
	g func(B) either.Either[L, C],
) func(A) either.Either[L, C] {
	return func(a A) either.Either[L, C] {
		return either.FlatMapEither(f(a), g)
	}
}

// KleisliResult composes A -> Result[B] with B -> Result[C]
func KleisliResult[A, B, C any](
	f func(A) result.Result[B],
	g func(B) result.Result[C],
) func(A) result.Result[C] {
	return func(a A) result.Result[C] {
		return result.FlatMapResult(f(a), g)
	}
}

// KleisliOption composes A -> Option[B] with B -> Option[C]
func KleisliOption[A, B, C any](
	f func(A) monoid.Option[B],
	g func(B) monoid.Option[C],
) func(A) monoid.Option[C] {
	return func(a A) monoid.Option[C] {
		return monoid.FlatMapOption(f(a), g)
	}
}

// KleisliReader composes Reader arrows
func KleisliReader[R, A, B, C any](
	f func(A) effect.Reader[R, B],
	g func(B) effect.Reader[R, C],
) func(A) effect.Reader[R, C] {
	return func(a A) effect.Reader[R, C] {
		return func(r R) C {
			b := f(a)(r)
			return g(b)(r)
		}
	}
}

// KleisliState composes State arrows
func KleisliState[S, A, B, C any](
	f func(A) effect.State[S, B],
	g func(B) effect.State[S, C],
) func(A) effect.State[S, C] {
	return func(a A) effect.State[S, C] {
		return func(s S) (C, S) {
			b, s2 := f(a)(s)
			return g(b)(s2)
		}
	}
}
