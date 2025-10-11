// Package fold provides functional loop replacements via folding operations
package fold

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
