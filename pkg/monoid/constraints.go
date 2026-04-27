// pkg/monoid/constraints.go
//
// Local type-constraint definitions for monoid implementations.
//
// We used to depend on golang.org/x/exp/constraints for Integer,
// Float, and Ordered. As of Go 1.21 the standard library ships
// cmp.Ordered, which covers the ordered case. The integer and
// float cases inline cleanly as small union types — they're well
// under a screen of code, never change, and remove an external
// dependency that didn't earn its keep.
//
// Why this matters: the only thing purekernels imported from
// outside the standard library was constraints. Defining them
// locally drops the entire third-party require block from go.mod.
// Builds get faster, supply-chain surface area goes to zero.

package monoid

// Integer is the union of every signed and unsigned integer type,
// matching golang.org/x/exp/constraints.Integer. The ~T form means
// "T or any named type whose underlying type is T", so e.g. a
// custom `type Count int` satisfies the constraint.
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Float is the union of float32 and float64.
type Float interface {
	~float32 | ~float64
}

// Numeric is Integer | Float — the values that support + and *.
// Used by Sum and Product monoids where either kind works.
type Numeric interface {
	Integer | Float
}
