package monoid_test

import (
	"fmt"
	"testing"

	"github.com/vinodhalaharvi/purekernels/pkg/monoid"
)

func ExampleSumMonoid() {
	m := monoid.NewSumMonoid[int]()
	result := monoid.Reduce(m, []int{1, 2, 3, 4, 5})
	fmt.Println(result)
	// Output: 15
}

func ExampleProductMonoid() {
	m := monoid.NewProductMonoid[int]()
	result := monoid.Reduce(m, []int{2, 3, 4})
	fmt.Println(result)
	// Output: 24
}

func ExampleAndMonoid() {
	m := monoid.NewAndMonoid()
	result := monoid.Reduce(m, []bool{true, true, false})
	fmt.Println(result)
	// Output: false
}

func ExampleOrMonoid() {
	m := monoid.NewOrMonoid()
	result := monoid.Reduce(m, []bool{false, false, true})
	fmt.Println(result)
	// Output: true
}

func ExampleStringConcatMonoid() {
	m := monoid.NewStringConcatMonoid()
	result := monoid.Reduce(m, []string{"Hello", " ", "World"})
	fmt.Println(result)
	// Output: Hello World
}

func ExampleStringJoinMonoid() {
	m := monoid.NewStringJoinMonoid(", ")
	result := monoid.Reduce(m, []string{"apple", "banana", "cherry"})
	fmt.Println(result)
	// Output: apple, banana, cherry
}

func ExampleListMonoid() {
	m := monoid.NewListMonoid[int]()
	lists := [][]int{{1, 2}, {3, 4}, {5}}
	result := monoid.Reduce(m, lists)
	fmt.Println(result)
	// Output: [1 2 3 4 5]
}

func ExampleMapMergeMonoid() {
	m := monoid.NewMapMergeMonoid[string, int]()
	maps := []map[string]int{
		{"a": 1, "b": 2},
		{"b": 3, "c": 4},
	}
	result := monoid.Reduce(m, maps)
	fmt.Println(result["a"], result["b"], result["c"])
	// Output: 1 3 4
}

func ExampleEndoMonoid() {
	m := monoid.NewEndoMonoid[int]()

	double := func(x int) int { return x * 2 }
	addTen := func(x int) int { return x + 10 }

	// Compose: (double ∘ addTen)(5) = double(addTen(5)) = double(15) = 30
	composed := m.Combine(double, addTen)
	fmt.Println(composed(5))
	// Output: 30
}

func ExampleOptionFirstMonoid() {
	m := monoid.NewOptionFirstMonoid[string]()

	opts := []monoid.Option[string]{
		monoid.None[string](),
		monoid.Some("first"),
		monoid.Some("second"),
	}

	result := monoid.Reduce(m, opts)
	fmt.Println(result.GetOrElse("default"))
	// Output: first
}
