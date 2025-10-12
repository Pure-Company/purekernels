// pkg/comonad/examples_test.go
package comonad_test

import (
	"fmt"
	"github.com/vinodhalaharvi/purekernels/pkg/comonad"
	"github.com/vinodhalaharvi/purekernels/pkg/monoid"
)

func ExampleStore_cellularAutomata() {
	// Conway's Game of Life cell
	type Cell bool

	// 1D automaton: compute next state based on neighbors
	rule110 := func(store comonad.Store[int, Cell]) Cell {
		left := store.Peek(store.Pos() - 1)
		center := store.Extract()
		right := store.Peek(store.Pos() + 1)

		// Rule 110: binary representation determines next state
		pattern := 0
		if left {
			pattern += 4
		}
		if center {
			pattern += 2
		}
		if right {
			pattern += 1
		}

		// Rule 110 lookup
		return pattern == 1 || pattern == 2 || pattern == 3 || pattern == 5 || pattern == 6
	}

	// Initial configuration
	grid := comonad.NewStore(
		func(pos int) Cell {
			// Start with single cell
			return pos == 0
		},
		0,
	)

	// Evolve one step
	nextGen := grid.Extend(rule110)

	fmt.Println("Current:", grid.Extract())
	fmt.Println("Next:", nextGen.Extract())

	// Output:
	// Current: true
	// Next: true
}

func ExampleEnv_configuration() {
	type Config struct {
		Multiplier int
	}

	// Computation with access to config
	env := comonad.NewEnv(Config{Multiplier: 10}, 42)

	// Extract current value
	fmt.Println("Value:", env.Extract())

	// Extend: compute using both env and value
	doubled := env.Extend(func(e comonad.Env[Config, int]) int {
		return e.Extract() * e.Ask().Multiplier
	})

	fmt.Println("Doubled:", doubled.Extract())

	// Output:
	// Value: 42
	// Doubled: 420
}

func ExampleTraced_memoization() {
	sumMonoid := monoid.NewSumMonoid[int]()

	// Traced computation: memoized addition
	traced := comonad.NewTraced(sumMonoid, func(offset int) string {
		return fmt.Sprintf("Result at %d", offset)
	})

	// Extract at identity (0)
	fmt.Println(traced.Extract())

	// Trace at different positions
	fmt.Println(traced.Trace(5))
	fmt.Println(traced.Trace(10))

	// Output:
	// Result at 0
	// Result at 5
	// Result at 10
}
