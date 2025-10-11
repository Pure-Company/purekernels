package monoid

import (
	"math"
	"testing"
	"testing/quick"
)

// testMonoidLaws verifies identity and associativity for any monoid
func testMonoidLaws[T any](t *testing.T, m Monoid[T], gen func() T, eq func(T, T) bool) {
	t.Helper()

	// Law 1: Left Identity - Combine(Empty(), x) == x
	t.Run("LeftIdentity", func(t *testing.T) {
		f := func() bool {
			x := gen()
			result := m.Combine(m.Empty(), x)
			return eq(result, x)
		}
		if err := quick.Check(f, nil); err != nil {
			t.Errorf("Left identity law failed: %v", err)
		}
	})

	// Law 2: Right Identity - Combine(x, Empty()) == x
	t.Run("RightIdentity", func(t *testing.T) {
		f := func() bool {
			x := gen()
			result := m.Combine(x, m.Empty())
			return eq(result, x)
		}
		if err := quick.Check(f, nil); err != nil {
			t.Errorf("Right identity law failed: %v", err)
		}
	})

	// Law 3: Associativity - Combine(Combine(x, y), z) == Combine(x, Combine(y, z))
	t.Run("Associativity", func(t *testing.T) {
		f := func() bool {
			x, y, z := gen(), gen(), gen()
			left := m.Combine(m.Combine(x, y), z)
			right := m.Combine(x, m.Combine(y, z))
			return eq(left, right)
		}
		if err := quick.Check(f, nil); err != nil {
			t.Errorf("Associativity law failed: %v", err)
		}
	})
}

// Test generators for different types
func genInt() int {
	return int(quick.Value(nil, nil).Interface().(int))
}

func genFloat() float64 {
	return quick.Value(nil, nil).Interface().(float64)
}

func genBool() bool {
	return quick.Value(nil, nil).Interface().(bool)
}

func genString() string {
	return quick.Value(nil, nil).Interface().(string)
}

func genIntSlice() []int {
	return quick.Value(nil, nil).Interface().([]int)
}

// Equality functions
func intEq(a, b int) bool       { return a == b }
func boolEq(a, b bool) bool     { return a == b }
func stringEq(a, b string) bool { return a == b }

func floatEq(a, b float64) bool {
	if math.IsNaN(a) && math.IsNaN(b) {
		return true
	}
	if math.IsInf(a, 0) && math.IsInf(b, 0) {
		return math.Signbit(a) == math.Signbit(b)
	}
	return math.Abs(a-b) < 1e-9
}

func sliceEq[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Monoid law tests

func TestSumMonoidLaws(t *testing.T) {
	m := NewSumMonoid[int]()
	testMonoidLaws(t, m, genInt, intEq)
}

func TestProductMonoidLaws(t *testing.T) {
	m := NewProductMonoid[int]()
	testMonoidLaws(t, m, genInt, intEq)
}

func TestAndMonoidLaws(t *testing.T) {
	m := NewAndMonoid()
	testMonoidLaws(t, m, genBool, boolEq)
}

func TestOrMonoidLaws(t *testing.T) {
	m := NewOrMonoid()
	testMonoidLaws(t, m, genBool, boolEq)
}

func TestStringConcatMonoidLaws(t *testing.T) {
	m := NewStringConcatMonoid()
	testMonoidLaws(t, m, genString, stringEq)
}

func TestListMonoidLaws(t *testing.T) {
	m := NewListMonoid[int]()
	testMonoidLaws(t, m, genIntSlice, sliceEq[int])
}

func TestMaxMonoidLaws(t *testing.T) {
	m := NewMaxMonoid[int](math.MinInt)
	testMonoidLaws(t, m, genInt, intEq)
}

func TestMinMonoidLaws(t *testing.T) {
	m := NewMinMonoid[int](math.MaxInt)
	testMonoidLaws(t, m, genInt, intEq)
}

func TestAvgMonoidLaws(t *testing.T) {
	m := NewAvgMonoid()
	gen := func() Avg {
		return Avg{
			Sum:   genFloat(),
			Count: math.Abs(genFloat()),
		}
	}
	eq := func(a, b Avg) bool {
		return floatEq(a.Sum, b.Sum) && floatEq(a.Count, b.Count)
	}
	testMonoidLaws(t, m, gen, eq)
}

func TestDualMonoidLaws(t *testing.T) {
	inner := NewSumMonoid[int]()
	m := NewDual(inner)
	testMonoidLaws(t, m, genInt, intEq)
}

func TestSetMonoidLaws(t *testing.T) {
	m := NewSetMonoid[int]()
	gen := func() SetMonoid[int] {
		s := NewSetMonoid[int]()
		slice := genIntSlice()
		for _, v := range slice {
			s = s.Insert(v)
		}
		return s
	}
	eq := func(a, b SetMonoid[int]) bool {
		if a.Size() != b.Size() {
			return false
		}
		for _, v := range a.ToSlice() {
			if !b.Contains(v) {
				return false
			}
		}
		return true
	}
	testMonoidLaws(t, m, gen, eq)
}

// Demonstrate dual reverses order
func TestDualReversesOrder(t *testing.T) {
	concat := NewStringConcatMonoid()
	dual := NewDual[string](concat)

	a, b := "Hello", "World"

	normal := concat.Combine(a, b)
	reversed := dual.Combine(a, b)

	if normal != "HelloWorld" {
		t.Errorf("Expected 'HelloWorld', got %s", normal)
	}
	if reversed != "WorldHello" {
		t.Errorf("Expected 'WorldHello', got %s", reversed)
	}
}
