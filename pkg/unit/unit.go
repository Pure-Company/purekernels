package unit

// Unit is the unit type, representing a value with no information.
// It's useful as a placeholder when a function needs to return something
// but has no meaningful value to return (similar to void in other languages).
//
// Laws: There is only one value of type Unit, so all Unit values are equal.
type Unit struct{}

// Void The singleton instance (optional, for convenience)
var Void = Unit{}

// String returns the string representation
func (Unit) String() string {
	return "()"
}

// MarshalJSON encodes Unit as JSON null or empty object
func (Unit) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}

// UnmarshalJSON decodes Unit from JSON
func (u *Unit) UnmarshalJSON(b []byte) error {
	return nil
}
