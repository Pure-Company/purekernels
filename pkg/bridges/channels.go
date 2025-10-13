package bridges

import (
	"context"
	"time"

	"github.com/vinodhalaharvi/purekernels/pkg/monoid"
	"github.com/vinodhalaharvi/purekernels/pkg/result"
)

// CollectFromChannel drains a channel into a slice
func CollectFromChannel[T any](ch <-chan T) []T {
	items := []T{}
	for item := range ch {
		items = append(items, item)
	}
	return items
}

// CollectWithTimeout drains a channel with timeout
func CollectWithTimeout[T any](
	ch <-chan T,
	timeout time.Duration,
) result.Result[[]T] {
	items := []T{}
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	for {
		select {
		case item, ok := <-ch:
			if !ok {
				return result.Ok(items)
			}
			items = append(items, item)
		case <-timer.C:
			return result.Err[[]T](context.DeadlineExceeded)
		}
	}
}

// TryReceive attempts to receive from channel without blocking
func TryReceive[T any](ch <-chan T) monoid.Option[T] {
	select {
	case val, ok := <-ch:
		if !ok {
			return monoid.None[T]()
		}
		return monoid.Some(val)
	default:
		return monoid.None[T]()
	}
}

// SendWithTimeout attempts to send with timeout
func SendWithTimeout[T any](
	ch chan<- T,
	value T,
	timeout time.Duration,
) result.Result[struct{}] {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case ch <- value:
		return result.Ok(struct{}{})
	case <-timer.C:
		return result.Err[struct{}](context.DeadlineExceeded)
	}
}

// FanOut sends a value to multiple channels
func FanOut[T any](value T, channels ...chan<- T) {
	for _, ch := range channels {
		ch <- value
	}
}

// FanIn merges multiple channels into one
func FanIn[T any](channels ...<-chan T) <-chan T {
	out := make(chan T)

	for _, ch := range channels {
		go func(c <-chan T) {
			for val := range c {
				out <- val
			}
		}(ch)
	}

	return out
}
