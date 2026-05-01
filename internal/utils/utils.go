package utils

func IsZero[T comparable](v T) bool {
	var zero T
	return v == zero
}
