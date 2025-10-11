package optic

// Traversal focuses on multiple targets
type Traversal[S, A any] struct {
	ModifyF func(func(A) A) func(S) S
}

// NewTraversal creates a traversal
func NewTraversal[S, A any](modifyF func(func(A) A) func(S) S) Traversal[S, A] {
	return Traversal[S, A]{ModifyF: modifyF}
}

// Modify applies a function to all focuses
func (t Traversal[S, A]) Modify(f func(A) A) func(S) S {
	return t.ModifyF(f)
}

// Each creates a traversal for slice elements
func Each[A any]() Traversal[[]A, A] {
	return Traversal[[]A, A]{
		ModifyF: func(f func(A) A) func([]A) []A {
			return func(xs []A) []A {
				result := make([]A, len(xs))
				for i, x := range xs {
					result[i] = f(x)
				}
				return result
			}
		},
	}
}

// ComposeTraversal composes traversals
func ComposeTraversal[S, A, B any](outer Traversal[S, A], inner Traversal[A, B]) Traversal[S, B] {
	return Traversal[S, B]{
		ModifyF: func(f func(B) B) func(S) S {
			return outer.Modify(inner.Modify(f))
		},
	}
}
