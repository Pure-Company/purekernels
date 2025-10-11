package monoid

// Dual wraps a monoid and reverses the argument order in Combine
// This is a classic categorical construction that transforms any monoid
// into its "opposite" monoid by flipping the operation order
//
// Laws (inherited from wrapped monoid):
//   - Identity: Combine(x, Empty()) == x && Combine(Empty(), x) == x
//   - Associativity: Combine(Combine(x, y), z) == Combine(x, Combine(y, z))
//
// Note: Combine(a, b) in Dual = Combine(b, a) in the wrapped monoid
type Dual[T any] struct {
	inner Monoid[T]
}

// NewDual creates a dual monoid that reverses combination order
func NewDual[T any](m Monoid[T]) Dual[T] {
	return Dual[T]{inner: m}
}

// Empty returns the identity from the wrapped monoid
func (d Dual[T]) Empty() T {
	return d.inner.Empty()
}

// Combine reverses the argument order from the wrapped monoid
// Dual.Combine(a, b) = Inner.Combine(b, a)
func (d Dual[T]) Combine(a, b T) T {
	return d.inner.Combine(b, a)
}

// Inner returns the wrapped monoid
func (d Dual[T]) Inner() Monoid[T] {
	return d.inner
}
