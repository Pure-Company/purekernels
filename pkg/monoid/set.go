package monoid

// SetMonoid represents a set with union as the operation
type SetMonoid[T comparable] struct {
	items map[T]struct{}
}

// NewSetMonoid creates a new empty set monoid
func NewSetMonoid[T comparable]() SetMonoid[T] {
	return SetMonoid[T]{
		items: make(map[T]struct{}),
	}
}

// Empty returns an empty set
func (s SetMonoid[T]) Empty() SetMonoid[T] {
	return NewSetMonoid[T]()
}

// Combine performs set union
func (s SetMonoid[T]) Combine(a, b SetMonoid[T]) SetMonoid[T] {
	result := NewSetMonoid[T]()

	// Add all from a
	for k := range a.items {
		result.items[k] = struct{}{}
	}

	// Add all from b
	for k := range b.items {
		result.items[k] = struct{}{}
	}

	return result
}

// Insert adds an element to the set
func (s SetMonoid[T]) Insert(item T) SetMonoid[T] {
	result := NewSetMonoid[T]()

	// Copy existing items
	for k := range s.items {
		result.items[k] = struct{}{}
	}

	// Add new item
	result.items[item] = struct{}{}
	return result
}

// Contains checks if element exists
func (s SetMonoid[T]) Contains(item T) bool {
	_, ok := s.items[item]
	return ok
}

// ToSlice converts set to slice
func (s SetMonoid[T]) ToSlice() []T {
	result := make([]T, 0, len(s.items))
	for k := range s.items {
		result = append(result, k)
	}
	return result
}

// Size returns the number of elements
func (s SetMonoid[T]) Size() int {
	return len(s.items)
}

// FromSlice creates a set from a slice
func FromSlice[T comparable](items []T) SetMonoid[T] {
	result := NewSetMonoid[T]()
	for _, item := range items {
		result.items[item] = struct{}{}
	}
	return result
}
