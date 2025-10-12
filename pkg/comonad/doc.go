// pkg/comonad/doc.go
// Package comonad provides comonadic structures - the categorical duals of monads.
//
// # Comonads vs Monads
//
// Monads are about putting values INTO context:
//   - Pure: A -> M[A]
//   - FlatMap: M[A] -> (A -> M[B]) -> M[B]
//
// Comonads are about extracting values FROM context:
//   - Extract: W[A] -> A
//   - Extend: W[A] -> (W[A] -> B) -> W[B]
//
// Note: Due to Go's generic limitations with recursive type instantiation,
// we provide Extract and Extend but not Duplicate (which would require
// types like W[W[A]]). Extend is the more fundamental operation anyway.
//
// # Core Types
//
// Store[S, A] - Dual of State
//   - Models a position in a space with a lookup function
//   - Use for: spreadsheet cells, game boards, zipper structures
//
// Env[E, A] - Dual of Reader
//   - Models a value paired with environment
//   - Use for: configuration with computed values, annotated data
//
// Traced[M, A] - Dual of Writer
//   - Models position-dependent computation over a monoid
//   - Use for: memoization, demand-driven computation
//
// # Comonad Laws
//
// All comonads must satisfy:
//   - Extract-Extend: Extract(Extend(w, f)) ≡ f(w)
//   - Extend-Extract: Extend(w, Extract) ≡ w
//   - Extend-Extend: Extend(Extend(w, f), g) ≡ Extend(w, func(w2) { g(Extend(w2, f)) })
//
// # Example: Store (2D Grid Navigation)
//
//	type Grid = Store[Pos, Cell]
//
//	grid := NewStore(
//	    func(pos Pos) Cell { return lookupCell(pos) },
//	    Pos{x: 0, y: 0},
//	)
//
//	// Extract current cell
//	current := grid.Extract()
//
//	// Extend to compute based on neighborhood
//	evolved := grid.Extend(func(g Store[Pos, Cell]) Cell {
//	    neighbors := g.Experiment(getNeighbors)
//	    return computeNextState(neighbors)
//	})
//
// # Example: Env (Configuration with Value)
//
//	type Config = Env[Settings, Result]
//
//	env := NewEnv(settings, computedResult)
//
//	// Extract the result
//	result := env.Extract()
//
//	// Extend with access to settings
//	transformed := env.Extend(func(e Env[Settings, Result]) Result {
//	    return processWithSettings(e.Ask(), e.Extract())
//	})
package comonad
