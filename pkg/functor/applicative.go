// Package functor provides categorical abstractions for composable computations
package functor

// Applicative represents a computation that accumulates values
// This is the foundation for dependency tracking in AST extraction
//
// Laws (for mathematical correctness):
//   - Identity: Pure(empty) <*> v = v
//   - Associativity: (u <*> v) <*> w = u <*> (v <*> w)
//   - Homomorphism: Pure(x) <*> Pure(y) = Pure(x ⊕ y)
type Applicative[T any] interface {
	// Apply combines this computation with another using monoid operation
	Apply(other Applicative[T]) Applicative[T]

	// Value extracts the accumulated result
	Value() T

	// Map transforms the accumulated value
	Map(f func(T) T) Applicative[T]
}
