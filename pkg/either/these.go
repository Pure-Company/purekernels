package either

// These represents a value that can be Left, Right, or Both
type These[L, R any] struct {
	left    L
	right   R
	variant uint8 // 0=Left, 1=Right, 2=Both
}

// ThisOnly creates a These with only left value
func ThisOnly[L, R any](left L) These[L, R] {
	return These[L, R]{left: left, variant: 0}
}

// ThatOnly creates a These with only right value
func ThatOnly[L, R any](right R) These[L, R] {
	return These[L, R]{right: right, variant: 1}
}

// Both creates a These with both values
func Both[L, R any](left L, right R) These[L, R] {
	return These[L, R]{left: left, right: right, variant: 2}
}

// IsThis returns true if this is left only
func (t These[L, R]) IsThis() bool {
	return t.variant == 0
}

// IsThat returns true if this is right only
func (t These[L, R]) IsThat() bool {
	return t.variant == 1
}

// IsBoth returns true if this has both values
func (t These[L, R]) IsBoth() bool {
	return t.variant == 2
}

// Fold applies appropriate function based on variant
func (t These[L, R]) Fold(
	onThis func(L) any,
	onThat func(R) any,
	onBoth func(L, R) any,
) any {
	switch t.variant {
	case 0:
		return onThis(t.left)
	case 1:
		return onThat(t.right)
	default:
		return onBoth(t.left, t.right)
	}
}

// FromEither converts Either to These
func FromEitherToThese[L, R any](e Either[L, R]) These[L, R] {
	if e.IsLeft() {
		return ThisOnly[L, R](e.GetLeft())
	}
	return ThatOnly[L, R](e.GetRight())
}

// ToEither converts These to Either (drops left if both)
func (t These[L, R]) ToEither() Either[L, R] {
	if t.variant == 1 || t.variant == 2 {
		return Right[L, R](t.right)
	}
	return Left[L, R](t.left)
}

// MapThis transforms the left value
func (t These[L, R]) MapThis(f func(L) L) These[L, R] {
	if t.variant == 0 || t.variant == 2 {
		return These[L, R]{left: f(t.left), right: t.right, variant: t.variant}
	}
	return t
}

// MapThat transforms the right value
func (t These[L, R]) MapThat(f func(R) R) These[L, R] {
	if t.variant == 1 || t.variant == 2 {
		return These[L, R]{left: t.left, right: f(t.right), variant: t.variant}
	}
	return t
}
