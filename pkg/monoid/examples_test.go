package monoid_test

import (
	"fmt"
	"math"

	"github.com/Pure-Company/purekernels/pkg/monoid"
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
func ExampleMaxMonoid() {
	m := monoid.NewMaxMonoid[int](math.MinInt)
	result := monoid.Reduce(m, []int{5, 2, 9, 1, 7})
	fmt.Println(result)
	// Output: 9
}

func ExampleMinMonoid() {
	m := monoid.NewMinMonoid[int](math.MaxInt)
	result := monoid.Reduce(m, []int{5, 2, 9, 1, 7})
	fmt.Println(result)
	// Output: 1
}

func ExampleAvgMonoid() {
	m := monoid.NewAvgMonoid()

	// Create averages from individual measurements
	avgs := []monoid.Avg{
		monoid.FromValue(100.0),
		monoid.FromValue(200.0),
		monoid.FromValue(300.0),
	}

	result := monoid.Reduce(m, avgs)
	fmt.Printf("%.1f\n", result.Value())
	// Output: 200.0
}

func ExampleAvgMonoid_parallelAggregation() {
	m := monoid.NewAvgMonoid()

	// Simulate parallel computation of averages
	batch1 := monoid.FromValues(10, 20, 30) // avg = 20
	batch2 := monoid.FromValues(40, 50, 60) // avg = 50

	combined := m.Combine(batch1, batch2)
	fmt.Printf("Combined average: %.1f\n", combined.Value())
	// Output: Combined average: 35.0
}

func ExampleDual() {
	concat := monoid.NewStringConcatMonoid()
	dual := monoid.NewDual[string](concat)

	// Normal concatenation: left-to-right
	normal := monoid.Reduce(concat, []string{"A", "B", "C"})
	fmt.Println("Normal:", normal)

	// Dual concatenation: right-to-left
	reversed := monoid.Reduce(dual, []string{"A", "B", "C"})
	fmt.Println("Dual:", reversed)

	// Output:
	// Normal: ABC
	// Dual: CBA
}

func ExampleDual_listAppend() {
	listMonoid := monoid.NewListMonoid[int]()
	dualList := monoid.NewDual[[]int](listMonoid)

	lists := [][]int{{1, 2}, {3, 4}, {5, 6}}

	normal := monoid.Reduce(listMonoid, lists)
	reversed := monoid.Reduce(dualList, lists)

	fmt.Println("Normal:", normal)
	fmt.Println("Dual:", reversed)

	// Output:
	// Normal: [1 2 3 4 5 6]
	// Dual: [5 6 3 4 1 2]
}
