package monoid

// AndMonoid represents logical AND with true as identity
// Laws:
//   - Identity: Combine(x, true) == x && Combine(true, x) == x
//   - Associativity: Combine(Combine(x, y), z) == Combine(x, Combine(y, z))
type AndMonoid struct{}

// NewAndMonoid creates an AND monoid
func NewAndMonoid() AndMonoid {
	return AndMonoid{}
}

// Empty returns true (AND identity)
func (AndMonoid) Empty() bool {
	return true
}

// Combine performs logical AND
func (AndMonoid) Combine(a, b bool) bool {
	return a && b
}

// OrMonoid represents logical OR with false as identity
// Laws:
//   - Identity: Combine(x, false) == x && Combine(false, x) == x
//   - Associativity: Combine(Combine(x, y), z) == Combine(x, Combine(y, z))
type OrMonoid struct{}

// NewOrMonoid creates an OR monoid
func NewOrMonoid() OrMonoid {
	return OrMonoid{}
}

// Empty returns false (OR identity)
func (OrMonoid) Empty() bool {
	return false
}

// Combine performs logical OR
func (OrMonoid) Combine(a, b bool) bool {
	return a || b
}
