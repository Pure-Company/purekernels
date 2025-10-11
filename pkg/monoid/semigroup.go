package monoid

// Semigroup represents a set with an associative binary operation
// Unlike Monoid, it doesn't require an identity element
// Laws:
//   - Associativity: Combine(Combine(x, y), z) == Combine(x, Combine(y, z))
type Semigroup[T any] interface {
	// Combine performs the associative binary operation
	Combine(a, b T) T
}

// Every Monoid is a Semigroup
// ToSemigroup converts a monoid to a semigroup (identity is simply not used)
func ToSemigroup[T any](m Monoid[T]) Semigroup[T] {
	return semigroupAdapter[T]{m}
}

type semigroupAdapter[T any] struct {
	monoid Monoid[T]
}

func (s semigroupAdapter[T]) Combine(a, b T) T {
	return s.monoid.Combine(a, b)
}

// ReduceSemigroup folds a non-empty slice using a semigroup
// Panics if the slice is empty (no identity element to return)
func ReduceSemigroup[T any](s Semigroup[T], items []T) T {
	if len(items) == 0 {
		panic("ReduceSemigroup: cannot reduce empty slice (no identity)")
	}

	result := items[0]
	for i := 1; i < len(items); i++ {
		result = s.Combine(result, items[i])
	}
	return result
}

// ReduceSemigroupOpt safely reduces with Option result
func ReduceSemigroupOpt[T any](s Semigroup[T], items []T) Option[T] {
	if len(items) == 0 {
		return None[T]()
	}
	return Some(ReduceSemigroup(s, items))
}
