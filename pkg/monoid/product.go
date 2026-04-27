package monoid

// ProductMonoid represents multiplication with 1 as identity
// Laws:
//   - Identity: Combine(x, 1) == x && Combine(1, x) == x
//   - Associativity: Combine(Combine(x, y), z) == Combine(x, Combine(y, z))
type ProductMonoid[T Numeric] struct{}

// NewProductMonoid creates a product monoid
func NewProductMonoid[T Numeric]() ProductMonoid[T] {
	return ProductMonoid[T]{}
}

// Empty returns 1 (multiplicative identity)
func (ProductMonoid[T]) Empty() T {
	return 1
}

// Combine performs multiplication
func (ProductMonoid[T]) Combine(a, b T) T {
	return a * b
}
