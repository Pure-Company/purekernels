package monoid

import "golang.org/x/exp/constraints"

// SumMonoid represents addition with 0 as identity
// Laws:
//   - Identity: Combine(x, 0) == x && Combine(0, x) == x
//   - Associativity: Combine(Combine(x, y), z) == Combine(x, Combine(y, z))
type SumMonoid[T constraints.Integer | constraints.Float] struct{}

// NewSumMonoid creates a sum monoid
func NewSumMonoid[T constraints.Integer | constraints.Float]() SumMonoid[T] {
	return SumMonoid[T]{}
}

// Empty returns 0 (additive identity)
func (SumMonoid[T]) Empty() T {
	return 0
}

// Combine performs addition
func (SumMonoid[T]) Combine(a, b T) T {
	return a + b
}
