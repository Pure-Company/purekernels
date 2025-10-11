// Package optic provides composable data accessors and transformations
package optic

// Lens focuses on a part of a structure (getter + setter)
// Laws:
//   - GetPut: Set(s, Get(s)) == s
//   - PutGet: Get(Set(s, a)) == a
//   - PutPut: Set(Set(s, a), b) == Set(s, b)
type Lens[S, A any] struct {
	Get func(S) A
	Set func(S, A) S
}

// NewLens creates a lens from getter and setter
func NewLens[S, A any](get func(S) A, set func(S, A) S) Lens[S, A] {
	return Lens[S, A]{Get: get, Set: set}
}

// Modify applies a function through the lens
func (l Lens[S, A]) Modify(f func(A) A) func(S) S {
	return func(s S) S {
		return l.Set(s, f(l.Get(s)))
	}
}

// ComposeLens composes two lenses
func ComposeLens[S, A, B any](outer Lens[S, A], inner Lens[A, B]) Lens[S, B] {
	return Lens[S, B]{
		Get: func(s S) B {
			return inner.Get(outer.Get(s))
		},
		Set: func(s S, b B) S {
			a := outer.Get(s)
			newA := inner.Set(a, b)
			return outer.Set(s, newA)
		},
	}
}

// Identity lens (focuses on the whole structure)
func Identity[S any]() Lens[S, S] {
	return Lens[S, S]{
		Get: func(s S) S { return s },
		Set: func(_ S, s S) S { return s },
	}
}
