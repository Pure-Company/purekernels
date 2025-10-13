// pkg/functor/contravariant.go
package functor

// Predicate represents a contravariant predicate
// Contravariant maps in the opposite direction: (B -> A) -> F[A] -> F[B]
type Predicate[A any] struct {
	test func(A) bool
}

// NewPredicate creates a predicate
func NewPredicate[A any](test func(A) bool) Predicate[A] {
	return Predicate[A]{test: test}
}

// Test runs the predicate
func (p Predicate[A]) Test(a A) bool {
	return p.test(a)
}

// Contramap transforms the input (contravariant functor)
// Given f: B -> A and Predicate[A], produces Predicate[B]
func Contramap[A, B any](p Predicate[A], f func(B) A) Predicate[B] {
	return Predicate[B]{
		test: func(b B) bool {
			return p.test(f(b))
		},
	}
}

// Comparison represents a contravariant comparison
type Comparison[A any] struct {
	compare func(A, A) int
}

// NewComparison creates a comparison
func NewComparison[A any](compare func(A, A) int) Comparison[A] {
	return Comparison[A]{compare: compare}
}

// Compare runs the comparison
func (c Comparison[A]) Compare(a1, a2 A) int {
	return c.compare(a1, a2)
}

// ContramapComparison transforms the input (contravariant functor)
// Given f: B -> A and Comparison[A], produces Comparison[B]
func ContramapComparison[A, B any](c Comparison[A], f func(B) A) Comparison[B] {
	return Comparison[B]{
		compare: func(b1, b2 B) int {
			return c.compare(f(b1), f(b2))
		},
	}
}

// And combines two predicates with logical AND
func And[A any](p1, p2 Predicate[A]) Predicate[A] {
	return Predicate[A]{
		test: func(a A) bool {
			return p1.test(a) && p2.test(a)
		},
	}
}

// Or combines two predicates with logical OR
func Or[A any](p1, p2 Predicate[A]) Predicate[A] {
	return Predicate[A]{
		test: func(a A) bool {
			return p1.test(a) || p2.test(a)
		},
	}
}

// Not negates a predicate
func Not[A any](p Predicate[A]) Predicate[A] {
	return Predicate[A]{
		test: func(a A) bool {
			return !p.test(a)
		},
	}
}
