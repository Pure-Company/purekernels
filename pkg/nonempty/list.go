package nonempty

import "github.com/vinodhalaharvi/purekernels/pkg/monoid"

// NonEmptyList guarantees at least one element
type NonEmptyList[T any] struct {
	head T
	tail []T
}

// New creates a non-empty list
func New[T any](head T, tail ...T) NonEmptyList[T] {
	return NonEmptyList[T]{head: head, tail: tail}
}

// FromSlice converts a slice to NonEmptyList (returns Option)
func FromSlice[T any](items []T) monoid.Option[NonEmptyList[T]] {
	if len(items) == 0 {
		return monoid.None[NonEmptyList[T]]()
	}
	return monoid.Some(NonEmptyList[T]{
		head: items[0],
		tail: items[1:],
	})
}

// Head returns the first element (always safe!)
func (nel NonEmptyList[T]) Head() T {
	return nel.head
}

// Tail returns remaining elements
func (nel NonEmptyList[T]) Tail() []T {
	return nel.tail
}

// ToSlice converts to regular slice
func (nel NonEmptyList[T]) ToSlice() []T {
	result := make([]T, 0, 1+len(nel.tail))
	result = append(result, nel.head)
	result = append(result, nel.tail...)
	return result
}

// Size returns the number of elements (always >= 1)
func (nel NonEmptyList[T]) Size() int {
	return 1 + len(nel.tail)
}

// Map transforms all elements
func (nel NonEmptyList[T]) Map(f func(T) T) NonEmptyList[T] {
	newTail := make([]T, len(nel.tail))
	for i, v := range nel.tail {
		newTail[i] = f(v)
	}
	return NonEmptyList[T]{
		head: f(nel.head),
		tail: newTail,
	}
}

// FlatMap transforms and flattens
func (nel NonEmptyList[T]) FlatMap(f func(T) NonEmptyList[T]) NonEmptyList[T] {
	first := f(nel.head)
	result := first.ToSlice()

	for _, item := range nel.tail {
		result = append(result, f(item).ToSlice()...)
	}

	return NonEmptyList[T]{
		head: result[0],
		tail: result[1:],
	}
}

// Append adds elements (always non-empty!)
func (nel NonEmptyList[T]) Append(items ...T) NonEmptyList[T] {
	return NonEmptyList[T]{
		head: nel.head,
		tail: append(nel.tail, items...),
	}
}

// Concat combines two non-empty lists
func (nel NonEmptyList[T]) Concat(other NonEmptyList[T]) NonEmptyList[T] {
	return nel.Append(other.ToSlice()...)
}

// Reduce is always safe - no need for empty case!
func (nel NonEmptyList[T]) Reduce(f func(T, T) T) T {
	result := nel.head
	for _, item := range nel.tail {
		result = f(result, item)
	}
	return result
}

// NonEmptyListSemigroup - concatenation forms a semigroup
type NonEmptyListSemigroup[T any] struct{}

func NewNonEmptyListSemigroup[T any]() NonEmptyListSemigroup[T] {
	return NonEmptyListSemigroup[T]{}
}

func (NonEmptyListSemigroup[T]) Combine(a, b NonEmptyList[T]) NonEmptyList[T] {
	return a.Concat(b)
}
