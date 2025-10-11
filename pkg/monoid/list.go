package monoid

// ListMonoid represents list concatenation with empty list as identity
// Laws:
//   - Identity: Combine(x, []) == x && Combine([], x) == x
//   - Associativity: Combine(Combine(x, y), z) == Combine(x, Combine(y, z))
type ListMonoid[T any] struct{}

// NewListMonoid creates a list monoid
func NewListMonoid[T any]() ListMonoid[T] {
	return ListMonoid[T]{}
}

// Empty returns an empty slice
func (ListMonoid[T]) Empty() []T {
	return nil
}

// Combine appends two slices
func (ListMonoid[T]) Combine(a, b []T) []T {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	result := make([]T, len(a)+len(b))
	copy(result, a)
	copy(result[len(a):], b)
	return result
}

// FromItems creates a list from varargs
func FromItems[T any](items ...T) []T {
	return items
}
