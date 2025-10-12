// pkg/comonad/store.go
// Package comonad provides comonadic structures - duals of monads
package comonad

// Store is the dual of State - a position in a space with a lookup function
// This models "context with a focus point"
// Laws:
//   - Extract-Extend: Extract(Extend(w, f)) ≡ f(w)
//   - Extend-Extract: Extend(w, Extract) ≡ w
//   - Extend-Extend: Extend(Extend(w, f), g) ≡ Extend(w, func(w2) { return g(Extend(w2, f)) })
type Store[S, A any] struct {
	lookup   func(S) A
	position S
}

// NewStore creates a Store with a lookup function and initial position
func NewStore[S, A any](lookup func(S) A, position S) Store[S, A] {
	return Store[S, A]{
		lookup:   lookup,
		position: position,
	}
}

// Extract gets the value at the current position (comonad extract)
func (s Store[S, A]) Extract() A {
	return s.lookup(s.position)
}

// Extend applies a function that sees the whole Store (comonad extend)
func (s Store[S, A]) Extend(f func(Store[S, A]) A) Store[S, A] {
	return Store[S, A]{
		lookup: func(pos S) A {
			return f(Store[S, A]{lookup: s.lookup, position: pos})
		},
		position: s.position,
	}
}

// Map transforms the values (functor)
func (s Store[S, A]) Map(f func(A) A) Store[S, A] {
	return Store[S, A]{
		lookup:   func(pos S) A { return f(s.lookup(pos)) },
		position: s.position,
	}
}

// Seek moves to a new position
func (s Store[S, A]) Seek(newPos S) Store[S, A] {
	return Store[S, A]{
		lookup:   s.lookup,
		position: newPos,
	}
}

// Pos returns the current position
func (s Store[S, A]) Pos() S {
	return s.position
}

// Peek looks at a value without moving
func (s Store[S, A]) Peek(pos S) A {
	return s.lookup(pos)
}

// Experiment applies a function to position and extracts values
func (s Store[S, A]) Experiment(f func(S) []S) []A {
	positions := f(s.position)
	results := make([]A, len(positions))
	for i, pos := range positions {
		results[i] = s.lookup(pos)
	}
	return results
}
