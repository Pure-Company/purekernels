// pkg/comonad/traced.go
package comonad

import "github.com/vinodhalaharvi/purekernels/pkg/monoid"

// Traced is the dual of Writer - a function from monoid to value
// This represents "position-dependent computation"
// Laws: Same comonad laws
type Traced[M, A any] struct {
	trace  func(M) A
	monoid monoid.Monoid[M]
}

// NewTraced creates a Traced computation
func NewTraced[M, A any](m monoid.Monoid[M], trace func(M) A) Traced[M, A] {
	return Traced[M, A]{
		trace:  trace,
		monoid: m,
	}
}

// Extract evaluates at the identity (comonad extract)
func (t Traced[M, A]) Extract() A {
	return t.trace(t.monoid.Empty())
}

// Extend applies a function that can probe different positions
func (t Traced[M, A]) Extend(f func(Traced[M, A]) A) Traced[M, A] {
	return Traced[M, A]{
		trace: func(m M) A {
			shifted := Traced[M, A]{
				trace: func(m2 M) A {
					return t.trace(t.monoid.Combine(m, m2))
				},
				monoid: t.monoid,
			}
			return f(shifted)
		},
		monoid: t.monoid,
	}
}

// Map transforms the value
func (t Traced[M, A]) Map(f func(A) A) Traced[M, A] {
	return Traced[M, A]{
		trace:  func(m M) A { return f(t.trace(m)) },
		monoid: t.monoid,
	}
}

// Trace evaluates at a specific position
func (t Traced[M, A]) Trace(m M) A {
	return t.trace(m)
}

// Traces applies multiple positions
func (t Traced[M, A]) Traces(positions []M) []A {
	results := make([]A, len(positions))
	for i, pos := range positions {
		results[i] = t.trace(pos)
	}
	return results
}

// Listen returns both the value and the traced position
func (t Traced[M, A]) Listen(m M) (A, M) {
	return t.trace(m), m
}
