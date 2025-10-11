package effect

// State represents a stateful computation: S -> (A, S)
// Threads state through a computation
type State[S, A any] func(S) (A, S)

// PureState creates a State that returns a value without modifying state
func PureState[S, A any](a A) State[S, A] {
	return func(s S) (A, S) {
		return a, s
	}
}

// Get creates a State that returns the current state
func Get[S any]() State[S, S] {
	return func(s S) (S, S) {
		return s, s
	}
}

// Put creates a State that sets new state
func Put[S any](newState S) State[S, struct{}] {
	return func(S) (struct{}, S) {
		return struct{}{}, newState
	}
}

// Modify creates a State that modifies state
func Modify[S any](f func(S) S) State[S, struct{}] {
	return func(s S) (struct{}, S) {
		return struct{}{}, f(s)
	}
}

// Map transforms the result value
func (st State[S, A]) Map(f func(A) A) State[S, A] {
	return func(s S) (A, S) {
		a, s2 := st(s)
		return f(a), s2
	}
}

// FlatMap chains stateful computations
func (st State[S, A]) FlatMap(f func(A) State[S, A]) State[S, A] {
	return func(s S) (A, S) {
		a, s2 := st(s)
		return f(a)(s2)
	}
}

// Run executes the stateful computation with initial state
func (st State[S, A]) Run(initial S) (A, S) {
	return st(initial)
}

// Eval runs and returns only the value
func (st State[S, A]) Eval(initial S) A {
	a, _ := st(initial)
	return a
}

// Exec runs and returns only the final state
func (st State[S, A]) Exec(initial S) S {
	_, s := st(initial)
	return s
}
