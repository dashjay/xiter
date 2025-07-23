//go:build go1.21
// +build go1.21

package xmap

import (
	"maps"

	"github.com/dashjay/xiter/pkg/union"
	"github.com/dashjay/xiter/pkg/xiter"
)

func Clone[M ~map[K]V, K comparable, V any](m M) M {
	return maps.Clone(m)
}

func Equal[M1, M2 ~map[K]V, K, V comparable](m1 M1, m2 M2) bool {
	return maps.Equal(m1, m2)
}

func EqualFunc[M1 ~map[K]V1, M2 ~map[K]V2, K comparable, V1, V2 any](m1 M1, m2 M2, eq func(V1, V2) bool) bool {
	return maps.EqualFunc(m1, m2, eq)
}

func Copy[M1 ~map[K]V, M2 ~map[K]V, K comparable, V any](dst M1, src M2) {
	maps.Copy(dst, src)
}

func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	return xiter.ToSlice(xiter.FromMapKeys(m))
}

func Values[M ~map[K]V, K comparable, V any](m M) []V {
	return xiter.ToSlice(xiter.FromMapValues(m))
}

func ToUnionSlice[M ~map[K]V, K comparable, V any](m M) []union.U2[K, V] {
	return xiter.ToSlice(xiter.Seq2ToSeqUnion(xiter.FromMapKeyAndValues(m)))
}
