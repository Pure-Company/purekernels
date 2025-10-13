package effect

import "sync"

// Lazy represents deferred computation with memoization
type Lazy[A any] struct {
	compute func() A
	cached  *A
	mu      sync.Mutex
}

// NewLazy creates a lazy computation
func NewLazy[A any](f func() A) *Lazy[A] {
	return &Lazy[A]{compute: f}
}

// Force evaluates the computation (memoized)
func (l *Lazy[A]) Force() A {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.cached == nil {
		value := l.compute()
		l.cached = &value
	}
	return *l.cached
}

// Map transforms the result lazily
func (l *Lazy[A]) Map(f func(A) A) *Lazy[A] {
	return NewLazy(func() A {
		return f(l.Force())
	})
}

// FlatMap chains lazy computations
func (l *Lazy[A]) FlatMap(f func(A) *Lazy[A]) *Lazy[A] {
	return NewLazy(func() A {
		return f(l.Force()).Force()
	})
}

// Value is an alias for Force (more idiomatic)
func (l *Lazy[A]) Value() A {
	return l.Force()
}
