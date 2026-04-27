package effect

import "github.com/Pure-Company/purekernels/pkg/monoid"

// Writer represents a computation that produces a value and accumulates a log
// This is the Writer monad: (A, W) where W is a monoid
type Writer[W, A any] struct {
	value A
	log   W
}

// NewWriter creates a Writer with value and log
func NewWriter[W, A any](value A, log W) Writer[W, A] {
	return Writer[W, A]{value: value, log: log}
}

// PureWriter creates a Writer with just a value and empty log
func PureWriter[W, A any](m monoid.Monoid[W], value A) Writer[W, A] {
	return Writer[W, A]{value: value, log: m.Empty()}
}

// Tell creates a Writer with no value, just a log entry
func Tell[W, A any](log W) Writer[W, A] {
	var zero A
	return Writer[W, A]{value: zero, log: log}
}

// Map transforms the value
func (w Writer[W, A]) Map(f func(A) A) Writer[W, A] {
	return Writer[W, A]{value: f(w.value), log: w.log}
}

// FlatMap chains Writers, combining logs
func (w Writer[W, A]) FlatMap(m monoid.Monoid[W], f func(A) Writer[W, A]) Writer[W, A] {
	next := f(w.value)
	return Writer[W, A]{
		value: next.value,
		log:   m.Combine(w.log, next.log),
	}
}

// Run extracts value and log
func (w Writer[W, A]) Run() (A, W) {
	return w.value, w.log
}

// Value gets just the value
func (w Writer[W, A]) Value() A {
	return w.value
}

// Log gets just the log
func (w Writer[W, A]) Log() W {
	return w.log
}

// MapLog transforms the log
func (w Writer[W, A]) MapLog(f func(W) W) Writer[W, A] {
	return Writer[W, A]{value: w.value, log: f(w.log)}
}
