package monoid

// Avg represents a running average as a composite of sum and count
type Avg struct {
	Sum   float64
	Count float64
}

// AvgMonoid is a composite monoid for computing averages
// This demonstrates how monoids can be composed to create derived statistics
// Laws:
//   - Identity: Combine(x, {0,0}) == x && Combine({0,0}, x) == x
//   - Associativity: Combine(Combine(x, y), z) == Combine(x, Combine(y, z))
type AvgMonoid struct{}

// NewAvgMonoid creates an average monoid
func NewAvgMonoid() AvgMonoid {
	return AvgMonoid{}
}

// Empty returns zero sum and zero count
func (AvgMonoid) Empty() Avg {
	return Avg{Sum: 0, Count: 0}
}

// Combine adds sums and counts
func (AvgMonoid) Combine(a, b Avg) Avg {
	return Avg{
		Sum:   a.Sum + b.Sum,
		Count: a.Count + b.Count,
	}
}

// Value computes the average (returns 0 if count is 0)
func (a Avg) Value() float64 {
	if a.Count == 0 {
		return 0
	}
	return a.Sum / a.Count
}

// FromValue creates an Avg from a single value
func FromValue(v float64) Avg {
	return Avg{Sum: v, Count: 1}
}

// FromValues creates an Avg from multiple values
func FromValues(values ...float64) Avg {
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return Avg{Sum: sum, Count: float64(len(values))}
}
