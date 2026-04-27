// pkg/category/kleisli.go
// Package category provides categorical morphisms and arrows
package category

import (
	"github.com/Pure-Company/purekernels/pkg/comonad"
	"github.com/Pure-Company/purekernels/pkg/effect"
	"github.com/Pure-Company/purekernels/pkg/either"
	"github.com/Pure-Company/purekernels/pkg/monoid"
	"github.com/Pure-Company/purekernels/pkg/result"
)

// Kleisli represents a morphism A -> M[B] in the Kleisli category
type Kleisli[M, A, B any] struct {
	run func(A) M
}

// NewKleisli creates a Kleisli arrow
func NewKleisli[M, A, B any](f func(A) M) Kleisli[M, A, B] {
	return Kleisli[M, A, B]{run: f}
}

// Run executes the Kleisli arrow
func (k Kleisli[M, A, B]) Run(a A) M {
	return k.run(a)
}

// KleisliEither is a concrete Kleisli arrow for Either
type KleisliEither[L, A, B any] struct {
	run func(A) either.Either[L, B]
}

// NewKleisliEither creates a Kleisli arrow for Either
func NewKleisliEither[L, A, B any](f func(A) either.Either[L, B]) KleisliEither[L, A, B] {
	return KleisliEither[L, A, B]{run: f}
}

// Run executes the arrow
func (k KleisliEither[L, A, B]) Run(a A) either.Either[L, B] {
	return k.run(a)
}

// ComposeKleisliEither composes two Kleisli arrows: A -> Either[L,B] >>> B -> Either[L,C]
func ComposeKleisliEither[L, A, B, C any](
	first KleisliEither[L, A, B],
	second KleisliEither[L, B, C],
) KleisliEither[L, A, C] {
	return KleisliEither[L, A, C]{
		run: func(a A) either.Either[L, C] {
			return either.FlatMapEither(first.run(a), second.run)
		},
	}
}

// KleisliResult is a concrete Kleisli arrow for Result
type KleisliResult[A, B any] struct {
	run func(A) result.Result[B]
}

// NewKleisliResult creates a Kleisli arrow for Result
func NewKleisliResult[A, B any](f func(A) result.Result[B]) KleisliResult[A, B] {
	return KleisliResult[A, B]{run: f}
}

// Run executes the arrow
func (k KleisliResult[A, B]) Run(a A) result.Result[B] {
	return k.run(a)
}

// ComposeKleisliResult composes two Kleisli arrows
func ComposeKleisliResult[A, B, C any](
	first KleisliResult[A, B],
	second KleisliResult[B, C],
) KleisliResult[A, C] {
	return KleisliResult[A, C]{
		run: func(a A) result.Result[C] {
			return result.FlatMapResult(first.run(a), second.run)
		},
	}
}

// KleisliOption is a concrete Kleisli arrow for Option
type KleisliOption[A, B any] struct {
	run func(A) monoid.Option[B]
}

// NewKleisliOption creates a Kleisli arrow for Option
func NewKleisliOption[A, B any](f func(A) monoid.Option[B]) KleisliOption[A, B] {
	return KleisliOption[A, B]{run: f}
}

// Run executes the arrow
func (k KleisliOption[A, B]) Run(a A) monoid.Option[B] {
	return k.run(a)
}

// ComposeKleisliOption composes two Kleisli arrows
func ComposeKleisliOption[A, B, C any](
	first KleisliOption[A, B],
	second KleisliOption[B, C],
) KleisliOption[A, C] {
	return KleisliOption[A, C]{
		run: func(a A) monoid.Option[C] {
			return monoid.FlatMapOption(first.run(a), second.run)
		},
	}
}

// KleisliState is a concrete Kleisli arrow for State
type KleisliState[S, A, B any] struct {
	run func(A) effect.State[S, B]
}

// NewKleisliState creates a Kleisli arrow for State
func NewKleisliState[S, A, B any](f func(A) effect.State[S, B]) KleisliState[S, A, B] {
	return KleisliState[S, A, B]{run: f}
}

// Run executes the arrow
func (k KleisliState[S, A, B]) Run(a A) effect.State[S, B] {
	return k.run(a)
}

// ComposeKleisliState composes two Kleisli arrows
func ComposeKleisliState[S, A, B, C any](
	first KleisliState[S, A, B],
	second KleisliState[S, B, C],
) KleisliState[S, A, C] {
	return KleisliState[S, A, C]{
		run: func(a A) effect.State[S, C] {
			return func(s S) (C, S) {
				b, s2 := first.run(a)(s)
				return second.run(b)(s2)
			}
		},
	}
}

// CoKleisli represents a morphism W[A] -> B in the CoKleisli category
type CoKleisli[W, A, B any] struct {
	run func(W) B
}

// NewCoKleisli creates a CoKleisli arrow
func NewCoKleisli[W, A, B any](f func(W) B) CoKleisli[W, A, B] {
	return CoKleisli[W, A, B]{run: f}
}

// Run executes the CoKleisli arrow
func (c CoKleisli[W, A, B]) Run(w W) B {
	return c.run(w)
}

// CoKleisliStore is a concrete CoKleisli arrow for Store
type CoKleisliStore[S, A, B any] struct {
	run func(comonad.Store[S, A]) B
}

// NewCoKleisliStore creates a CoKleisli arrow for Store
func NewCoKleisliStore[S, A, B any](
	f func(comonad.Store[S, A]) B,
) CoKleisliStore[S, A, B] {
	return CoKleisliStore[S, A, B]{run: f}
}

// Run executes the arrow
func (c CoKleisliStore[S, A, B]) Run(w comonad.Store[S, A]) B {
	return c.run(w)
}

// CoKleisliEnv is a concrete CoKleisli arrow for Env
type CoKleisliEnv[E, A, B any] struct {
	run func(comonad.Env[E, A]) B
}

// NewCoKleisliEnv creates a CoKleisli arrow for Env
func NewCoKleisliEnv[E, A, B any](
	f func(comonad.Env[E, A]) B,
) CoKleisliEnv[E, A, B] {
	return CoKleisliEnv[E, A, B]{run: f}
}

// Run executes the arrow
func (c CoKleisliEnv[E, A, B]) Run(w comonad.Env[E, A]) B {
	return c.run(w)
}

// IdentityKleisliEither returns the identity Kleisli arrow
func IdentityKleisliEither[L, A any]() KleisliEither[L, A, A] {
	return NewKleisliEither[L, A, A](func(a A) either.Either[L, A] {
		return either.Right[L, A](a)
	})
}

// IdentityCoKleisliStore returns the identity CoKleisli arrow (extract)
func IdentityCoKleisliStore[S, A any]() CoKleisliStore[S, A, A] {
	return NewCoKleisliStore[S, A, A](func(w comonad.Store[S, A]) A {
		return w.Extract()
	})
}
