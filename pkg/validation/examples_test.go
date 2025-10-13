// pkg/validation/examples_test.go
package validation_test

import (
	"fmt"
	"strings"

	"github.com/vinodhalaharvi/purekernels/pkg/either"
	"github.com/vinodhalaharvi/purekernels/pkg/validation"
)

func ExampleBasicValidation() {
	// Single field validation
	result := validation.NotEmpty("username")("alice")
	fmt.Println("Valid:", result.IsValid())

	result2 := validation.NotEmpty("username")("")
	fmt.Println("Invalid:", result2.IsInvalid())
	if result2.IsInvalid() {
		fmt.Println(result2.GetErrors()[0])
	}

	// Output:
	// Valid: true
	// Invalid: true
	// username: cannot be empty
}

func ExampleAllCombinator() {
	// Combine multiple validators - accumulates ALL errors
	validateUsername := validation.All(
		validation.NotEmpty("username"),
		validation.MinLength("username", 3),
		validation.MaxLength("username", 20),
	)

	result := validateUsername("ab")
	if result.IsInvalid() {
		errors := result.GetErrors()
		for _, err := range errors {
			fmt.Println(err)
		}
	}

	// Output:
	// username: must be at least 3 characters
}

func ExampleAp3_ValidateStruct() {
	type User struct {
		Username string
		Email    string
		Age      int
	}

	// Define validators for each field
	validateUsername := validation.All(
		validation.NotEmpty("username"),
		validation.MinLength("username", 3),
		validation.MaxLength("username", 20),
	)

	validateEmail := validation.All(
		validation.NotEmpty("email"),
		validation.Email("email"),
	)

	validateAge := validation.Between("age", 18, 120)

	// Validate struct using Ap3 - accumulates ALL field errors
	validateUser := func(username, email string, age int) either.Validation[validation.Errors, User] {
		m := validation.NewErrorsMonoid()
		return either.Ap3(
			m,
			func(u string, e string, a int) User {
				return User{Username: u, Email: e, Age: a}
			},
			validateUsername(username),
			validateEmail(email),
			validateAge(age),
		)
	}

	// Test with invalid data
	result := validateUser("ab", "notanemail", 15)
	if result.IsInvalid() {
		errors := result.GetErrors()
		fmt.Printf("Found %d errors:\n", len(errors))
		for _, err := range errors {
			fmt.Println(err)
		}
	}

	// Output:
	// Found 3 errors:
	// username: must be at least 3 characters
	// email: must be a valid email address
	// age: must be at least 18
}

func ExampleTraverseValidation() {
	// Validate a list of emails
	emails := []string{
		"alice@example.com",
		"invalid-email",
		"bob@example.com",
		"another-bad",
	}

	validateEmail := validation.Email("email")

	m := validation.NewErrorsMonoid()
	result := either.TraverseValidation(
		m,
		func(email string) either.Validation[validation.Errors, string] {
			return validateEmail(email)
		},
		emails,
	)

	if result.IsInvalid() {
		errors := result.GetErrors()
		fmt.Printf("Invalid emails: %d\n", len(errors))
		for _, err := range errors {
			fmt.Println(err)
		}
	}

	// Output:
	// Invalid emails: 2
	// email: must be a valid email address
	// email: must be a valid email address
}

func ExampleOptional() {
	// Optional middle name
	validateMiddleName := validation.Optional(
		validation.MaxLength("middleName", 50),
	)

	// Empty is valid
	result1 := validateMiddleName("")
	fmt.Println("Empty valid:", result1.IsValid())

	// Too long is invalid
	result2 := validateMiddleName(strings.Repeat("x", 60))
	fmt.Println("Too long valid:", result2.IsValid())

	// Output:
	// Empty valid: true
	// Too long valid: false
}

func ExampleEnsure() {
	// Custom validation predicate
	validatePassword := validation.All(
		validation.MinLength("password", 8),
		validation.Ensure("password",
			func(s string) bool {
				// Must contain at least one digit
				for _, ch := range s {
					if ch >= '0' && ch <= '9' {
						return true
					}
				}
				return false
			},
			"must contain at least one digit",
		),
	)

	result := validatePassword("abcdefgh")
	if result.IsInvalid() {
		fmt.Println(result.GetErrors()[0])
	}

	// Output:
	// password: must contain at least one digit
}

func ExampleBetween_AccumulatesErrors() {
	// Between accumulates both min and max violations
	result := validation.Between("age", 18, 65)(-5)

	if result.IsInvalid() {
		errors := result.GetErrors()
		for _, err := range errors {
			fmt.Println(err)
		}
	}

	// Output:
	// age: must be at least 18
}
