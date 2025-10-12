// pkg/profunctor/profunctor.go
// Package profunctor provides profunctor structures for optics composition
package profunctor

import (
	"github.com/vinodhalaharvi/purekernels/pkg/monoid"
	"github.com/vinodhalaharvi/purekernels/pkg/optic"
)

// Star represents a profunctor A -> F[B] (Kleisli-like)
type Star[F, A, B any] struct {
	run func(A) F
}

// NewStar creates a Star profunctor
func NewStar[F, A, B any](f func(A) F) Star[F, A, B] {
	return Star[F, A, B]{run: f}
}

// Dimap is the profunctor operation: contramap on input, map on output
func (s Star[F, A, B]) Dimap(
	pre func(A) A,
	post func(F) F,
) Star[F, A, B] {
	return Star[F, A, B]{
		run: func(a A) F {
			return post(s.run(pre(a)))
		},
	}
}

// Costar represents a profunctor F[A] -> B (CoKleisli-like)
type Costar[F, A, B any] struct {
	run func(F) B
}

// NewCostar creates a Costar profunctor
func NewCostar[F, A, B any](f func(F) B) Costar[F, A, B] {
	return Costar[F, A, B]{run: f}
}

// Dimap is the profunctor operation
func (c Costar[F, A, B]) Dimap(
	pre func(F) F,
	post func(B) B,
) Costar[F, A, B] {
	return Costar[F, A, B]{
		run: func(f F) B {
			return post(c.run(pre(f)))
		},
	}
}

// Strong profunctor - for lenses
type Strong[P, A, B, C any] struct {
	profunctor P
	first      func(P) any // First[P, A, B, C] - pairs first component
}

// Choice profunctor - for prisms
type Choice[P, A, B, C any] struct {
	profunctor P
	left       func(P) any // Left[P, A, B, C] - eithers left component
}

// Wander profunctor - for traversals
type Wander[P, A, B any] struct {
	profunctor P
	traverse   func(P) any // Traverse[P, A, B]
}

// LensP is a lens via profunctors
type LensP[S, A any] struct {
	view   func(S) A
	update func(S, A) S
}

// NewLensP creates a profunctor-based lens
func NewLensP[S, A any](
	view func(S) A,
	update func(S, A) S,
) LensP[S, A] {
	return LensP[S, A]{
		view:   view,
		update: update,
	}
}

// ToLens converts profunctor lens to concrete lens
func (l LensP[S, A]) ToLens() optic.Lens[S, A] {
	return optic.NewLens(l.view, l.update)
}

// PrismP is a prism via profunctors
type PrismP[S, A any] struct {
	match func(S) monoid.Option[A]
	build func(A) S
}

// NewPrismP creates a profunctor-based prism
func NewPrismP[S, A any](
	match func(S) monoid.Option[A],
	build func(A) S,
) PrismP[S, A] {
	return PrismP[S, A]{
		match: match,
		build: build,
	}
}

// ToPrism converts profunctor prism to concrete prism
func (p PrismP[S, A]) ToPrism() optic.Prism[S, A] {
	return optic.NewPrism(p.match, p.build)
}

// ComposeLensP composes two profunctor lenses
func ComposeLensP[S, A, B any](
	outer LensP[S, A],
	inner LensP[A, B],
) LensP[S, B] {
	return LensP[S, B]{
		view: func(s S) B {
			return inner.view(outer.view(s))
		},
		update: func(s S, b B) S {
			a := outer.view(s)
			newA := inner.update(a, b)
			return outer.update(s, newA)
		},
	}
}

// ComposePrismP composes two profunctor prisms
func ComposePrismP[S, A, B any](
	outer PrismP[S, A],
	inner PrismP[A, B],
) PrismP[S, B] {
	return PrismP[S, B]{
		match: func(s S) monoid.Option[B] {
			if a, ok := outer.match(s).Get(); ok {
				return inner.match(a)
			}
			return monoid.None[B]()
		},
		build: func(b B) S {
			return outer.build(inner.build(b))
		},
	}
}
