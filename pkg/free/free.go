// pkg/free/free.go
// Package free provides free constructions for generating lawful functors
package free

import "github.com/Pure-Company/purekernels/pkg/monoid"

// Free[F, A] is the free monad over functor F
// It represents a computation that can:
//   - Return a pure value (Pure)
//   - Suspend in a functor layer (Free)
//
// This lets you build ASTs/programs that can be interpreted multiple ways
type Free[F, A any] struct {
	value   *A
	suspend *F
	isPure  bool
}

// Pure creates a Free monad from a value
func Pure[F, A any](a A) Free[F, A] {
	return Free[F, A]{
		value:  &a,
		isPure: true,
	}
}

// Suspend creates a Free monad from a functor layer
func Suspend[F, A any](f F) Free[F, A] {
	return Free[F, A]{
		suspend: &f,
		isPure:  false,
	}
}

// IsPure returns true if this is a Pure value
func (f Free[F, A]) IsPure() bool {
	return f.isPure
}

// GetPure extracts the Pure value (panics if not Pure)
func (f Free[F, A]) GetPure() A {
	if !f.isPure {
		panic("GetPure called on Suspend")
	}
	return *f.value
}

// GetSuspend extracts the suspended functor (panics if Pure)
func (f Free[F, A]) GetSuspend() F {
	if f.isPure {
		panic("GetSuspend called on Pure")
	}
	return *f.suspend
}

// Map transforms the value (functor)
func (f Free[F, A]) Map(fn func(A) A, fmap func(F) F) Free[F, A] {
	if f.isPure {
		return Pure[F, A](fn(*f.value))
	}
	return Suspend[F, A](fmap(*f.suspend))
}

// FlatMap chains Free computations (monad)
func (f Free[F, A]) FlatMap(
	fn func(A) Free[F, A],
	fmap func(F) F,
) Free[F, A] {
	if f.isPure {
		return fn(*f.value)
	}
	// For suspended values, we'd need to map over F
	// which requires F to be a functor - this is a limitation of concrete types
	return Suspend[F, A](fmap(*f.suspend))
}

// Fold interprets a Free monad by providing:
//   - pure: how to handle Pure values
//   - suspend: how to handle suspended computations
func (f Free[F, A]) Fold(
	pure func(A) any,
	suspend func(F) any,
) any {
	if f.isPure {
		return pure(*f.value)
	}
	return suspend(*f.suspend)
}

// FreeOption is a concrete Free monad for Option
type FreeOption[A any] struct {
	value   *A
	suspend *monoid.Option[FreeOption[A]]
	isPure  bool
}

// PureOption creates a pure FreeOption
func PureOption[A any](a A) FreeOption[A] {
	return FreeOption[A]{
		value:  &a,
		isPure: true,
	}
}

// SuspendOption suspends in an Option layer
func SuspendOption[A any](opt monoid.Option[FreeOption[A]]) FreeOption[A] {
	return FreeOption[A]{
		suspend: &opt,
		isPure:  false,
	}
}

// Interpret runs a FreeOption computation
func (f FreeOption[A]) Interpret() monoid.Option[A] {
	if f.isPure {
		return monoid.Some(*f.value)
	}

	if val, ok := f.suspend.Get(); ok {
		return val.Interpret()
	}
	return monoid.None[A]()
}
