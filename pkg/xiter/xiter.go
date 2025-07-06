//go:build go1.23
// +build go1.23

package xiter

import (
	"github.com/dashjay/xiter/pkg/internal/constraints"
	"github.com/dashjay/xiter/pkg/optional"
	"iter"
	"maps"
	"math/rand/v2"
	"strings"
)

// Seq is a sequence of elements provided by an iterator-like function.
// We made an Alias Seq to iter.Seq for providing a compatible interface in lower go versions.
type Seq[V any] iter.Seq[V]

// Seq2 is a sequence of key/value pair provided by an iterator-like function.
// We made an Alias Seq2 to iter.Seq2 for providing a compatible interface in lower go versions.
type Seq2[K, V any] iter.Seq2[K, V]

// ToSlice returns the elements in seq as a slice.
func ToSlice[T any](seq Seq[T]) (out []T) {
	for v := range seq {
		out = append(out, v)
	}
	return out
}

func ToSliceSeq2Key[K, V any](seq Seq2[K, V]) (out []K) {
	for k := range seq {
		out = append(out, k)
	}
	return
}

func ToSliceSeq2Value[K, V any](seq Seq2[K, V]) (out []V) {
	for _, v := range seq {
		out = append(out, v)
	}
	return
}

func ToMap[K comparable, V any](seq Seq2[K, V]) (out map[K]V) {
	return maps.Collect(iter.Seq2[K, V](seq))
}

func FromMapKeys[K comparable, V any](m map[K]V) Seq[K] {
	return Seq[K](maps.Keys(m))
}

func FromMapValues[K comparable, V any](m map[K]V) Seq[V] {
	return Seq[V](maps.Values(m))
}

func FromMapKeyAndValues[K comparable, V any](m map[K]V) Seq2[K, V] {
	return Seq2[K, V](maps.All(m))
}

func Pull[V any](seq Seq[V]) (next func() (V, bool), stop func()) {
	return iter.Pull(iter.Seq[V](seq))
}

func Pull2[K, V any](seq Seq2[K, V]) (next func() (K, V, bool), stop func()) {
	return iter.Pull2(iter.Seq2[K, V](seq))
}

// AllFromSeq return true if all elements from seq satisfy the condition evaluated by f.
func AllFromSeq[T any](seq Seq[T], f func(T) bool) bool {
	for t := range seq {
		if !f(t) {
			return false
		}
	}
	return true
}

// AnyFromSeq return true if any elements from seq satisfy the condition evaluated by f.
func AnyFromSeq[T any](seq Seq[T], f func(T) bool) bool {
	for t := range seq {
		if f(t) {
			return true
		}
	}
	return false
}

// AvgFromSeq return the average value of all elements from seq.
func AvgFromSeq[T constraints.Number](seq Seq[T]) float64 {
	var sum T
	count := 0
	for t := range seq {
		sum += t
		count++
	}
	if count == 0 {
		return 0
	}
	return float64(sum) / float64(count)
}

// AvgByFromSeq return the average value of all elements from seq, evaluated by f.
func AvgByFromSeq[V any, T constraints.Number](seq Seq[V], f func(V) T) float64 {
	var sum T
	count := 0
	for v := range seq {
		sum += f(v)
		count++
	}
	if count == 0 {
		return 0
	}
	return float64(sum) / float64(count)
}

// Contains return true if v is in seq.
func Contains[T comparable](seq Seq[T], in T) bool {
	for v := range seq {
		if in == v {
			return true
		}
	}
	return false
}

// ContainsBy return true if any element from seq satisfies the condition evaluated by f.
func ContainsBy[T any](seq Seq[T], f func(T) bool) bool {
	for v := range seq {
		if f(v) {
			return true
		}
	}
	return false
}

// ContainsAny return true if any element from seq is in vs.
func ContainsAny[T comparable](seq Seq[T], in []T) bool {
	if len(in) == 0 {
		return false
	}
	m := make(map[T]struct{}, len(in))
	for _, v := range in {
		m[v] = struct{}{}
	}
	for v := range seq {
		if _, exists := m[v]; exists {
			return true
		}
	}
	return false
}

// ContainsAll return true if all elements from seq is in vs.
func ContainsAll[T comparable](seq Seq[T], in []T) bool {
	if len(in) == 0 {
		return true
	}
	m := make(map[T]struct{}, len(in))
	for _, v := range in {
		m[v] = struct{}{}
	}
	for v := range seq {
		if _, exists := m[v]; exists {
			delete(m, v)
			if len(m) == 0 {
				return true
			}
		}
	}
	return len(m) == 0
}

// Count return the number of elements in seq.
func Count[T any](seq Seq[T]) int {
	var count int
	for range seq {
		count++
	}
	return count
}

// Find return the first element from seq that satisfies the condition evaluated by f with a boolean representing whether it exists.
func Find[T any](seq Seq[T], f func(T) bool) (val T, found bool) {
	for v := range seq {
		if f(v) {
			val = v
			found = true
			return
		}
	}
	return
}

// FindO return the first element from seq that satisfies the condition evaluated by f.
func FindO[T any](seq Seq[T], f func(T) bool) optional.O[T] {
	for v := range seq {
		if f(v) {
			return optional.FromValue(v)
		}
	}
	return optional.Empty[T]()
}

// ForEach execute f for each element in seq.
func ForEach[T any](seq Seq[T], f func(T) bool) {
	for v := range seq {
		if !f(v) {
			break
		}
	}
}

// ForEachIdx execute f for each element in seq with its index.
func ForEachIdx[T any](seq Seq[T], f func(idx int, v T) bool) {
	idx := 0
	for v := range seq {
		if !f(idx, v) {
			break
		}
		idx++
	}
}

// HeadO return the first element from seq.
func HeadO[T any](seq Seq[T]) optional.O[T] {
	for v := range seq {
		return optional.FromValue(v)
	}
	return optional.Empty[T]()
}

// Head return the first element from seq with a boolean representing whether it is at least one element in seq.
func Head[T any](seq Seq[T]) (v T, hasOne bool) {
	for t := range seq {
		v = t
		hasOne = true
		return
	}
	return
}

// Join return the concatenation of all elements in seq with sep.
func Join[T ~string](seq Seq[T], sep T) T {
	elems := make([]string, 0, 10)
	for v := range seq {
		elems = append(elems, string(v))
	}
	return T(strings.Join(elems, string(sep)))
}

// Max returns the maximum element in seq.
func Max[T constraints.Ordered](seq Seq[T]) (r optional.O[T]) {
	first := true
	var _max T
	for v := range seq {
		if first {
			_max = v
			first = false
		} else if _max < v {
			_max = v
		}
	}
	if first {
		return
	}
	return optional.FromValue(_max)
}

// MaxBy return the maximum element in seq, evaluated by f.
func MaxBy[T constraints.Ordered](seq Seq[T], less func(T, T) bool) (r optional.O[T]) {
	first := true
	var _max T
	for v := range seq {
		if first {
			_max = v
			first = false
		} else if less(_max, v) {
			_max = v
		}
	}
	if first {
		return
	}
	return optional.FromValue(_max)
}

// Min return the minimum element in seq.
func Min[T constraints.Ordered](seq Seq[T]) (r optional.O[T]) {
	first := true
	var _min T
	for v := range seq {
		if first {
			_min = v
			first = false
		} else if _min > v {
			_min = v
		}
	}
	if first {
		return
	}
	return optional.FromValue(_min)
}

// MinBy return the minimum element in seq, evaluated by f.
func MinBy[T constraints.Ordered](seq Seq[T], less func(T, T) bool) (r optional.O[T]) {
	first := true
	var _min T
	for v := range seq {
		if first {
			_min = v
			first = false
		} else if less(v, _min) {
			_min = v
		}
	}
	if first {
		return
	}
	return optional.FromValue(_min)
}

// PullOut pull out n elements from seq.
func PullOut[T any](seq Seq[T], n int) (out []T) {
	if n == 0 {
		return
	} else if n > 0 {
		out = make([]T, 0, n)
		for v := range seq {
			if n == 0 {
				break
			}
			out = append(out, v)
			n--
		}
		return out
	} else { // n < 0 means no limit
		out = make([]T, 0)
		for v := range seq {
			out = append(out, v)
		}
		return out
	}
}

// Skip return a seq that skip n elements from seq.
func Skip[T any](seq Seq[T], n int) Seq[T] {
	return func(yield func(T) bool) {
		for v := range seq {
			if n == 0 {
				if !yield(v) {
					break
				}
			} else {
				n--
			}
		}
	}
}

// Replace return a seq that replace from -> to
func Replace[T comparable](seq Seq[T], from, to T, n int) Seq[T] {
	return func(yield func(T) bool) {
		for v := range seq {
			// n == 0 means we have no more elements need to be replaced
			if n == 0 {
				if !yield(v) {
					break
				}
				continue
			} else if n > 0 { // we have n elements need to be replaced
				n--
			} else { // n < 0 means we need to replace all elements

			}
			if v == from {
				if !yield(to) {
					break
				}
			} else {
				if !yield(v) {
					break
				}
			}
		}
	}
}

// ReplaceAll return a seq that replace all from -> to
func ReplaceAll[T comparable](seq Seq[T], from, to T) Seq[T] {
	return Replace(seq, from, to, -1)
}

// FromSliceShuffle return a seq that shuffle the elements in the input slice.
func FromSliceShuffle[T any](in []T) Seq[T] {
	randPerm := rand.Perm(len(in))
	return func(yield func(T) bool) {
		for i := 0; i < len(randPerm); i++ {
			if !yield(in[randPerm[i]]) {
				break
			}
		}
	}
}
