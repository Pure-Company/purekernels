// pkg/comonad/env.go
package comonad

// Env is the dual of Reader - a value with context
// This represents "computation with environment"
// Laws: Same comonad laws as Store
type Env[E, A any] struct {
	env   E
	value A
}

// NewEnv creates an Env with environment and value
func NewEnv[E, A any](env E, value A) Env[E, A] {
	return Env[E, A]{
		env:   env,
		value: value,
	}
}

// Extract gets the value (comonad extract)
func (e Env[E, A]) Extract() A {
	return e.value
}

// Extend applies a function that can see both env and value
func (e Env[E, A]) Extend(f func(Env[E, A]) A) Env[E, A] {
	return Env[E, A]{
		env:   e.env,
		value: f(e),
	}
}

// Map transforms the value
func (e Env[E, A]) Map(f func(A) A) Env[E, A] {
	return Env[E, A]{
		env:   e.env,
		value: f(e.value),
	}
}

// Ask returns the environment
func (e Env[E, A]) Ask() E {
	return e.env
}

// Local modifies the environment
func (e Env[E, A]) Local(f func(E) E) Env[E, A] {
	return Env[E, A]{
		env:   f(e.env),
		value: e.value,
	}
}
