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

// Reverse is safe - guaranteed non-empty!
func (nel NonEmptyList[T]) Reverse() NonEmptyList[T] {
	all := nel.ToSlice()
	for i, j := 0, len(all)-1; i < j; i, j = i+1, j-1 {
		all[i], all[j] = all[j], all[i]
	}
	return NonEmptyList[T]{head: all[0], tail: all[1:]}
}

// Last is always safe - no Option needed!
func (nel NonEmptyList[T]) Last() T {
	if len(nel.tail) == 0 {
		return nel.head
	}
	return nel.tail[len(nel.tail)-1]
}

// Zip combines two non-empty lists pairwise
// Standalone function because it introduces new type parameters B and C
func Zip[A, B, C any](
	nel1 NonEmptyList[A],
	nel2 NonEmptyList[B],
	f func(A, B) C,
) NonEmptyList[C] {
	head := f(nel1.head, nel2.head)

	n := min(len(nel1.tail), len(nel2.tail))
	tail := make([]C, n)
	for i := 0; i < n; i++ {
		tail[i] = f(nel1.tail[i], nel2.tail[i])
	}

	return NonEmptyList[C]{head: head, tail: tail}
}

// ZipWith is an alias with more conventional argument order
func ZipWith[A, B, C any](
	f func(A, B) C,
	nel1 NonEmptyList[A],
	nel2 NonEmptyList[B],
) NonEmptyList[C] {
	return Zip(nel1, nel2, f)
}

// NonEmptyListSemigroup - concatenation forms a semigroup
type NonEmptyListSemigroup[T any] struct{}

func NewNonEmptyListSemigroup[T any]() NonEmptyListSemigroup[T] {
	return NonEmptyListSemigroup[T]{}
}

func (NonEmptyListSemigroup[T]) Combine(a, b NonEmptyList[T]) NonEmptyList[T] {
	return a.Concat(b)
}
