# 📚 PureKernels - Comprehensive README.md

```markdown
# 🔮 PureKernels

[![Go Version](https://img.shields.io/badge/Go-1.25.0+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

**Pure categorical abstractions for composable, declarative Go code.**

PureKernels is a foundational library providing core functional programming abstractions built on category theory. It replaces imperative loops and mutable state with composable, type-safe operations using Monoids, Folds, and Applicative Functors.

## 🎯 Philosophy

> "Make illegal states unrepresentable, make correct code inevitable."

PureKernels follows these principles:

- **Composition over Inheritance** - Build complex operations from simple, reusable pieces
- **Declarative over Imperative** - Express *what* you want, not *how* to get it
- **Immutability by Design** - All operations return new values, never mutate
- **Type Safety** - Leverage Go's generics to catch errors at compile time
- **Zero Dependencies** - Pure standard library, no external deps
- **Practical Category Theory** - Mathematical correctness without the jargon

## 🚀 Why PureKernels?

### The Problem

Traditional Go code often looks like this:

```go
// Imperative: Hard to reason about, easy to get wrong
func extractDependencies(nodes []Node) Dependencies {
    deps := Dependencies{
        types:     make(map[string]bool),
        functions: make(map[string]bool),
    }
    
    for _, node := range nodes {
        // Nested loops, mutable state
        for _, child := range node.Children {
            if child.Type == "function" {
                deps.functions[child.Name] = true
                // More nested logic...
            }
        }
    }
    
    return deps
}
```

**Problems:**
- Mutable state (`deps` changes throughout)
- Nested loops (hard to test, hard to reuse)
- Error-prone (forgot to initialize map? crash!)
- Not composable (can't easily combine with other operations)

### The Solution

With PureKernels:

```go
// Declarative: Clear intent, impossible to get wrong
func extractDependencies(nodes []Node) Dependencies {
    return functor.Traverse(
        extractNodeDeps,
        nodes,
    ).Value()
}

func extractNodeDeps(node Node) functor.Accumulator[Dependencies] {
    return functor.NewEmptyAccumulator(DepsMonoid).
        Apply(extractTypes(node)).
        Apply(extractFunctions(node))
}
```

**Benefits:**
- Immutable (no mutation, ever)
- Composable (combine operations easily)
- Testable (each function is pure)
- Type-safe (compiler catches errors)
- Reusable (functions work everywhere)

## 📦 Installation

```bash
go get github.com/vinodhalaharvi/purekernels
```

**Requirements:**
- Go 1.25.0 or higher

## 🧩 Core Concepts

### 1. Monoid - The Foundation

A **Monoid** is a set with an associative binary operation and an identity element.

**Think of it as:** "A way to combine things that always works the same way."

```go
import "github.com/vinodhalaharvi/purekernels/pkg/monoid"

// Example: Set union is a monoid
set1 := monoid.FromSlice([]string{"a", "b"})
set2 := monoid.FromSlice([]string{"b", "c"})
empty := monoid.NewSetMonoid[string]()

// Combine operation: union
result := set1.Combine(set2)
// result contains: {"a", "b", "c"}

// Identity element: empty set
same := set1.Combine(empty)
// same == set1 (unchanged)

// Associativity: order doesn't matter
a := set1.Combine(set2).Combine(empty)
b := set1.Combine(set2.Combine(empty))
// a == b (always true)
```

**Laws (automatically satisfied):**
```go
// Identity
monoid.Combine(x, empty) == x
monoid.Combine(empty, x) == x

// Associativity  
monoid.Combine(monoid.Combine(x, y), z) == monoid.Combine(x, monoid.Combine(y, z))
```

**Built-in Monoids:**
- `SetMonoid[T]` - Set union
- More coming (Sum, Product, String, List, etc.)

### 2. Fold - Loop Elimination

**Fold** replaces imperative loops with declarative transformations.

**Think of it as:** "Process a list by combining elements one at a time."

```go
import "github.com/vinodhalaharvi/purekernels/pkg/fold"

// Instead of:
sum := 0
for _, n := range numbers {
    sum += n
}

// Write:
sum := fold.FoldLeft(
    func(acc, n int) int { return acc + n },
    0,
    numbers,
)

// Or use built-in operations:
doubled := fold.Map(func(n int) int { return n * 2 }, numbers)
evens := fold.Filter(func(n int) bool { return n%2 == 0 }, numbers)
```

**Common Operations:**

| Operation | Imperative | Declarative |
|-----------|-----------|-------------|
| Sum | `for { sum += x }` | `fold.FoldLeft(add, 0, xs)` |
| Filter | `for { if pred(x) { ... }}` | `fold.Filter(pred, xs)` |
| Map | `for { result = append(...) }` | `fold.Map(f, xs)` |
| FlatMap | `for { for { ... }}` | `fold.FlatMap(f, xs)` |
| GroupBy | nested loops + maps | `fold.GroupBy(keyFn, xs)` |

**Key Functions:**

```go
// FoldLeft - process left to right
func FoldLeft[A, B any](f func(B, A) B, initial B, items []A) B

// FoldRight - process right to left  
func FoldRight[A, B any](f func(A, B) B, initial B, items []A) B

// Map - transform each element
func Map[A, B any](f func(A) B, items []A) []B

// Filter - keep matching elements
func Filter[A any](pred func(A) bool, items []A) []A

// FlatMap - transform and flatten
func FlatMap[A, B any](f func(A) []B, items []A) []B

// FoldList - fold with nil filtering (crucial for AST traversal)
func FoldList[A, B any](items []A, f func(A) *B) []*B

// GroupBy - group by key function
func GroupBy[A any, K comparable](keyFn func(A) K, items []A) map[K][]A

// Partition - split into two lists
func Partition[A any](pred func(A) bool, items []A) ([]A, []A)

// Any - check if any element matches
func Any[A any](pred func(A) bool, items []A) bool

// All - check if all elements match
func All[A any](pred func(A) bool, items []A) bool
```

### 3. Compose - Function Pipelines

**Compose** builds complex functions from simple ones.

**Think of it as:** "Chain operations together like UNIX pipes."

```go
import "github.com/vinodhalaharvi/purekernels/pkg/compose"

// Mathematical composition: (g ∘ f)(x) = g(f(x))
double := func(n int) int { return n * 2 }
addTen := func(n int) int { return n + 10 }

f := compose.Compose(addTen, double)
result := f(5) // double(5) = 10, then addTen(10) = 20

// Pipeline (more readable for multiple operations)
process := compose.Pipe(
    double,
    addTen,
    func(n int) int { return n * n },
)
result := process(5) // 5 -> 10 -> 20 -> 400
```

**Type-Safe Pipelines:**

```go
// Pipe2: A -> B -> C
parse := func(s string) int { /* ... */ }
validate := func(n int) bool { /* ... */ }

validator := compose.Pipe2(parse, validate)
valid := validator("42") // string -> int -> bool

// Pipe3: A -> B -> C -> D
transform := compose.Pipe3(
    parseJSON,      // string -> map
    extractField,   // map -> value
    formatOutput,   // value -> string
)
```

**Utility Functions:**

```go
// Identity - returns input unchanged
compose.Identity(x) == x

// Const - always returns same value
alwaysFive := compose.Const[int, int](5)
alwaysFive(100) == 5

// Flip - swap argument order
divide := func(a, b int) int { return a / b }
divideFlipped := compose.Flip(divide)
divide(10, 2) == 5
divideFlipped(2, 10) == 5

// Curry - convert to curried form
add := func(a, b int) int { return a + b }
curriedAdd := compose.Curry(add)
addFive := curriedAdd(5)
addFive(3) == 8
```

### 4. Applicative Functor - The Star of the Show ⭐

**Applicative** is the most powerful abstraction for combining independent computations.

**Think of it as:** "Run multiple computations and combine their results automatically."

#### Why Applicative?

When you have multiple independent operations whose results need to be combined:

```go
// Without Applicative (imperative, fragile)
func analyze(code string) Analysis {
    result := Analysis{}
    
    // Three independent operations
    types := extractTypes(code)
    result.types = append(result.types, types...)
    
    funcs := extractFunctions(code)
    result.functions = append(result.functions, funcs...)
    
    imports := extractImports(code)
    result.imports = append(result.imports, imports...)
    
    return result
}

// With Applicative (declarative, bulletproof)
func analyze(code string) Analysis {
    return functor.CombineAll(
        extractTypes(code),
        extractFunctions(code),
        extractImports(code),
    ).Value()
}
```

#### Core API

```go
import "github.com/vinodhalaharvi/purekernels/pkg/functor"

// Create an accumulator with a monoid
acc := functor.NewEmptyAccumulator(myMonoid)

// Add values
acc = acc.Add(value1)

// Combine with other accumulators  
combined := acc.Apply(otherAcc)

// Extract final result
result := acc.Value()
```

#### Real-World Example: Dependency Extraction

```go
// Define your domain types
type Dependencies struct {
    Types      monoid.SetMonoid[string]
    Functions  monoid.SetMonoid[string]
    Imports    monoid.SetMonoid[string]
}

// Implement Monoid for Dependencies
type DepsMonoid struct{}

func (DepsMonoid) Empty() Dependencies {
    return Dependencies{
        Types:     monoid.NewSetMonoid[string](),
        Functions: monoid.NewSetMonoid[string](),
        Imports:   monoid.NewSetMonoid[string](),
    }
}

func (DepsMonoid) Combine(a, b Dependencies) Dependencies {
    return Dependencies{
        Types:     a.Types.Combine(b.Types),
        Functions: a.Functions.Combine(b.Functions),
        Imports:   a.Imports.Combine(b.Imports),
    }
}

// Now extract dependencies declaratively!
func ExtractDeps(file *ast.File) Dependencies {
    empty := functor.NewEmptyAccumulator(DepsMonoid{})
    
    return functor.Traverse(
        extractDeclDeps,
        file.Decls,
    ).Value()
}

func extractDeclDeps(decl ast.Decl) functor.Accumulator[Dependencies] {
    empty := functor.NewEmptyAccumulator(DepsMonoid{})
    
    switch d := decl.(type) {
    case *ast.FuncDecl:
        return empty.
            Apply(extractFuncTypes(d)).
            Apply(extractFuncCalls(d))
    
    case *ast.GenDecl:
        return extractGenDeclDeps(d)
    
    default:
        return empty
    }
}
```

#### Applicative Operators

```go
// CombineAll - combine multiple accumulators
result := functor.CombineAll(acc1, acc2, acc3, acc4)

// Traverse - map and combine in one go
// Replaces: for item in items { acc = acc.Combine(f(item)) }
acc := functor.Traverse(processNode, nodes)

// Fold - accumulate with a function
acc := functor.Fold(empty, extractInfo, items)

// Conditional application
acc := acc.When(condition, value)        // add if true
acc := acc.Unless(condition, value)      // add if false
acc := functor.ApplyIf(cond, base, opt)  // conditional combine
```

## 🎨 Design Patterns

### Pattern 1: Replace Nested Loops

**Before:**
```go
func findMatches(items []Item) []Match {
    var results []Match
    for _, item := range items {
        for _, child := range item.Children {
            if child.IsMatch() {
                results = append(results, Match{item, child})
            }
        }
    }
    return results
}
```

**After:**
```go
func findMatches(items []Item) []Match {
    return fold.FlatMap(
        func(item Item) []Match {
            return fold.Map(
                func(child Child) Match { return Match{item, child} },
                fold.Filter(Child.IsMatch, item.Children),
            )
        },
        items,
    )
}
```

### Pattern 2: Accumulate with Type Safety

**Before:**
```go
func analyzeCode(files []File) Stats {
    stats := Stats{
        lines: make(map[string]int),
        funcs: make(map[string][]string),
    }
    
    for _, file := range files {
        // Easy to forget to initialize maps
        // Easy to have nil pointer errors
        for _, fn := range file.Functions {
            stats.funcs[file.Name] = append(stats.funcs[file.Name], fn.Name)
        }
    }
    return stats
}
```

**After:**
```go
func analyzeCode(files []File) Stats {
    return functor.Traverse(
        analyzeFile,
        files,
    ).Value()
}

// Each function is independently testable
func analyzeFile(file File) functor.Accumulator[Stats] {
    return functor.NewEmptyAccumulator(StatsMonoid).
        Apply(countLines(file)).
        Apply(extractFunctions(file))
}
```

### Pattern 3: Conditional Accumulation

**Before:**
```go
func extractInfo(node Node) Info {
    info := Info{}
    
    if node.HasType {
        info.types = append(info.types, node.Type)
    }
    
    if node.IsExported {
        info.exported = append(info.exported, node.Name)
    }
    
    return info
}
```

**After:**
```go
func extractInfo(node Node) Info {
    return functor.NewEmptyAccumulator(InfoMonoid).
        When(node.HasType, infoFromType(node.Type)).
        When(node.IsExported, infoFromExport(node.Name)).
        Value()
}
```

### Pattern 4: Pipeline Transformations

**Before:**
```go
func processData(raw string) Result {
    parsed := parse(raw)
    validated := validate(parsed)
    normalized := normalize(validated)
    enriched := enrich(normalized)
    return enriched
}
```

**After:**
```go
func processData(raw string) Result {
    return compose.Pipe3(
        parse,
        validate,
        compose.Pipe2(normalize, enrich),
    )(raw)
}

// Or more readable:
var processData = compose.Pipe3(
    parse,
    validate,
    compose.Pipe2(normalize, enrich),
)
```


---

### 🧠 The MIT + Functorial Brain-Melt License (FBML)
This project is licensed under the **MIT + Functorial Brain-Melt License (FBML)**.  
See the [LICENSE](./LICENSE) file for full terms, side-effects, and possible category-theoretic hallucinations.

---

# purekernels
