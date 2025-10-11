package functor

import (
	"github.com/vinodhalaharvi/purekernels/pkg/monoid"
)

// Sequential is a sequential computation that accumulates values using a monoid
// Unlike Concurrent, this executes computations one after another
//
// Laws (same as Concurrent):
//   - Identity: Pure(empty).Apply(v) ≡ v
//   - Composition: u.Apply(v).Apply(w) ≡ u.Apply(v.Apply(w))
//   - Homomorphism: Pure(x).Apply(Pure(y)) ≡ Pure(combine(x,y))
type Sequential[T any] struct {
	compute func() T
	monoid  monoid.Monoid[T]
}

// NewSequential creates a sequential computation
func NewSequential[T any](m monoid.Monoid[T], compute func() T) Sequential[T] {
	return Sequential[T]{
		compute: compute,
		monoid:  m,
	}
}

// PureSequential wraps a value in a sequential computation
func PureSequential[T any](m monoid.Monoid[T], value T) Sequential[T] {
	return Sequential[T]{
		compute: func() T { return value },
		monoid:  m,
	}
}

// Apply combines two sequential computations (executes left then right)
func (s Sequential[T]) Apply(other Sequential[T]) Sequential[T] {
	return Sequential[T]{
		compute: func() T {
			left := s.compute()
			right := other.compute()
			return s.monoid.Combine(left, right)
		},
		monoid: s.monoid,
	}
}

// Map transforms the result of a computation
func (s Sequential[T]) Map(f func(T) T) Sequential[T] {
	return Sequential[T]{
		compute: func() T {
			return f(s.compute())
		},
		monoid: s.monoid,
	}
}

// Value executes the computation and returns the result
func (s Sequential[T]) Value() T {
	return s.compute()
}

// CombineAllSequential combines multiple sequential computations
func CombineAllSequential[T any](first Sequential[T], rest ...Sequential[T]) Sequential[T] {
	result := first
	for _, comp := range rest {
		result = result.Apply(comp)
	}
	return result
}

// SequenceSequential executes multiple computations sequentially and extracts values
func SequenceSequential[T any](comps []Sequential[T]) []T {
	result := make([]T, len(comps))
	for i, comp := range comps {
		result[i] = comp.Value()
	}
	return result
}

// TraverseSequential maps a function over items and combines results sequentially
func TraverseSequential[A any, T any](
	m monoid.Monoid[T],
	f func(A) Sequential[T],
	items []A,
) Sequential[T] {
	if len(items) == 0 {
		return PureSequential(m, m.Empty())
	}

	result := f(items[0])
	for i := 1; i < len(items); i++ {
		result = result.Apply(f(items[i]))
	}
	return result
}

// TraverseSequentialDirect processes items sequentially with direct function
func TraverseSequentialDirect[A any, T any](
	m monoid.Monoid[T],
	f func(A) T,
	items []A,
) Sequential[T] {
	return Sequential[T]{
		compute: func() T {
			result := m.Empty()
			for _, item := range items {
				result = m.Combine(result, f(item))
			}
			return result
		},
		monoid: m,
	}
}

// ApplyIfSequential conditionally applies a computation
func ApplyIfSequential[T any](
	condition bool,
	base Sequential[T],
	conditional Sequential[T],
) Sequential[T] {
	if condition {
		return base.Apply(conditional)
	}
	return base
}

// ApplyWhenSequential applies a computation if predicate is true
func ApplyWhenSequential[T any](
	predicate func() bool,
	base Sequential[T],
	conditional Sequential[T],
) Sequential[T] {
	return ApplyIfSequential(predicate(), base, conditional)
}

// Fold accumulates using sequential applicative
func FoldSequential[A any, T any](
	m monoid.Monoid[T],
	f func(A) T,
	items []A,
) Sequential[T] {
	return Sequential[T]{
		compute: func() T {
			result := m.Empty()
			for _, item := range items {
				result = m.Combine(result, f(item))
			}
			return result
		},
		monoid: m,
	}
}
