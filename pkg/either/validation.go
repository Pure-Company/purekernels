package either

import "github.com/vinodhalaharvi/purekernels/pkg/monoid"

// Validation is an applicative functor for accumulating errors
// Unlike Either, it does NOT short-circuit - it accumulates all errors
// This makes it NOT a monad (no lawful FlatMap), but a perfect applicative
//
// Laws (Applicative):
//   - Identity: Pure(id).Apply(v) ≡ v
//   - Composition: u.Apply(v).Apply(w) ≡ u.Apply(v.Apply(w))
//   - Homomorphism: Pure(x).Apply(Pure(y)) ≡ Pure(f(x,y))
type Validation[E, A any] struct {
	errors  E
	value   A
	isValid bool
	monoid  monoid.Monoid[E]
}

// Valid creates a valid Validation with a value
func Valid[E, A any](m monoid.Monoid[E], value A) Validation[E, A] {
	return Validation[E, A]{
		value:   value,
		isValid: true,
		monoid:  m,
	}
}

// Invalid creates an invalid Validation with errors
func Invalid[E, A any](m monoid.Monoid[E], errors E) Validation[E, A] {
	return Validation[E, A]{
		errors:  errors,
		isValid: false,
		monoid:  m,
	}
}

// FromEither converts Either to Validation
func FromEither[E, A any](m monoid.Monoid[E], e Either[E, A]) Validation[E, A] {
	if e.isRight {
		return Valid[E, A](m, e.right)
	}
	return Invalid[E, A](m, e.left)
}

// IsValid returns true if validation succeeded
func (v Validation[E, A]) IsValid() bool {
	return v.isValid
}

// IsInvalid returns true if validation failed
func (v Validation[E, A]) IsInvalid() bool {
	return !v.isValid
}

// Map transforms the value if valid (functor)
func (v Validation[E, A]) Map(f func(A) A) Validation[E, A] {
	if !v.isValid {
		return Validation[E, A]{errors: v.errors, isValid: false, monoid: v.monoid}
	}
	return Valid[E, A](v.monoid, f(v.value))
}

// MapErrors transforms the errors if invalid
func (v Validation[E, A]) MapErrors(f func(E) E) Validation[E, A] {
	if v.isValid {
		return Validation[E, A]{value: v.value, isValid: true, monoid: v.monoid}
	}
	return Invalid[E, A](v.monoid, f(v.errors))
}

// Apply combines two validations (applicative) - ACCUMULATES errors
// This is the key difference from Either: errors accumulate instead of short-circuiting
func (v Validation[E, A]) Apply(other Validation[E, A], combine func(A, A) A) Validation[E, A] {
	switch {
	case v.isValid && other.isValid:
		// Both valid: combine values
		return Valid[E, A](v.monoid, combine(v.value, other.value))
	case !v.isValid && !other.isValid:
		// Both invalid: ACCUMULATE errors using monoid
		return Invalid[E, A](v.monoid, v.monoid.Combine(v.errors, other.errors))
	case !v.isValid:
		// First invalid
		return Invalid[E, A](v.monoid, v.errors)
	default:
		// Second invalid
		return Invalid[E, A](v.monoid, other.errors)
	}
}

// Get returns value, errors, and validity
func (v Validation[E, A]) Get() (A, E, bool) {
	return v.value, v.errors, v.isValid
}

// GetValueOr returns the value or a default
func (v Validation[E, A]) GetValueOr(defaultValue A) A {
	if v.isValid {
		return v.value
	}
	return defaultValue
}

// ToEither converts Validation to Either
func (v Validation[E, A]) ToEither() Either[E, A] {
	if v.isValid {
		return Right[E, A](v.value)
	}
	return Left[E, A](v.errors)
}

// Fold applies one of two functions
func (v Validation[E, A]) Fold(onInvalid func(E) any, onValid func(A) any) any {
	if v.isValid {
		return onValid(v.value)
	}
	return onInvalid(v.errors)
}

// Ensure adds a validation condition
func (v Validation[E, A]) Ensure(pred func(A) bool, error E) Validation[E, A] {
	if !v.isValid {
		return v
	}
	if pred(v.value) {
		return v
	}
	return Invalid[E, A](v.monoid, error)
}

// Ap2 applies a binary function to two validations
func Ap2[E, A, B, C any](
	m monoid.Monoid[E],
	f func(A, B) C,
	va Validation[E, A],
	vb Validation[E, B],
) Validation[E, C] {
	switch {
	case va.isValid && vb.isValid:
		return Valid[E, C](m, f(va.value, vb.value))
	case !va.isValid && !vb.isValid:
		return Invalid[E, C](m, m.Combine(va.errors, vb.errors))
	case !va.isValid:
		return Invalid[E, C](m, va.errors)
	default:
		return Invalid[E, C](m, vb.errors)
	}
}

// Ap3 applies a ternary function to three validations
func Ap3[E, A, B, C, D any](
	m monoid.Monoid[E],
	f func(A, B, C) D,
	va Validation[E, A],
	vb Validation[E, B],
	vc Validation[E, C],
) Validation[E, D] {
	errors := m.Empty()
	hasError := false

	if !va.isValid {
		errors = m.Combine(errors, va.errors)
		hasError = true
	}
	if !vb.isValid {
		errors = m.Combine(errors, vb.errors)
		hasError = true
	}
	if !vc.isValid {
		errors = m.Combine(errors, vc.errors)
		hasError = true
	}

	if hasError {
		return Invalid[E, D](m, errors)
	}
	return Valid[E, D](m, f(va.value, vb.value, vc.value))
}

// pkg/either/validation.go

// ValidateAll combines multiple validations, accumulating ALL errors
func ValidateAll[E, A, B any](
	m monoid.Monoid[E],
	validations []Validation[E, A],
	combine func([]A) B,
) Validation[E, B] {
	if len(validations) == 0 {
		return Invalid[E, B](m, m.Empty())
	}

	values := make([]A, 0, len(validations))
	allErrors := m.Empty()
	hasError := false

	for _, v := range validations {
		if v.isValid {
			values = append(values, v.value)
		} else {
			allErrors = m.Combine(allErrors, v.errors)
			hasError = true
		}
	}

	if hasError {
		return Invalid[E, B](m, allErrors)
	}
	return Valid[E, B](m, combine(values))
}

// SequenceValidation converts []Validation[E, A] to Validation[E, []A]
// Accumulates ALL errors if any validation fails
func SequenceValidation[E, A any](
	m monoid.Monoid[E],
	validations []Validation[E, A],
) Validation[E, []A] {
	return ValidateAll(m, validations, func(values []A) []A { return values })
}

// TraverseValidation maps and validates, accumulating errors
func TraverseValidation[E, A, B any](
	m monoid.Monoid[E],
	f func(A) Validation[E, B],
	items []A,
) Validation[E, []B] {
	validations := make([]Validation[E, B], len(items))
	for i, item := range items {
		validations[i] = f(item)
	}
	return SequenceValidation(m, validations)
}

// pkg/either/validation.go

// GetValue returns the value (panics if invalid)
func (v Validation[E, A]) GetValue() A {
	if !v.isValid {
		panic("GetValue called on Invalid")
	}
	return v.value
}

// GetErrors returns the errors (panics if valid)
func (v Validation[E, A]) GetErrors() E {
	if v.isValid {
		panic("GetErrors called on Valid")
	}
	return v.errors
}
