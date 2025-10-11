// Package functor provides concurrent and compositional computation patterns
// based on category theory.
//
// # Core Types
//
// This package provides two concrete applicative types:
//
//   - Concurrent[T]: Parallel computation using goroutines
//   - Sequential[T]: Sequential computation (left-to-right evaluation)
//
// Both types accumulate results using a Monoid[T] and follow applicative laws.
//
// # Why No Interface?
//
// Earlier versions used an Applicative[T] interface, but this caused type erasure:
//
//	var c Concurrent[T] = NewConcurrent(...)
//	result := c.Apply(other)  // Returns Concurrent[T] ✓
//
//	var a Applicative[T] = c
//	result := a.Apply(other)  // Returns Applicative[T] ✗ (lost concrete type!)
//
// Go's generics don't support higher-kinded types (F[_]), so interfaces
// erase the important type-constructor information. By using concrete types:
//
//   - No type erasure
//   - Better method chaining
//   - More idiomatic Go
//   - The algebra is rigid anyway - only one lawful implementation
//
// # Mathematical Laws
//
// Both Concurrent and Sequential satisfy applicative laws:
//
//	Identity:     Pure(empty).Apply(v) ≡ v
//	Composition:  u.Apply(v).Apply(w) ≡ u.Apply(v.Apply(w))
//	Homomorphism: Pure(x).Apply(Pure(y)) ≡ Pure(combine(x,y))
//
// These laws ensure:
//   - Computations compose predictably
//   - Order doesn't affect correctness (associativity)
//   - Pure values behave as identity elements
//
// # Basic Usage: Concurrent
//
//	// Create a monoid for your type
//	setMonoid := monoid.NewSetMonoid[string]()
//
//	// Create parallel computations
//	c1 := functor.NewConcurrent(setMonoid, func() monoid.SetMonoid[string] {
//	    // expensive computation 1
//	    return result1
//	})
//
//	c2 := functor.NewConcurrent(setMonoid, func() monoid.SetMonoid[string] {
//	    // expensive computation 2
//	    return result2
//	})
//
//	// Combine them (runs in parallel!)
//	combined := c1.Apply(c2)
//	result := combined.Value()
//
// # Basic Usage: Sequential
//
//	// For deterministic, ordered execution
//	s1 := functor.NewSequential(setMonoid, func() monoid.SetMonoid[string] {
//	    return result1
//	})
//
//	s2 := functor.NewSequential(setMonoid, func() monoid.SetMonoid[string] {
//	    return result2
//	})
//
//	// Combines left-to-right sequentially
//	result := s1.Apply(s2).Value()
//
// # Higher-Level Operations
//
//	// Parallel processing with worker pool
//	result := functor.TraverseConcurrent(
//	    myMonoid,
//	    expensiveFunc,
//	    items,
//	    0, // use runtime.NumCPU() workers
//	).Value()
//
//	// Batch processing
//	result := functor.ConcurrentBatch(
//	    myMonoid,
//	    100, // batch size
//	    batchProcessor,
//	    items,
//	).Value()
//
//	// Sequential fold
//	result := functor.FoldSequential(
//	    myMonoid,
//	    processFunc,
//	    items,
//	).Value()
//
// # Choosing Between Concurrent and Sequential
//
// Use Concurrent when:
//   - Computations are independent
//   - You want parallelism
//   - Order doesn't matter for correctness (monoid is commutative)
//
// Use Sequential when:
//   - You need deterministic ordering
//   - Side effects matter
//   - Single-threaded execution is required
//
// Both use the same monoid-based composition, just different execution strategies.
package functor
