package optic

import "github.com/vinodhalaharvi/purekernels/pkg/monoid"

// Prism focuses on a part that may not exist (for sum types)
// Laws:
//   - ReviewPreview: Preview(Review(a)) == Some(a)
//   - PreviewReview: match Preview(s): Some(a) => Review(a) == s
type Prism[S, A any] struct {
	Preview func(S) monoid.Option[A]
	Review  func(A) S
}

// NewPrism creates a prism from preview and review
func NewPrism[S, A any](
	preview func(S) monoid.Option[A],
	review func(A) S,
) Prism[S, A] {
	return Prism[S, A]{Preview: preview, Review: review}
}

// Modify applies a function if the focus exists
func (p Prism[S, A]) Modify(f func(A) A) func(S) S {
	return func(s S) S {
		if a, ok := p.Preview(s).Get(); ok {
			return p.Review(f(a))
		}
		return s
	}
}

// ComposePrism composes two prisms
func ComposePrism[S, A, B any](outer Prism[S, A], inner Prism[A, B]) Prism[S, B] {
	return Prism[S, B]{
		Preview: func(s S) monoid.Option[B] {
			if a, ok := outer.Preview(s).Get(); ok {
				return inner.Preview(a)
			}
			return monoid.None[B]()
		},
		Review: func(b B) S {
			return outer.Review(inner.Review(b))
		},
	}
}
