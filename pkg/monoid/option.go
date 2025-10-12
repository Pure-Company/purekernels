package monoid

// Option represents an optional value (Maybe/Option type)
type Option[T any] struct {
	Value T
	Valid bool
}

// Some creates an Option with a value
func Some[T any](value T) Option[T] {
	return Option[T]{Value: value, Valid: true}
}

// None creates an empty Option
func None[T any]() Option[T] {
	return Option[T]{Valid: false}
}

type OptionFirstMonoid[T any] struct{}

// NewOptionFirstMonoid creates a first-wins option monoid
func NewOptionFirstMonoid[T any]() OptionFirstMonoid[T] {
	return OptionFirstMonoid[T]{}
}

// Empty returns None
func (OptionFirstMonoid[T]) Empty() Option[T] {
	return None[T]()
}

// Combine returns the first valid option
func (OptionFirstMonoid[T]) Combine(a, b Option[T]) Option[T] {
	if a.Valid {
		return a
	}
	return b
}

// OptionLastMonoid takes the last valid option
type OptionLastMonoid[T any] struct{}

// NewOptionLastMonoid creates a last-wins option monoid
func NewOptionLastMonoid[T any]() OptionLastMonoid[T] {
	return OptionLastMonoid[T]{}
}

// Empty returns None
func (OptionLastMonoid[T]) Empty() Option[T] {
	return None[T]()
}

// Combine returns the last valid option
func (OptionLastMonoid[T]) Combine(a, b Option[T]) Option[T] {
	if b.Valid {
		return b
	}
	return a
}

// Get returns the value and validity
func (o Option[T]) Get() (T, bool) {
	return o.Value, o.Valid
}

// GetOrElse returns the value or a default
func (o Option[T]) GetOrElse(defaultValue T) T {
	if o.Valid {
		return o.Value
	}
	return defaultValue
}

// FlatMapOption allows type transformation A -> B
func FlatMapOption[A, B any](
	opt Option[A],
	f func(A) Option[B],
) Option[B] {
	if !opt.Valid {
		return None[B]()
	}
	return f(opt.Value)
}

// FlatMap Keep existing FlatMapOption in functor/map.go or add method for convenience
func (o Option[T]) FlatMap(f func(T) Option[T]) Option[T] {
	if !o.Valid {
		return None[T]()
	}
	return f(o.Value)
}
