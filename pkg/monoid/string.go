package monoid

// StringConcatMonoid represents string concatenation with empty string as identity
// Laws:
//   - Identity: Combine(x, "") == x && Combine("", x) == x
//   - Associativity: Combine(Combine(x, y), z) == Combine(x, Combine(y, z))
type StringConcatMonoid struct{}

// NewStringConcatMonoid creates a string concatenation monoid
func NewStringConcatMonoid() StringConcatMonoid {
	return StringConcatMonoid{}
}

// Empty returns empty string (concatenation identity)
func (StringConcatMonoid) Empty() string {
	return ""
}

// Combine concatenates two strings
func (StringConcatMonoid) Combine(a, b string) string {
	return a + b
}

// StringJoinMonoid represents string concatenation with a separator
// Laws:
//   - Identity: Combine(x, "") == x (when separator is applied)
//   - Associativity: Combine(Combine(x, y), z) == Combine(x, Combine(y, z))
type StringJoinMonoid struct {
	separator string
}

// NewStringJoinMonoid creates a string join monoid with separator
func NewStringJoinMonoid(sep string) StringJoinMonoid {
	return StringJoinMonoid{separator: sep}
}

// Empty returns empty string
func (StringJoinMonoid) Empty() string {
	return ""
}

// Combine joins strings with separator (only if both are non-empty)
func (s StringJoinMonoid) Combine(a, b string) string {
	if a == "" {
		return b
	}
	if b == "" {
		return a
	}
	return a + s.separator + b
}
