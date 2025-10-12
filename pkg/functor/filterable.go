package functor

import (
	"github.com/vinodhalaharvi/purekernels/pkg/either"
	"github.com/vinodhalaharvi/purekernels/pkg/monoid"
	"github.com/vinodhalaharvi/purekernels/pkg/result"
)

// FilterMap combines filter and map
func FilterMap[A, B any](f func(A) monoid.Option[B], items []A) []B {
	result := make([]B, 0, len(items))
	for _, item := range items {
		if opt := f(item); opt.Valid {
			result = append(result, opt.Value)
		}
	}
	return result
}

// Compact removes None values from []Option[T]
func Compact[T any](opts []monoid.Option[T]) []T {
	result := make([]T, 0, len(opts))
	for _, opt := range opts {
		if opt.Valid {
			result = append(result, opt.Value)
		}
	}
	return result
}

// CompactResults removes Err values from []Result[T]
func CompactResults[T any](results []result.Result[T]) []T {
	values := make([]T, 0, len(results))
	for _, r := range results {
		if val, err := r.Get(); err == nil {
			values = append(values, val)
		}
	}
	return values
}

// Separate splits []Either[L,R] into ([]L, []R)
func Separate[L, R any](eithers []either.Either[L, R]) ([]L, []R) {
	return either.Partition(eithers)
}

// FilterMapResult filters and transforms with Result
func FilterMapResult[A, B any](f func(A) result.Result[B], items []A) []B {
	result := make([]B, 0, len(items))
	for _, item := range items {
		if val, err := f(item).Get(); err == nil {
			result = append(result, val)
		}
	}
	return result
}
