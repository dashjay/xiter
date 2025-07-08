package optional

import "fmt"

type O[T any] struct {
	value T
	ok    bool
}

// FromValue creates an Optional from a value.
func FromValue[T any](v T) O[T] {
	return O[T]{value: v, ok: true}
}

// FromValue2 creates an Optional from a value and a boolean meaning whether the value is valid.
//
// HINT:
//
//	if we have a function defined as fn () (T, bool),
//	we can use FromValue2(fn()) instead of
//	if v, ok := fn(); ok { FromValue(v) } else { Empty[T]() }
func FromValue2[T any](v T, ok bool) O[T] {
	return O[T]{value: v, ok: ok}
}

// Empty creates an empty Optional with no value.
func Empty[T any]() O[T] {
	return O[T]{ok: false}
}

// Ptr returns a pointer to the value of the Optional.
// return nil if the Optional has no value.
func (o O[T]) Ptr() *T {
	if o.ok {
		return &o.value
	}
	return nil
}

// Ok returns whether the Optional has a valid value.
func (o O[T]) Ok() bool {
	return o.ok
}

// Must directly return the value of the Optional.
//
// ❌WARNING: Panic if the Optional has no value.
func (o O[T]) Must() T {
	if o.ok {
		return o.value
	}
	panic(fmt.Sprintf("Optional.O[%T] has no valid value", o.value))
}

// ValueOr returns the value of the Optional if it has a valid value, otherwise returns the given default value.
func (o O[T]) ValueOr(dft T) T {
	if o.ok {
		return o.value
	}
	return dft
}

// ValueOrZero returns the value of the Optional if it has a valid value, otherwise returns the zero value of T.
func (o O[T]) ValueOrZero() T {
	if o.ok {
		return o.value
	}
	var empty T
	return empty
}
