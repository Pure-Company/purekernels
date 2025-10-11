// Package monoid provides core categorical abstractions for combining values
package monoid

// Monoid represents a set with an associative binary operation and identity element
// Laws:
//   - Identity: Combine(x, Empty()) == x && Combine(Empty(), x) == x
//   - Associativity: Combine(Combine(x, y), z) == Combine(x, Combine(y, z))
type Monoid[T any] interface {
	// Empty returns the identity element
	Empty() T

	// Combine performs the associative binary operation
	Combine(a, b T) T
}

// Reduce folds a slice using a monoid - replaces imperative loops
func Reduce[T any](m Monoid[T], items []T) T {
	result := m.Empty()
	for _, item := range items {
		result = m.Combine(result, item)
	}
	return result
}

// Map transforms and combines - replaces map+reduce loops
func Map[A, B any](m Monoid[B], f func(A) B, items []A) B {
	result := m.Empty()
	for _, item := range items {
		result = m.Combine(result, f(item))
	}
	return result
}
