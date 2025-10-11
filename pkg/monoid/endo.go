package monoid

// EndoMonoid represents function composition (endomorphisms: T -> T)
// Identity is the identity function, and composition is function composition
// Laws:
//   - Identity: Combine(f, id) == f && Combine(id, f) == f
//   - Associativity: Combine(Combine(f, g), h) == Combine(f, Combine(g, h))
type EndoMonoid[T any] struct{}

// NewEndoMonoid creates an endomorphism monoid
func NewEndoMonoid[T any]() EndoMonoid[T] {
	return EndoMonoid[T]{}
}

// Empty returns the identity function
func (EndoMonoid[T]) Empty() func(T) T {
	return func(x T) T { return x }
}

// Combine composes two functions: (f ∘ g)(x) = f(g(x))
// Note: Applies g first, then f (standard mathematical composition)
func (EndoMonoid[T]) Combine(f, g func(T) T) func(T) T {
	return func(x T) T {
		return f(g(x))
	}
}
