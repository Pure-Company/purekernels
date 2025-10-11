package monoid

// MapMergeMonoid represents map merging with empty map as identity
// When keys conflict, the right map's value wins (right-biased)
// Laws:
//   - Identity: Combine(x, {}) == x && Combine({}, x) == x
//   - Associativity: Combine(Combine(x, y), z) == Combine(x, Combine(y, z))
type MapMergeMonoid[K comparable, V any] struct{}

// NewMapMergeMonoid creates a map merge monoid
func NewMapMergeMonoid[K comparable, V any]() MapMergeMonoid[K, V] {
	return MapMergeMonoid[K, V]{}
}

// Empty returns an empty map
func (MapMergeMonoid[K, V]) Empty() map[K]V {
	return make(map[K]V)
}

// Combine merges two maps (right-biased on conflicts)
func (MapMergeMonoid[K, V]) Combine(a, b map[K]V) map[K]V {
	result := make(map[K]V, len(a)+len(b))

	// Copy from a
	for k, v := range a {
		result[k] = v
	}

	// Copy from b (overwrites conflicts)
	for k, v := range b {
		result[k] = v
	}

	return result
}

// MapMergeWithMonoid represents map merging where values are combined using a monoid
// This allows combining values when keys conflict rather than just taking one
type MapMergeWithMonoid[K comparable, V any] struct {
	valueMonoid Monoid[V]
}

// NewMapMergeWithMonoid creates a map merge monoid that combines conflicting values
func NewMapMergeWithMonoid[K comparable, V any](vm Monoid[V]) MapMergeWithMonoid[K, V] {
	return MapMergeWithMonoid[K, V]{valueMonoid: vm}
}

// Empty returns an empty map
func (MapMergeWithMonoid[K, V]) Empty() map[K]V {
	return make(map[K]V)
}

// Combine merges two maps, combining values for duplicate keys using the value monoid
func (m MapMergeWithMonoid[K, V]) Combine(a, b map[K]V) map[K]V {
	result := make(map[K]V, len(a)+len(b))

	// Copy from a
	for k, v := range a {
		result[k] = v
	}

	// Merge from b, combining conflicts
	for k, v := range b {
		if existing, ok := result[k]; ok {
			result[k] = m.valueMonoid.Combine(existing, v)
		} else {
			result[k] = v
		}
	}

	return result
}
