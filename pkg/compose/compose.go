// Package compose provides function composition and pipeline operations
package compose

// Compose composes two functions: (g ∘ f)(x) = g(f(x))
func Compose[A, B, C any](g func(B) C, f func(A) B) func(A) C {
	return func(a A) C {
		return g(f(a))
	}
}

// Pipe chains multiple functions left to right
// Usage: Pipe(f, g, h)(x) == h(g(f(x)))
func Pipe[T any](fns ...func(T) T) func(T) T {
	return func(initial T) T {
		result := initial
		for _, fn := range fns {
			result = fn(result)
		}
		return result
	}
}

// Pipe2 pipes with type transformations A -> B -> C
func Pipe2[A, B, C any](f func(A) B, g func(B) C) func(A) C {
	return func(a A) C {
		return g(f(a))
	}
}

// Pipe3 pipes A -> B -> C -> D
func Pipe3[A, B, C, D any](
	f func(A) B,
	g func(B) C,
	h func(C) D,
) func(A) D {
	return func(a A) D {
		return h(g(f(a)))
	}
}

// Identity returns the input unchanged
func Identity[T any](x T) T {
	return x
}

// Const returns a function that always returns the same value
func Const[A, B any](b B) func(A) B {
	return func(_ A) B {
		return b
	}
}

// Flip swaps the order of arguments
func Flip[A, B, C any](f func(A, B) C) func(B, A) C {
	return func(b B, a A) C {
		return f(a, b)
	}
}

// Curry converts a 2-arg function to curried form
func Curry[A, B, C any](f func(A, B) C) func(A) func(B) C {
	return func(a A) func(B) C {
		return func(b B) C {
			return f(a, b)
		}
	}
}

// Uncurry converts curried function to 2-arg form
func Uncurry[A, B, C any](f func(A) func(B) C) func(A, B) C {
	return func(a A, b B) C {
		return f(a)(b)
	}
}
