//go:build go1.23
// +build go1.23

package cmp

import "cmp"

type Ordered cmp.Ordered

func Less[T Ordered](x, y T) bool {
	return cmp.Less(x, y)
}

func Compare[T Ordered](x, y T) int {
	return cmp.Compare(x, y)
}

func Or[T comparable](vals ...T) T {
	return cmp.Or(vals...)
}
