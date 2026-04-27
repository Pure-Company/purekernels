package functor

import (
	"github.com/Pure-Company/purekernels/pkg/monoid"
	"github.com/Pure-Company/purekernels/pkg/result"
)

// MapSlice transforms each element in a slice
// This is the functor law for slices: fmap f xs
func MapSlice[A, B any](f func(A) B, xs []A) []B {
	result := make([]B, len(xs))
	for i, x := range xs {
		result[i] = f(x)
	}
	return result
}

// MapOption transforms the value inside an Option if present
func MapOption[A, B any](f func(A) B, opt monoid.Option[A]) monoid.Option[B] {
	if val, ok := opt.Get(); ok {
		return monoid.Some(f(val))
	}
	return monoid.None[B]()
}

// MapResult transforms the value inside a Result if Ok
func MapResult[A, B any](f func(A) B, r result.Result[A]) result.Result[B] {
	if val, err := r.Get(); err == nil {
		return result.Ok(f(val))
	} else {
		return result.Err[B](err)
	}
}

// MapMap transforms values in a map
func MapMap[K comparable, A, B any](f func(A) B, m map[K]A) map[K]B {
	result := make(map[K]B, len(m))
	for k, v := range m {
		result[k] = f(v)
	}
	return result
}

// FlatMapSlice transforms and flattens slices (monadic bind)
func FlatMapSlice[A, B any](f func(A) []B, xs []A) []B {
	result := []B{}
	for _, x := range xs {
		result = append(result, f(x)...)
	}
	return result
}

// FlatMapOption chains optional computations
func FlatMapOption[A, B any](f func(A) monoid.Option[B], opt monoid.Option[A]) monoid.Option[B] {
	if val, ok := opt.Get(); ok {
		return f(val)
	}
	return monoid.None[B]()
}

// Sequence converts []Result[T] to Result[[]T]
func SequenceResults[T any](results []result.Result[T]) result.Result[[]T] {
	return result.Collect(results)
}

// Sequence converts []Option[T] to Option[[]T]
func SequenceOptions[T any](opts []monoid.Option[T]) monoid.Option[[]T] {
	values := make([]T, 0, len(opts))
	for _, opt := range opts {
		if val, ok := opt.Get(); ok {
			values = append(values, val)
		} else {
			return monoid.None[[]T]()
		}
	}
	return monoid.Some(values)
}
