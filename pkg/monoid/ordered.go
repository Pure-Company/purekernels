package monoid

import "cmp"

// MaxMonoid represents the maximum operation with the minimum value as identity
// Laws:
//   - Identity: Combine(x, minValue) == x && Combine(minValue, x) == x
//   - Associativity: Combine(Combine(x, y), z) == Combine(x, Combine(y, z))
type MaxMonoid[T cmp.Ordered] struct {
	minValue T
}

// NewMaxMonoid creates a max monoid with specified minimum value
func NewMaxMonoid[T cmp.Ordered](minValue T) MaxMonoid[T] {
	return MaxMonoid[T]{minValue: minValue}
}

// Empty returns the minimum value (identity for max)
func (m MaxMonoid[T]) Empty() T {
	return m.minValue
}

// Combine returns the maximum of two values
func (MaxMonoid[T]) Combine(a, b T) T {
	if a > b {
		return a
	}
	return b
}

// MinMonoid represents the minimum operation with the maximum value as identity
// Laws:
//   - Identity: Combine(x, maxValue) == x && Combine(maxValue, x) == x
//   - Associativity: Combine(Combine(x, y), z) == Combine(x, Combine(y, z))
type MinMonoid[T cmp.Ordered] struct {
	maxValue T
}

// NewMinMonoid creates a min monoid with specified maximum value
func NewMinMonoid[T cmp.Ordered](maxValue T) MinMonoid[T] {
	return MinMonoid[T]{maxValue: maxValue}
}

// Empty returns the maximum value (identity for min)
func (m MinMonoid[T]) Empty() T {
	return m.maxValue
}

// Combine returns the minimum of two values
func (MinMonoid[T]) Combine(a, b T) T {
	if a < b {
		return a
	}
	return b
}
