package functor

import (
	"runtime"
	"sync"

	"github.com/vinodhalaharvi/purekernels/pkg/monoid"
)

// Concurrent is a parallel computation that accumulates values using a monoid
// This follows the Applicative pattern but as a concrete type
//
// Laws (for mathematical correctness):
//   - Identity: Pure(empty).Apply(v) ≡ v
//   - Composition: u.Apply(v).Apply(w) ≡ u.Apply(v.Apply(w))
//   - Homomorphism: Pure(x).Apply(Pure(y)) ≡ Pure(monoid.Combine(x, y))
type Concurrent[T any] struct {
	compute func() T
	monoid  monoid.Monoid[T]
}

// NewConcurrent creates a concurrent computation
func NewConcurrent[T any](m monoid.Monoid[T], compute func() T) Concurrent[T] {
	return Concurrent[T]{
		compute: compute,
		monoid:  m,
	}
}

// Pure wraps a value in a concurrent computation (identity)
func Pure[T any](m monoid.Monoid[T], value T) Concurrent[T] {
	return Concurrent[T]{
		compute: func() T { return value },
		monoid:  m,
	}
}

// Apply combines two concurrent computations in parallel
func (c Concurrent[T]) Apply(other Concurrent[T]) Concurrent[T] {
	return Concurrent[T]{
		compute: func() T {
			var left, right T
			var wg sync.WaitGroup
			wg.Add(2)

			// Run both computations in parallel
			go func() {
				defer wg.Done()
				left = c.compute()
			}()

			go func() {
				defer wg.Done()
				right = other.compute()
			}()

			wg.Wait()
			return c.monoid.Combine(left, right)
		},
		monoid: c.monoid,
	}
}

// Map transforms the result of a computation
func (c Concurrent[T]) Map(f func(T) T) Concurrent[T] {
	return Concurrent[T]{
		compute: func() T {
			return f(c.compute())
		},
		monoid: c.monoid,
	}
}

// Value executes the computation and returns the result
func (c Concurrent[T]) Value() T {
	return c.compute()
}

// CombineAll combines multiple concurrent computations
func CombineAll[T any](first Concurrent[T], rest ...Concurrent[T]) Concurrent[T] {
	result := first
	for _, comp := range rest {
		result = result.Apply(comp)
	}
	return result
}

// Sequence executes multiple computations and extracts their values
func Sequence[T any](comps []Concurrent[T]) []T {
	result := make([]T, len(comps))
	for i, comp := range comps {
		result[i] = comp.Value()
	}
	return result
}

// SequenceParallel executes multiple computations in parallel and extracts their values
func SequenceParallel[T any](comps []Concurrent[T]) []T {
	result := make([]T, len(comps))
	var wg sync.WaitGroup
	wg.Add(len(comps))

	for i := range comps {
		go func(idx int) {
			defer wg.Done()
			result[idx] = comps[idx].Value()
		}(i)
	}

	wg.Wait()
	return result
}

// Traverse maps a function over items and combines results using a monoid
func Traverse[A any, T any](
	m monoid.Monoid[T],
	f func(A) Concurrent[T],
	items []A,
) Concurrent[T] {
	if len(items) == 0 {
		return Pure(m, m.Empty())
	}

	result := f(items[0])
	for i := 1; i < len(items); i++ {
		result = result.Apply(f(items[i]))
	}
	return result
}

// TraverseConcurrent processes items in parallel using a worker pool
// workers=0 means use runtime.NumCPU()
func TraverseConcurrent[A any, T any](
	m monoid.Monoid[T],
	f func(A) T,
	items []A,
	workers int,
) Concurrent[T] {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	return Concurrent[T]{
		compute: func() T {
			if len(items) == 0 {
				return m.Empty()
			}

			results := make(chan T, len(items))
			jobs := make(chan A, len(items))

			// Start worker pool
			var wg sync.WaitGroup
			for i := 0; i < workers; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for item := range jobs {
						results <- f(item)
					}
				}()
			}

			// Send jobs
			for _, item := range items {
				jobs <- item
			}
			close(jobs)

			// Wait for completion
			go func() {
				wg.Wait()
				close(results)
			}()

			// Combine results
			result := m.Empty()
			for r := range results {
				result = m.Combine(result, r)
			}

			return result
		},
		monoid: m,
	}
}

// ParMap is a parallel map that returns a slice of results
func ParMap[A any, T any](
	f func(A) T,
	items []A,
) []T {
	return parMapWithWorkers(f, items, runtime.NumCPU())
}

// ParMapWithWorkers is ParMap with configurable worker count
func ParMapWithWorkers[A any, T any](
	f func(A) T,
	items []A,
	workers int,
) []T {
	return parMapWithWorkers(f, items, workers)
}

func parMapWithWorkers[A any, T any](
	f func(A) T,
	items []A,
	workers int,
) []T {
	if len(items) == 0 {
		return []T{}
	}

	results := make([]T, len(items))
	jobs := make(chan int, len(items))

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				results[idx] = f(items[idx])
			}
		}()
	}

	for i := range items {
		jobs <- i
	}
	close(jobs)

	wg.Wait()
	return results
}

// ConcurrentBatch processes items in batches concurrently
func ConcurrentBatch[A any, T any](
	m monoid.Monoid[T],
	batchSize int,
	f func([]A) T,
	items []A,
) Concurrent[T] {
	return Concurrent[T]{
		compute: func() T {
			if len(items) == 0 {
				return m.Empty()
			}

			// Split into batches
			batches := [][]A{}
			for i := 0; i < len(items); i += batchSize {
				end := i + batchSize
				if end > len(items) {
					end = len(items)
				}
				batches = append(batches, items[i:end])
			}

			// Process batches in parallel
			results := make(chan T, len(batches))
			var wg sync.WaitGroup

			for _, batch := range batches {
				wg.Add(1)
				go func(b []A) {
					defer wg.Done()
					results <- f(b)
				}(batch)
			}

			go func() {
				wg.Wait()
				close(results)
			}()

			// Combine results
			result := m.Empty()
			for r := range results {
				result = m.Combine(result, r)
			}

			return result
		},
		monoid: m,
	}
}

// ApplyIf conditionally applies a computation
func ApplyIf[T any](
	condition bool,
	base Concurrent[T],
	conditional Concurrent[T],
) Concurrent[T] {
	if condition {
		return base.Apply(conditional)
	}
	return base
}

// ApplyWhen applies a computation if predicate is true
func ApplyWhen[T any](
	predicate func() bool,
	base Concurrent[T],
	conditional Concurrent[T],
) Concurrent[T] {
	return ApplyIf(predicate(), base, conditional)
}
