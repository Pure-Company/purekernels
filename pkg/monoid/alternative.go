package monoid

// OrElse for Option - returns first Some, or None
func (o Option[T]) OrElse(alternative Option[T]) Option[T] {
	if o.Valid {
		return o
	}
	return alternative
}

// OrElseLazy for Option - lazily computes alternative
func (o Option[T]) OrElseLazy(alternative func() Option[T]) Option[T] {
	if o.Valid {
		return o
	}
	return alternative()
}

// GetOrElseF returns value or computes default
func (o Option[T]) GetOrElseF(f func() T) T {
	if o.Valid {
		return o.Value
	}
	return f()
}
