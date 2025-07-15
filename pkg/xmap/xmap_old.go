//go:build !go1.21
// +build !go1.21

package xmap

import (
	"github.com/dashjay/xiter/pkg/union"
	"github.com/dashjay/xiter/pkg/xiter"
)

func Clone[M ~map[K]V, K comparable, V any](m M) M {
	return xiter.ToMap(xiter.FromMapKeyAndValues(m))
}

func Equal[M1, M2 ~map[K]V, K, V comparable](m1 M1, m2 M2) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v1 := range m1 {
		if v2, ok := m2[k]; !ok || v1 != v2 {
			return false
		}
	}
	return true
}

func EqualFunc[M1 ~map[K]V1, M2 ~map[K]V2, K comparable, V1, V2 any](m1 M1, m2 M2, eq func(V1, V2) bool) bool {
	if len(m1) != len(m2) {
		return false
	}
	for k, v1 := range m1 {
		if v2, ok := m2[k]; !ok || !eq(v1, v2) {
			return false
		}
	}
	return true
}

func Copy[M1 ~map[K]V, M2 ~map[K]V, K comparable, V any](dst M1, src M2) {
	for k, v := range src {
		dst[k] = v
	}
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
