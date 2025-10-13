package bridges

import (
	"context"

	"github.com/vinodhalaharvi/purekernels/pkg/effect"
	"github.com/vinodhalaharvi/purekernels/pkg/monoid"
)

// FromContext extracts a typed value from context
func FromContext[T any](ctx context.Context, key any) monoid.Option[T] {
	val := ctx.Value(key)
	if val == nil {
		return monoid.None[T]()
	}
	if typed, ok := val.(T); ok {
		return monoid.Some(typed)
	}
	return monoid.None[T]()
}

// ToContext adds a value to context
func ToContext[T any](ctx context.Context, key any, value T) context.Context {
	return context.WithValue(ctx, key, value)
}

// ReaderFromContext creates a Reader from context keys
func ReaderFromContext[R, A any](
	extract func(context.Context) monoid.Option[A],
) effect.Reader[context.Context, monoid.Option[A]] {
	return func(ctx context.Context) monoid.Option[A] {
		return extract(ctx)
	}
}

// WithCancel runs a computation with cancellable context
func WithCancel[A any](
	ctx context.Context,
	f func(context.Context) A,
) (A, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	return f(ctx), cancel
}
