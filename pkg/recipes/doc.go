// Package recipes pkg/recipes/doc.go
// Package recipes provides cookbook examples using purekernels primitives.
//
// This package contains NO code - only documentation showing how to compose
// the library's abstractions for common patterns.
//
// # Parallel Map-Reduce
//
// Use TraverseConcurrent with a monoid for parallel aggregation:
//
//	sumMonoid := monoid.NewSumMonoid[int]()
//	result := functor.TraverseConcurrent(
//	    sumMonoid,
//	    expensiveComputation,
//	    items,
//	    0, // use all CPUs
//	).Value()
//
// # Retry with Exponential Backoff
//
// Compose Task with recursion:
//
//	func withRetry[A any](f func() result.Result[A], attempts int) effect.Task[A] {
//	    return effect.NewTask(func() result.Result[A] {
//	        res := f()
//	        if res.IsOk() || attempts <= 1 {
//	            return res
//	        }
//	        time.Sleep(time.Second * time.Duration(attempts))
//	        return withRetry(f, attempts-1).Run()
//	    })
//	}
//
// # Pipeline with Error Handling
//
// Chain operations using Kleisli composition:
//
//	pipeline := compose.KleisliResult(
//	    compose.KleisliResult(parse, validate),
//	    transform,
//	)
//	result := pipeline(input)
//
// # Validation Pipeline
//
// Use Validation with Ap3 for error accumulation:
//
//	import "github.com/vinodhalaharvi/purekernels/pkg/validation"
//
//	result := either.Ap3(
//	    validation.ErrorsMonoid,
//	    NewUser,
//	    validateName(input.Name),
//	    validateEmail(input.Email),
//	    validateAge(input.Age),
//	)
//
// # Context-Aware Operations
//
//	import "github.com/vinodhalaharvi/purekernels/pkg/bridges"
//
//	userID := bridges.FromContext[string](ctx, "user_id")
//	userID.Map(func(id string) { /* use id */ })
//
// # Database Queries
//
//	users := bridges.QueryAll(db, "SELECT * FROM users", func(rows *sql.Rows) (User, error) {
//	    var u User
//	    err := rows.Scan(&u.ID, &u.Name)
//	    return u, err
//	})
//
// # Batch Processing
//
//	batchMonoid := monoid.NewListMonoid[Result]()
//	results := functor.ConcurrentBatch(
//	    batchMonoid,
//	    100, // batch size
//	    processBatch,
//	    allItems,
//	).Value()
package recipes
