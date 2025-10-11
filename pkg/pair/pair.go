// Package pair provides product types and their monoid instances
package pair

import "github.com/vinodhalaharvi/purekernels/pkg/monoid"

// Pair represents a product of two values (A, B)
type Pair[A, B any] struct {
	First  A
	Second B
}

// NewPair creates a new pair
func NewPair[A, B any](a A, b B) Pair[A, B] {
	return Pair[A, B]{First: a, Second: b}
}

// Map transforms the first component
func (p Pair[A, B]) MapFirst(f func(A) A) Pair[A, B] {
	return Pair[A, B]{First: f(p.First), Second: p.Second}
}

// MapSecond transforms the second component
func (p Pair[A, B]) MapSecond(f func(B) B) Pair[A, B] {
	return Pair[A, B]{First: p.First, Second: f(p.Second)}
}

// MapBoth transforms both components
func (p Pair[A, B]) MapBoth(fa func(A) A, fb func(B) B) Pair[A, B] {
	return Pair[A, B]{First: fa(p.First), Second: fb(p.Second)}
}

// Swap exchanges first and second
func (p Pair[A, B]) Swap() Pair[B, A] {
	return Pair[B, A]{First: p.Second, Second: p.First}
}

// PairMonoid combines pairs using monoids for each component
// Laws:
//   - Identity: Combine(p, empty) == p
//   - Associativity: Combine(Combine(p1, p2), p3) == Combine(p1, Combine(p2, p3))
type PairMonoid[A, B any] struct {
	monoidA monoid.Monoid[A]
	monoidB monoid.Monoid[B]
}

// NewPairMonoid creates a monoid for pairs
func NewPairMonoid[A, B any](ma monoid.Monoid[A], mb monoid.Monoid[B]) PairMonoid[A, B] {
	return PairMonoid[A, B]{monoidA: ma, monoidB: mb}
}

// Empty returns the identity pair (empty, empty)
func (m PairMonoid[A, B]) Empty() Pair[A, B] {
	return Pair[A, B]{
		First:  m.monoidA.Empty(),
		Second: m.monoidB.Empty(),
	}
}

// Combine combines pairs component-wise
func (m PairMonoid[A, B]) Combine(p1, p2 Pair[A, B]) Pair[A, B] {
	return Pair[A, B]{
		First:  m.monoidA.Combine(p1.First, p2.First),
		Second: m.monoidB.Combine(p1.Second, p2.Second),
	}
}

// Triple represents a product of three values
type Triple[A, B, C any] struct {
	First  A
	Second B
	Third  C
}

// NewTriple creates a new triple
func NewTriple[A, B, C any](a A, b B, c C) Triple[A, B, C] {
	return Triple[A, B, C]{First: a, Second: b, Third: c}
}

// TripleMonoid combines triples using monoids for each component
type TripleMonoid[A, B, C any] struct {
	monoidA monoid.Monoid[A]
	monoidB monoid.Monoid[B]
	monoidC monoid.Monoid[C]
}

// NewTripleMonoid creates a monoid for triples
func NewTripleMonoid[A, B, C any](
	ma monoid.Monoid[A],
	mb monoid.Monoid[B],
	mc monoid.Monoid[C],
) TripleMonoid[A, B, C] {
	return TripleMonoid[A, B, C]{monoidA: ma, monoidB: mb, monoidC: mc}
}

// Empty returns the identity triple
func (m TripleMonoid[A, B, C]) Empty() Triple[A, B, C] {
	return Triple[A, B, C]{
		First:  m.monoidA.Empty(),
		Second: m.monoidB.Empty(),
		Third:  m.monoidC.Empty(),
	}
}

// Combine combines triples component-wise
func (m TripleMonoid[A, B, C]) Combine(t1, t2 Triple[A, B, C]) Triple[A, B, C] {
	return Triple[A, B, C]{
		First:  m.monoidA.Combine(t1.First, t2.First),
		Second: m.monoidB.Combine(t1.Second, t2.Second),
		Third:  m.monoidC.Combine(t1.Third, t2.Third),
	}
}
