// Package result provides error handling via sum types instead of exceptions.
//
// # Result Type
//
// Result[T] represents a computation that can either succeed with a value
// or fail with an error, without using exceptions or panic/recover.
//
// Basic usage:
//
//	func divide(a, b float64) result.Result[float64] {
//	    if b == 0 {
//	        return result.Err[float64](errors.New("division by zero"))
//	    }
//	    return result.Ok(a / b)
//	}
//
//	res := divide(10, 2)
//	value := res.UnwrapOr(0)  // 5.0
//
// # Chaining Operations
//
//	result.Ok(5).
//	    Map(func(x int) int { return x * 2 }).
//	    FlatMap(func(x int) result.Result[int] {
//	        if x > 20 {
//	            return result.Err[int](errors.New("too large"))
//	        }
//	        return result.Ok(x)
//	    }).
//	    UnwrapOr(0)
//
// # Collecting Results
//
//	results := []result.Result[int]{
//	    result.Ok(1),
//	    result.Ok(2),
//	    result.Ok(3),
//	}
//
//	collected := result.Collect(results)  // Result[[]int]
//	if collected.IsOk() {
//	    values := collected.Unwrap()  // []int{1, 2, 3}
//	}
package result
