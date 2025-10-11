// Package effect provides effect types for contextual computations
package effect

// Reader represents a computation that depends on an environment R
// This is the Reader monad: (R -> A)
type Reader[R, A any] func(R) A

// Pure creates a Reader that ignores environment and returns a constant
func PureReader[R, A any](a A) Reader[R, A] {
	return func(R) A { return a }
}

// Ask creates a Reader that returns the environment
func Ask[R any]() Reader[R, R] {
	return func(r R) R { return r }
}

// Map transforms the result of a Reader
func (r Reader[R, A]) Map(f func(A) A) Reader[R, A] {
	return func(env R) A {
		return f(r(env))
	}
}

// FlatMap chains Readers (monadic bind)
func (r Reader[R, A]) FlatMap(f func(A) Reader[R, A]) Reader[R, A] {
	return func(env R) A {
		a := r(env)
		return f(a)(env)
	}
}

// Run executes the Reader with an environment
func (r Reader[R, A]) Run(env R) A {
	return r(env)
}

// Local modifies the environment before running
func (r Reader[R, A]) Local(f func(R) R) Reader[R, A] {
	return func(env R) A {
		return r(f(env))
	}
}
