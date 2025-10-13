// Package fold provides functional loop replacements via folding operations
package fold

import (
	"github.com/vinodhalaharvi/purekernels/pkg/monoid"
)

// FoldLeft folds from left to right - replaces for loops
func FoldLeft[A, B any](f func(B, A) B, initial B, items []A) B {
	result := initial
	for _, item := range items {
		result = f(result, item)
	}
	return result
}

// FoldRight folds from right to left
func FoldRight[A, B any](f func(A, B) B, initial B, items []A) B {
	result := initial
	for i := len(items) - 1; i >= 0; i-- {
		result = f(items[i], result)
	}
	return result
}

// Filter keeps elements matching predicate - replaces filter loops
func Filter[A any](pred func(A) bool, items []A) []A {
	return FoldLeft(
		func(acc []A, item A) []A {
			if pred(item) {
				return append(acc, item)
			}
			return acc
		},
		[]A{},
		items,
	)
}

// Map transforms each element - replaces map loops
func Map[A, B any](f func(A) B, items []A) []B {
	return FoldLeft(
		func(acc []B, item A) []B {
			return append(acc, f(item))
		},
		[]B{},
		items,
	)
}

// MapWithIndex transforms each element with its index (pure!)
func MapWithIndex[A, B any](f func(int, A) B, items []A) []B {
	type IndexedAcc struct {
		Result []B
		Index  int
	}

	indexed := FoldLeft(
		func(acc IndexedAcc, item A) IndexedAcc {
			return IndexedAcc{
				Result: append(acc.Result, f(acc.Index, item)),
				Index:  acc.Index + 1,
			}
		},
		IndexedAcc{Result: []B{}, Index: 0},
		items,
	)

	return indexed.Result
}

// FlatMap transforms and flattens - replaces nested loops
func FlatMap[A, B any](f func(A) []B, items []A) []B {
	return FoldLeft(
		func(acc []B, item A) []B {
			return append(acc, f(item)...)
		},
		[]B{},
		items,
	)
}

// FoldList folds a slice, keeping only non-nil results
// This is key for AST traversal - replaces filter+map patterns
func FoldList[A, B any](items []A, f func(A) *B) []*B {
	return FoldLeft(
		func(acc []*B, item A) []*B {
			if result := f(item); result != nil {
				return append(acc, result)
			}
			return acc
		},
		[]*B{},
		items,
	)
}

// Partition splits into two lists based on predicate
func Partition[A any](pred func(A) bool, items []A) ([]A, []A) {
	trueItems := []A{}
	falseItems := []A{}

	for _, item := range items {
		if pred(item) {
			trueItems = append(trueItems, item)
		} else {
			falseItems = append(falseItems, item)
		}
	}

	return trueItems, falseItems
}

// GroupBy groups items by a key function - replaces grouping loops
func GroupBy[A any, K comparable](keyFn func(A) K, items []A) map[K][]A {
	return FoldLeft(
		func(acc map[K][]A, item A) map[K][]A {
			key := keyFn(item)
			acc[key] = append(acc[key], item)
			return acc
		},
		make(map[K][]A),
		items,
	)
}

// Any checks if any element satisfies predicate
func Any[A any](pred func(A) bool, items []A) bool {
	for _, item := range items {
		if pred(item) {
			return true
		}
	}
	return false
}

// All checks if all elements satisfy predicate
func All[A any](pred func(A) bool, items []A) bool {
	for _, item := range items {
		if !pred(item) {
			return false
		}
	}
	return true
}

// Reduce combines elements using a monoid (categorical!)
func Reduce[A any](m monoid.Monoid[A], items []A) A {
	return FoldLeft(
		func(acc A, item A) A {
			return m.Combine(acc, item)
		},
		m.Empty(),
		items,
	)
}

// FoldMap is the canonical categorical fold - all others derive from this
func FoldMap[A, M any](m monoid.Monoid[M], f func(A) M, items []A) M {
	return FoldLeft(
		func(acc M, a A) M { return m.Combine(acc, f(a)) },
		m.Empty(),
		items,
	)
}

// ScanLeft produces intermediate accumulation results
func ScanLeft[A, B any](f func(B, A) B, initial B, items []A) []B {
	acc := initial
	results := make([]B, 0, len(items)+1)
	results = append(results, acc)
	for _, item := range items {
		acc = f(acc, item)
		results = append(results, acc)
	}
	return results
}

// ScanRight produces intermediate results from the right
func ScanRight[A, B any](f func(A, B) B, initial B, items []A) []B {
	results := make([]B, len(items)+1)
	results[len(items)] = initial
	for i := len(items) - 1; i >= 0; i-- {
		results[i] = f(items[i], results[i+1])
	}
	return results
}

// Zip combines two slices pairwise
func Zip[A, B, C any](as []A, bs []B, f func(A, B) C) []C {
	n := min(len(as), len(bs))
	result := make([]C, n)
	for i := 0; i < n; i++ {
		result[i] = f(as[i], bs[i])
	}
	return result
}

// ZipWith is an alias for Zip with different argument order (more conventional)
func ZipWith[A, B, C any](f func(A, B) C, as []A, bs []B) []C {
	return Zip(as, bs, f)
}
