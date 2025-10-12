// Package either pkg/either/doc.go
// Package either provides sum types for branching logic and validation.
//
// # Either Type
//
// Either[L, R] represents a value that can be one of two types.
// By convention, Left represents failure and Right represents success.
// This is more general than Result - L doesn't have to be an error type.
//
// Basic usage:
//
//	func parseAge(s string) either.Either[string, int] {
//	    age, err := strconv.Atoi(s)
//	    if err != nil {
//	        return either.Left[string, int]("invalid age")
//	    }
//	    if age < 0 || age > 150 {
//	        return either.Left[string, int]("age out of range")
//	    }
//	    return either.Right[string, int](age)
//	}
//
// # Validation Type - Error Accumulation
//
// Validation[E, A] is an applicative functor that accumulates errors
// instead of short-circuiting like Either. This is crucial for form
// validation where you want to show ALL errors, not just the first one.
//
// Key difference: Validation is NOT a monad (no FlatMap that obeys laws)
// because error accumulation breaks the monad laws. But it IS an applicative!
//
// Basic validation example:
//
//	// Validate a user registration form
//	type User struct {
//	    Name  string
//	    Email string
//	    Age   int
//	}
//
//	// Use list monoid to accumulate error strings
//	errorMonoid := monoid.NewListMonoid[string]()
//
//	func validateName(name string) either.Validation[[]string, string] {
//	    if len(name) < 2 {
//	        return either.Invalid(errorMonoid, []string{"name too short"})
//	    }
//	    return either.Valid(errorMonoid, name)
//	}
//
//	func validateEmail(email string) either.Validation[[]string, string] {
//	    if !strings.Contains(email, "@") {
//	        return either.Invalid(errorMonoid, []string{"invalid email"})
//	    }
//	    return either.Valid(errorMonoid, email)
//	}
//
//	func validateAge(age int) either.Validation[[]string, int] {
//	    if age < 18 {
//	        return either.Invalid(errorMonoid, []string{"must be 18+"})
//	    }
//	    return either.Valid(errorMonoid, age)
//	}
//
//	// Combine validations - accumulates ALL errors
//	func validateUser(name string, email string, age int) either.Validation[[]string, User] {
//	    return either.Ap3(
//	        errorMonoid,
//	        func(n string, e string, a int) User {
//	            return User{Name: n, Email: e, Age: a}
//	        },
//	        validateName(name),
//	        validateEmail(email),
//	        validateAge(age),
//	    )
//	}
//
//	// Usage:
//	result := validateUser("X", "bademail", 15)
//	// Returns Invalid with ALL three errors:
//	// ["name too short", "invalid email", "must be 18+"]
//
// # When to Use What
//
// Use Either when:
//   - You want to short-circuit on first error (fail-fast)
//   - You're doing sequential operations where later steps depend on earlier ones
//   - You need monadic composition (FlatMap/chaining)
//
// Use Validation when:
//   - You want to accumulate ALL errors before failing
//   - You're validating independent fields (like form inputs)
//   - You need applicative composition (parallel validation)
//   - Errors form a monoid (can be combined)
//
// Use Result when:
//   - Your errors are Go errors
//   - You want Either[error, T] specifically
//
// # Mathematical Background
//
// Either is a Monad:
//   - Has lawful Map, FlatMap, and Pure
//   - Short-circuits on Left (error)
//   - Sequential composition
//
// Validation is an Applicative (but NOT a Monad):
//   - Has lawful Map, Apply, and Pure
//   - Accumulates errors (breaks monad laws)
//   - Parallel composition
//
// The key insight: You can't have both error accumulation AND
// lawful monadic composition. Validation chooses accumulation.
package either
