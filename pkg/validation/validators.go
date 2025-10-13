// pkg/validation/validators.go
package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/vinodhalaharvi/purekernels/pkg/either"
)

// Error represents a validation error
type Error struct {
	Field   string
	Message string
}

func (e Error) String() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Errors is a list of validation errors
type Errors []Error

func (errs Errors) Strings() []string {
	result := make([]string, len(errs))
	for i, err := range errs {
		result[i] = err.String()
	}
	return result
}

// ErrorsMonoid is a monoid for combining validation errors
type ErrorsMonoid struct{}

func NewErrorsMonoid() ErrorsMonoid {
	return ErrorsMonoid{}
}

func (ErrorsMonoid) Empty() Errors {
	return Errors{}
}

func (ErrorsMonoid) Combine(a, b Errors) Errors {
	result := make(Errors, 0, len(a)+len(b))
	result = append(result, a...)
	result = append(result, b...)
	return result
}

// The global instance
var errorsMonoid = NewErrorsMonoid()

// Validator validates a value and returns Validation
type Validator[T any] func(T) either.Validation[Errors, T]

// Valid creates a valid result
func Valid[T any](value T) either.Validation[Errors, T] {
	return either.Valid[Errors, T](errorsMonoid, value)
}

// Invalid creates an invalid result with a single error
func Invalid[T any](field, message string) either.Validation[Errors, T] {
	return either.Invalid[Errors, T](
		errorsMonoid,
		Errors{{Field: field, Message: message}},
	)
}

// String Validators

// NotEmpty validates string is not empty
func NotEmpty(field string) Validator[string] {
	return func(s string) either.Validation[Errors, string] {
		if strings.TrimSpace(s) == "" {
			return Invalid[string](field, "cannot be empty")
		}
		return Valid(s)
	}
}

// MinLength validates minimum string length
func MinLength(field string, min int) Validator[string] {
	return func(s string) either.Validation[Errors, string] {
		if len(s) < min {
			return Invalid[string](field, fmt.Sprintf("must be at least %d characters", min))
		}
		return Valid(s)
	}
}

// MaxLength validates maximum string length
func MaxLength(field string, max int) Validator[string] {
	return func(s string) either.Validation[Errors, string] {
		if len(s) > max {
			return Invalid[string](field, fmt.Sprintf("must be at most %d characters", max))
		}
		return Valid(s)
	}
}

// Matches validates string matches regex
func Matches(field string, pattern *regexp.Regexp, message string) Validator[string] {
	return func(s string) either.Validation[Errors, string] {
		if !pattern.MatchString(s) {
			return Invalid[string](field, message)
		}
		return Valid(s)
	}
}

// Email validates email format
func Email(field string) Validator[string] {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return Matches(field, emailRegex, "must be a valid email address")
}

// Numeric Validators

// Min validates minimum value
func Min[T interface{ ~int | ~int64 | ~float64 }](field string, min T) Validator[T] {
	return func(val T) either.Validation[Errors, T] {
		if val < min {
			return Invalid[T](field, fmt.Sprintf("must be at least %v", min))
		}
		return Valid(val)
	}
}

// Max validates maximum value
func Max[T interface{ ~int | ~int64 | ~float64 }](field string, max T) Validator[T] {
	return func(val T) either.Validation[Errors, T] {
		if val > max {
			return Invalid[T](field, fmt.Sprintf("must be at most %v", max))
		}
		return Valid(val)
	}
}

// Between validates value is in range (accumulates both errors if outside range)
func Between[T interface{ ~int | ~int64 | ~float64 }](field string, min, max T) Validator[T] {
	return func(val T) either.Validation[Errors, T] {
		errors := Errors{}

		if val < min {
			errors = append(errors, Error{
				Field:   field,
				Message: fmt.Sprintf("must be at least %v", min),
			})
		}
		if val > max {
			errors = append(errors, Error{
				Field:   field,
				Message: fmt.Sprintf("must be at most %v", max),
			})
		}

		if len(errors) > 0 {
			return either.Invalid[Errors, T](errorsMonoid, errors)
		}
		return Valid(val)
	}
}

// Optional makes a validator optional (empty/zero value is valid)
func Optional[T comparable](validator Validator[T]) Validator[T] {
	var zero T
	return func(val T) either.Validation[Errors, T] {
		if val == zero {
			return Valid(val)
		}
		return validator(val)
	}
}

// Ensure adds a custom predicate check
func Ensure[T any](field string, pred func(T) bool, message string) Validator[T] {
	return func(val T) either.Validation[Errors, T] {
		if !pred(val) {
			return Invalid[T](field, message)
		}
		return Valid(val)
	}
}

// All runs multiple validators on same value (accumulates all errors)
func All[T any](validators ...Validator[T]) Validator[T] {
	return func(val T) either.Validation[Errors, T] {
		if len(validators) == 0 {
			return Valid(val)
		}

		// Manually accumulate errors from all validators
		allErrors := Errors{}
		hasError := false

		for _, validator := range validators {
			result := validator(val)
			if result.IsInvalid() {
				allErrors = errorsMonoid.Combine(allErrors, result.GetErrors())
				hasError = true
			}
		}

		if hasError {
			return either.Invalid[Errors, T](errorsMonoid, allErrors)
		}
		return Valid(val)
	}
}
