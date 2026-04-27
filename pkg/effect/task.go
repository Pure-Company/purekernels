package effect

import "github.com/Pure-Company/purekernels/pkg/result"

// Task represents an asynchronous computation
type Task[A any] struct {
	run func() result.Result[A]
}

func NewTask[A any](f func() result.Result[A]) Task[A] {
	return Task[A]{run: f}
}

func TaskOf[A any](value A) Task[A] {
	return Task[A]{run: func() result.Result[A] {
		return result.Ok(value)
	}}
}

func TaskErr[A any](err error) Task[A] {
	return Task[A]{run: func() result.Result[A] {
		return result.Err[A](err)
	}}
}

// Run executes async
func (t Task[A]) Run() <-chan result.Result[A] {
	ch := make(chan result.Result[A], 1)
	go func() {
		ch <- t.run()
		close(ch)
	}()
	return ch
}

func (t Task[A]) Map(f func(A) A) Task[A] {
	return Task[A]{
		run: func() result.Result[A] {
			return t.run().Map(f)
		},
	}
}

func (t Task[A]) FlatMap(f func(A) Task[A]) Task[A] {
	return Task[A]{
		run: func() result.Result[A] {
			res := t.run()
			if res.IsErr() {
				return res
			}
			return f(res.Unwrap()).run()
		},
	}
}
