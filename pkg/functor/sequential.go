package functor

import (
	"github.com/vinodhalaharvi/purekernels/pkg/fold"
)

// CombineAll combines multiple applicatives
// Works with both sequential and concurrent implementations!
func CombineAll[T any](first Applicative[T], rest ...Applicative[T]) Applicative[T] {
	result := first
	for _, app := range rest {
		result = result.Apply(app)
	}
	return result
}

// Sequence extracts values from multiple applicatives
func Sequence[T any](apps []Applicative[T]) []T {
	result := make([]T, len(apps))
	for i, app := range apps {
		result[i] = app.Value()
	}
	return result
}

// Traverse maps function over items and combines results
// The magic: if you pass ConcurrentApplicative, this runs in parallel!
func Traverse[A any, T any](
	f func(A) Applicative[T],
	items []A,
) Applicative[T] {
	if len(items) == 0 {
		panic("cannot traverse empty list without monoid")
	}

	apps := fold.Map(f, items)
	return fold.FoldLeft(
		func(acc Applicative[T], app Applicative[T]) Applicative[T] {
			return acc.Apply(app)
		},
		apps[0],
		apps[1:],
	)
}

// TraverseWithEmpty handles empty lists
func TraverseWithEmpty[A any, T any](
	empty Applicative[T],
	f func(A) Applicative[T],
	items []A,
) Applicative[T] {
	if len(items) == 0 {
		return empty
	}
	return Traverse(f, items)
}

// Fold accumulates using applicative
func Fold[A any, T any](
	empty Applicative[T],
	f func(A) Applicative[T],
	items []A,
) Applicative[T] {
	return fold.FoldLeft(
		func(acc Applicative[T], item A) Applicative[T] {
			return acc.Apply(f(item))
		},
		empty,
		items,
	)
}

// ApplyIf conditionally applies
func ApplyIf[T any](
	condition bool,
	base Applicative[T],
	conditional Applicative[T],
) Applicative[T] {
	if condition {
		return base.Apply(conditional)
	}
	return base
}

// ApplyWhen with predicate function
func ApplyWhen[T any](
	predicate func() bool,
	base Applicative[T],
	conditional Applicative[T],
) Applicative[T] {
	return ApplyIf(predicate(), base, conditional)
}
