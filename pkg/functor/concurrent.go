package functor

import (
	"runtime"
	"sync"

	"github.com/vinodhalaharvi/purekernels/pkg/monoid"
)

// Concurrent is a parallel applicative implementation
// Perfect for Go's concurrency model!
type Concurrent[T any] struct {
	compute func() T
	monoid  monoid.Monoid[T]
}

// NewConcurrent creates a concurrent applicative
func NewConcurrent[T any](m monoid.Monoid[T], compute func() T) Concurrent[T] {
	return Concurrent[T]{
		compute: compute,
		monoid:  m,
	}
}

// NewConcurrentValue wraps a value
func NewConcurrentValue[T any](m monoid.Monoid[T], value T) Concurrent[T] {
	return Concurrent[T]{
		compute: func() T { return value },
		monoid:  m,
	}
}

// Apply combines two concurrent computations in parallel
func (c Concurrent[T]) Apply(other Applicative[T]) Concurrent[T] {
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
				right = other.Value()
			}()

			wg.Wait()
			return c.monoid.Combine(left, right)
		},
		monoid: c.monoid,
	}
}

// Map transforms the result
func (c Concurrent[T]) Map(f func(T) T) Concurrent[T] {
	return Concurrent[T]{
		compute: func() T {
			return f(c.compute())
		},
		monoid: c.monoid,
	}
}

// Value executes the computation and returns result
func (c Concurrent[T]) Value() T {
	return c.compute()
}

// TraverseConcurrent processes items in parallel
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

// ParMap is a parallel map using concurrent applicative
func ParMap[A any, T any](
	m monoid.Monoid[T],
	f func(A) T,
	items []A,
) []T {
	return parMapWithWorkers(m, f, items, runtime.NumCPU())
}

func parMapWithWorkers[A any, T any](
	m monoid.Monoid[T],
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

// ConcurrentBatch processes in batches concurrently
func ConcurrentBatch[A any, T any](
	m monoid.Monoid[T],
	batchSize int,
	f func([]A) T,
	items []A,
) Concurrent[T] {
	return Concurrent[T]{
		compute: func() T {
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
