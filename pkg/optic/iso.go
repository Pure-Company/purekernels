package optic

import "github.com/vinodhalaharvi/purekernels/pkg/monoid"

// Iso represents an isomorphism (bidirectional conversion)
// Laws:
//   - FromTo: From(To(s)) == s
//   - ToFrom: To(From(a)) == a
type Iso[S, A any] struct {
	To   func(S) A
	From func(A) S
}

// NewIso creates an isomorphism
func NewIso[S, A any](to func(S) A, from func(A) S) Iso[S, A] {
	return Iso[S, A]{To: to, From: from}
}

// Reverse reverses the isomorphism
func (i Iso[S, A]) Reverse() Iso[A, S] {
	return Iso[A, S]{To: i.From, From: i.To}
}

// ComposeIso composes two isomorphisms
func ComposeIso[S, A, B any](outer Iso[S, A], inner Iso[A, B]) Iso[S, B] {
	return Iso[S, B]{
		To:   func(s S) B { return inner.To(outer.To(s)) },
		From: func(b B) S { return outer.From(inner.From(b)) },
	}
}

// IsoToLens converts an isomorphism to a lens
func IsoToLens[S, A any](i Iso[S, A]) Lens[S, A] {
	return Lens[S, A]{
		Get: i.To,
		Set: func(_ S, a A) S { return i.From(a) },
	}
}

// IsoToPrism converts an isomorphism to a prism
func IsoToPrism[S, A any](i Iso[S, A]) Prism[S, A] {
	return Prism[S, A]{
		Preview: func(s S) monoid.Option[A] {
			return monoid.Some(i.To(s))
		},
		Review: i.From,
	}
}
