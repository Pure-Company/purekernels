package either

import "github.com/Pure-Company/purekernels/pkg/monoid"

// Either represents a value that can be one of two types: Left or Right
// By convention: Left = error/failure, Right = success
// This is more general than Result - L doesn't have to be an error
type Either[L, R any] struct {
	left    L
	right   R
	isRight bool
}

// Left creates a Left value (typically error)
func Left[L, R any](left L) Either[L, R] {
	return Either[L, R]{left: left, isRight: false}
}

// Right creates a Right value (typically success)
func Right[L, R any](right R) Either[L, R] {
	return Either[L, R]{right: right, isRight: true}
}

// IsLeft returns true if this is a Left value
func (e Either[L, R]) IsLeft() bool {
	return !e.isRight
}

// IsRight returns true if this is a Right value
func (e Either[L, R]) IsRight() bool {
	return e.isRight
}

// Fold applies one of two functions depending on the case
func (e Either[L, R]) Fold(onLeft func(L) any, onRight func(R) any) any {
	if e.isRight {
		return onRight(e.right)
	}
	return onLeft(e.left)
}

// Map transforms the Right value (functor)
func (e Either[L, R]) Map(f func(R) R) Either[L, R] {
	if !e.isRight {
		return Either[L, R]{left: e.left, isRight: false}
	}
	return Right[L, R](f(e.right))
}

// MapLeft transforms the Left value
func (e Either[L, R]) MapLeft(f func(L) L) Either[L, R] {
	if e.isRight {
		return Either[L, R]{right: e.right, isRight: true}
	}
	return Left[L, R](f(e.left))
}

// GetLeft returns the Left value (panics if Right)
func (e Either[L, R]) GetLeft() L {
	if e.isRight {
		panic("GetLeft called on Right")
	}
	return e.left
}

// GetRight returns the Right value (panics if Left)
func (e Either[L, R]) GetRight() R {
	if !e.isRight {
		panic("GetRight called on Left")
	}
	return e.right
}

// Get returns both values and a boolean indicating which is valid
func (e Either[L, R]) Get() (L, R, bool) {
	return e.left, e.right, e.isRight
}

// GetRightOr returns the Right value or a default
func (e Either[L, R]) GetRightOr(defaultValue R) R {
	if e.isRight {
		return e.right
	}
	return defaultValue
}

// GetLeftOr returns the Left value or a default
func (e Either[L, R]) GetLeftOr(defaultValue L) L {
	if !e.isRight {
		return e.left
	}
	return defaultValue
}

// Swap exchanges Left and Right
func (e Either[L, R]) Swap() Either[R, L] {
	if e.isRight {
		return Left[R, L](e.right)
	}
	return Right[R, L](e.left)
}

// ToOption converts Either to Option, keeping only Right values
func (e Either[L, R]) ToOption() monoid.Option[R] {
	if e.isRight {
		return monoid.Some(e.right)
	}
	return monoid.None[R]()
}

// Bimap applies functions to both sides
func Bimap[L1, L2, R1, R2 any](
	e Either[L1, R1],
	fl func(L1) L2,
	fr func(R1) R2,
) Either[L2, R2] {
	if e.isRight {
		return Right[L2, R2](fr(e.right))
	}
	return Left[L2, R2](fl(e.left))
}

// Sequence converts []Either[L,R] to Either[L, []R]
// Returns first Left, or Right with all values
func Sequence[L, R any](eithers []Either[L, R]) Either[L, []R] {
	values := make([]R, 0, len(eithers))
	for _, e := range eithers {
		if !e.isRight {
			return Left[L, []R](e.left)
		}
		values = append(values, e.right)
	}
	return Right[L, []R](values)
}

// Partition separates Lefts and Rights
func Partition[L, R any](eithers []Either[L, R]) ([]L, []R) {
	lefts := []L{}
	rights := []R{}
	for _, e := range eithers {
		if e.isRight {
			rights = append(rights, e.right)
		} else {
			lefts = append(lefts, e.left)
		}
	}
	return lefts, rights
}

// Traverse maps and sequences
func Traverse[A, L, R any](f func(A) Either[L, R], items []A) Either[L, []R] {
	values := make([]R, 0, len(items))
	for _, item := range items {
		e := f(item)
		if !e.isRight {
			return Left[L, []R](e.left)
		}
		values = append(values, e.right)
	}
	return Right[L, []R](values)
}

// pkg/either/either.go

// FlatMapEither is a standalone function that allows type transformation
// Use this instead of the method when you need to change the Right type
func FlatMapEither[L, A, B any](
	e Either[L, A],
	f func(A) Either[L, B],
) Either[L, B] {
	if !e.isRight {
		return Left[L, B](e.left)
	}
	return f(e.right)
}

// FlatMap Keep the method for convenience when types don't change
// FlatMap chains Eithers (monadic bind) - short-circuits on Left
func (e Either[L, R]) FlatMap(f func(R) Either[L, R]) Either[L, R] {
	if !e.isRight {
		return Either[L, R]{left: e.left, isRight: false}
	}
	return f(e.right)
}
