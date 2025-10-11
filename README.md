---

# 🔮 PureKernels

[![Go Version](https://img.shields.io/badge/Go-1.25.0+-00ADD8?style=flat\&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-FBML-purple.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

**Composable, concurrent, and mathematically correct Go.**

PureKernels provides concrete, law-abiding functional abstractions for Go — **no reflection, no interfaces, no magic** — just pure categorical algebra mapped to Go’s type system.

---

## 🧠 Philosophy

> “Make illegal states unrepresentable, make correct code inevitable.”

PureKernels replaces imperative patterns with **mathematically lawful, composable building blocks**:

* **Composition over inheritance** — small types, big systems
* **Concrete over abstract** — no type erasure, ever
* **Immutable by default** — computations never mutate state
* **Deterministic algebra** — structure and behavior are law-driven
* **Concurrency as a first-class applicative**

---

## 🧩 Core Concepts

### 1. Monoid — The Foundation of Composition

A **Monoid** defines *how values combine* and *what their identity element is*.

```go
type Monoid[T any] interface {
    Empty() T
    Combine(a, b T) T
}
```

Think: **“How do I add two things together safely?”**

Example — Set union:

```go
set1 := monoid.FromSlice([]string{"a", "b"})
set2 := monoid.FromSlice([]string{"b", "c"})

result := set1.Combine(set2)
fmt.Println(result.ToSlice()) // [a b c]
```

**Laws:**

```go
Combine(x, Empty()) == x
Combine(Empty(), x) == x
Combine(Combine(x, y), z) == Combine(x, Combine(y, z))
```

✅ Used as **input strategy**, not as a **return type** — meaning no type erasure.
Different algebraic worlds, same interface.

---

### 2. Fold — Declarative Loop Replacement

**Fold** eliminates manual loops, turning them into composable transformations.

```go
sum := fold.FoldLeft(
	func(acc, n int) int { return acc + n },
	0,
	[]int{1, 2, 3, 4},
)

evens := fold.Filter(func(n int) bool { return n%2 == 0 }, []int{1,2,3,4})
doubled := fold.Map(func(n int) int { return n*2 }, evens)
```

Fold gives you a lawful basis for map/filter/reduce, implemented with pure functions.

---

### 3. Compose — The Function Pipeline Algebra

Compose any number of functions into declarative pipelines:

```go
pipeline := compose.Pipe(
	func(x int) int { return x + 1 },
	func(x int) int { return x * 2 },
	func(x int) int { return x - 3 },
)

fmt.Println(pipeline(5)) // ((5 + 1) * 2) - 3 = 9
```

**Pure function algebra:**

* `Compose(g, f)` = `g(f(x))`
* `Pipe(f, g, h)` = `h(g(f(x)))`
* `Curry`, `Uncurry`, `Flip`, and `Const` for functional convenience

---

### 4. Functor + Applicative — Concrete and Lawful

Instead of using interfaces (`Applicative[T]`), PureKernels uses **concrete applicative structs**:

#### Concurrent Applicative

Runs computations in parallel and combines results using a `Monoid[T]`.

```go
setMonoid := monoid.NewSetMonoid[string]()

c1 := functor.NewConcurrent(setMonoid, func() monoid.SetMonoid[string] {
    return monoid.FromSlice([]string{"a"})
})

c2 := functor.NewConcurrent(setMonoid, func() monoid.SetMonoid[string] {
    return monoid.FromSlice([]string{"b"})
})

combined := c1.Apply(c2).Value()
fmt.Println(combined.ToSlice()) // [a b]
```

**Key properties:**

* No type erasure
* Parallelism built-in (`goroutines`)
* Fully lawful (Identity, Composition, Homomorphism)

**Laws:**

```
Pure(empty).Apply(v) ≡ v
u.Apply(v).Apply(w) ≡ u.Apply(v.Apply(w))
Pure(x).Apply(Pure(y)) ≡ Pure(monoid.Combine(x, y))
```

#### Sequential Applicative

Runs computations one-by-one, left-to-right (deterministic and ordered).

```go
s1 := functor.NewSequential(setMonoid, func() monoid.SetMonoid[string] {
    return monoid.FromSlice([]string{"x"})
})

s2 := functor.NewSequential(setMonoid, func() monoid.SetMonoid[string] {
    return monoid.FromSlice([]string{"y"})
})

result := s1.Apply(s2).Value()
fmt.Println(result.ToSlice()) // [x y]
```

Both types (`Concurrent[T]` and `Sequential[T]`) are **lawful applicatives** —
the only difference is execution model.

---

### 5. Traversals — Declarative Parallelism

Map a function over items and combine results declaratively:

```go
items := []int{1, 2, 3, 4}
sumMonoid := IntSumMonoid{}

result := functor.Traverse(
	sumMonoid,
	func(x int) functor.Concurrent[int] {
		return functor.NewConcurrent(sumMonoid, func() int {
			return x * x // parallel square computation
		})
	},
	items,
).Value()

fmt.Println(result) // 30 (1² + 2² + 3² + 4²)
```

Or with a worker pool:

```go
result := functor.TraverseConcurrent(
	sumMonoid,
	func(x int) int { return x * 2 },
	items,
	4, // 4 workers
).Value()
```

---

## ⚙️ Design Rules

| Concept         | Go Implementation           | Lawfulness                          |
| --------------- | --------------------------- | ----------------------------------- |
| **Monoid**      | Interface (input strategy)  | ✅ Associative + Identity            |
| **Functor**     | Concrete struct + `Map()`   | ✅ Functor law                       |
| **Applicative** | Concrete struct + `Apply()` | ✅ Applicative laws                  |
| **Concurrent**  | Parallel computation        | ✅ Associative if monoid commutative |
| **Sequential**  | Ordered computation         | ✅ Deterministic sequencing          |
| **Fold**        | Declarative iteration       | ✅ Functional totality               |
| **Compose**     | Pure function composition   | ✅ Category composition law          |

---

## 🔬 Why No Interface for Applicatives?

Earlier versions used:

```go
type Applicative[T any] interface {
	Apply(other Applicative[T]) Applicative[T]
}
```

❌ **Problem:** returns interface → loses concrete type → can’t chain

✅ **Now:** use concrete types like `Concurrent[T]` and `Sequential[T]`

* No type erasure
* Concrete, law-abiding composition
* Compiler knows full type all the way down
* Simpler and more Go-idiomatic

---

## 🚀 Example — AST Dependency Extraction

```go
func ExtractDeps(file *ast.File) Dependencies {
	setM := monoid.NewSetMonoid[string]()
	depsM := DepsMonoid{}

	return functor.Traverse(
		depsM,
		func(d ast.Decl) functor.Concurrent[Dependencies] {
			return functor.NewConcurrent(depsM, func() Dependencies {
				return ExtractDecl(d, setM)
			})
		},
		file.Decls,
	).Value()
}
```

Each piece (`Monoid`, `Concurrent`, and `Fold`) works in perfect harmony —
**immutable**, **parallel**, **composable**, and **lawful**.

---

## 🧬 Folder Overview

```
purekernels/
├── pkg/
│   ├── compose/     # Function composition (Pipe, Compose, Curry)
│   ├── fold/        # Declarative folds, maps, filters
│   ├── monoid/      # Algebraic combination strategies
│   └── functor/     # Concrete applicatives (Concurrent, Sequential)
│
├── LICENSE
└── README.md
```

---

## 🧠 License

**The MIT + Functorial Brain-Melt License (FBML)**
This software is MIT-licensed — but be warned:
*reading the source may cause categorical enlightenment or mild brain melt.*

See [LICENSE](./LICENSE) for details.

---
