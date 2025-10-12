// Package result provides sum types for error handling without exceptions
package result

import "github.com/vinodhalaharvi/purekernels/pkg/monoid"

// Result represents a computation that may succeed (Ok) or fail (Err)
// This is a sum type: either Value is valid, or Err is present
type Result[T any] struct {
	value T
	err   error
	isOk  bool
}

// Ok creates a successful result
func Ok[T any](value T) Result[T] {
	return Result[T]{value: value, isOk: true}
}

// Err creates a failed result
func Err[T any](err error) Result[T] {
	return Result[T]{err: err, isOk: false}
}

// IsOk returns true if result is successful
func (r Result[T]) IsOk() bool {
	return r.isOk
}

// IsErr returns true if result is an error
func (r Result[T]) IsErr() bool {
	return !r.isOk
}

// Unwrap returns the value or panics if error
func (r Result[T]) Unwrap() T {
	if !r.isOk {
		panic("Unwrap called on Err result")
	}
	return r.value
}

// UnwrapOr returns the value or a default
func (r Result[T]) UnwrapOr(defaultValue T) T {
	if r.isOk {
		return r.value
	}
	return defaultValue
}

// Get returns value and error (Go-style)
func (r Result[T]) Get() (T, error) {
	return r.value, r.err
}

// Map transforms the value if Ok, propagates error otherwise
func (r Result[T]) Map(f func(T) T) Result[T] {
	if !r.isOk {
		return Result[T]{err: r.err, isOk: false}
	}
	return Ok(f(r.value))
}

// MapErr transforms the error if Err, propagates value otherwise
func (r Result[T]) MapErr(f func(error) error) Result[T] {
	if r.isOk {
		return r
	}
	return Err[T](f(r.err))
}

// Recover attempts to recover from error
func (r Result[T]) Recover(f func(error) T) Result[T] {
	if r.isOk {
		return r
	}
	return Ok(f(r.err))
}

// ToOption converts Result to Option, discarding error
func (r Result[T]) ToOption() monoid.Option[T] {
	if r.isOk {
		return monoid.Some(r.value)
	}
	return monoid.None[T]()
}

// FromOption converts Option to Result with custom error
func FromOption[T any](opt monoid.Option[T], err error) Result[T] {
	if value, ok := opt.Get(); ok {
		return Ok(value)
	}
	return Err[T](err)
}

// Collect transforms a slice of Results into a Result of slice
// Returns first error encountered, or Ok with all values
func Collect[T any](results []Result[T]) Result[[]T] {
	values := make([]T, 0, len(results))
	for _, r := range results {
		if !r.isOk {
			return Err[[]T](r.err)
		}
		values = append(values, r.value)
	}
	return Ok(values)
}

// Partition separates Ok and Err results
func Partition[T any](results []Result[T]) ([]T, []error) {
	oks := []T{}
	errs := []error{}
	for _, r := range results {
		if r.isOk {
			oks = append(oks, r.value)
		} else {
			errs = append(errs, r.err)
		}
	}
	return oks, errs
}

// FlatMapResult allows type transformation A -> B
func FlatMapResult[A, B any](
	r Result[A],
	f func(A) Result[B],
) Result[B] {
	if !r.isOk {
		return Err[B](r.err)
	}
	return f(r.value)
}

// Keep existing method for same-type chaining
func (r Result[T]) FlatMap(f func(T) Result[T]) Result[T] {
	if !r.isOk {
		return Result[T]{err: r.err, isOk: false}
	}
	return f(r.value)
}

// Error returns the error, returns nil if Ok
func (r Result[T]) Error() error {
	if r.isOk {
		return nil
	}
	return r.err
}
