//go:build !go1.23
// +build !go1.23

package xiter

import (
	"github.com/dashjay/xiter/pkg/internal/constraints"
	"github.com/dashjay/xiter/pkg/optional"
	"github.com/dashjay/xiter/pkg/union"
	"github.com/panjf2000/ants/v2"
	"math/rand"
	"strings"
)

var globalXiterPool *ants.Pool

func init() {
	globalXiterPool, _ = ants.NewPool(1_000_000)
}

// Seq is a sequence of elements provided by an iterator-like function.
// Before Go1.23, golang has not stabled iter package, so we had to define this type
type Seq[V any] func(yield func(V) bool)

// Seq2 is a sequence of key/value pair provided by an iterator-like function.
// Before Go1.23, golang has not stabled iter package, so we had to define this type
type Seq2[K, V any] func(yield func(K, V) bool)

// ToSlice return a slice containing all elements from seq.
func ToSlice[T any](seq Seq[T]) (out []T) {
	seq(func(t T) bool {
		out = append(out, t)
		return true
	})
	return
}

func ToSliceSeq2Key[K, V any](seq Seq2[K, V]) (out []K) {
	seq(func(k K, v V) bool {
		out = append(out, k)
		return true
	})
	return
}

func ToSliceSeq2Value[K, V any](seq Seq2[K, V]) (out []V) {
	seq(func(k K, v V) bool {
		out = append(out, v)
		return true
	})
	return
}

func ToMap[K comparable, V any](seq Seq2[K, V]) (out map[K]V) {
	out = make(map[K]V)
	seq(func(k K, v V) bool {
		out[k] = v
		return true
	})
	return out
}

func FromMapKeys[K comparable, V any](m map[K]V) Seq[K] {
	return func(yield func(K) bool) {
		for k := range m {
			if !yield(k) {
				break
			}
		}
	}
}

func FromMapValues[K comparable, V any](m map[K]V) Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range m {
			if !yield(v) {
				break
			}
		}
	}
}

func FromMapKeyAndValues[K comparable, V any](m map[K]V) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range m {
			if !yield(k, v) {
				break
			}
		}
	}
}

func Pull[V any](seq Seq[V]) (next func() (V, bool), stop func()) {
	ch := make(chan V)
	done := make(chan struct{})
	quit := make(chan struct{})
	err := globalXiterPool.Submit(
		func() {
			defer close(ch)
			seq(func(v V) bool {
				select {
				case ch <- v:
					return true
				case <-quit:
					return false
				}
			})
			select {
			case done <- struct{}{}:
			case <-quit:
			}
		},
	)
	if err != nil {
		panic(err)
	}

	next = func() (v V, ok bool) {
		select {
		case v, ok = <-ch:
			return v, ok
		case <-done:
			return v, false
		}
	}

	stop = func() {
		select {
		case <-quit:
			// Already closed
		default:
			close(quit)
		}
	}

	return next, stop
}

func Pull2[K, V any](seq Seq2[K, V]) (next func() (K, V, bool), stop func()) {
	ch := make(chan union.U2[K, V])
	done := make(chan struct{})
	quit := make(chan struct{})
	err := globalXiterPool.Submit(
		func() {
			defer close(ch)
			seq(func(k K, v V) bool {
				select {
				case ch <- union.U2[K, V]{T1: k, T2: v}:
					return true
				case <-quit:
					return false
				}
			})
			select {
			case done <- struct{}{}:
			case <-quit:
			}
		},
	)
	if err != nil {
		panic(err)
	}

	next = func() (k K, v V, ok bool) {
		select {
		case u2, ok := <-ch:
			return u2.T1, u2.T2, ok
		case <-done:
			return k, v, false
		}
	}

	stop = func() {
		select {
		case <-quit:
			// Already closed
		default:
			close(quit)
		}
	}

	return next, stop
}

// AllFromSeq return true if all elements from seq satisfy the condition evaluated by f.
func AllFromSeq[T any](seq Seq[T], f func(T) bool) bool {
	res := true
	seq(func(v T) bool {
		if !f(v) {
			res = false
			return false
		}
		return true
	})
	return res
}

// AnyFromSeq return true if any elements from seq satisfy the condition evaluated by f.
func AnyFromSeq[T any](seq Seq[T], f func(T) bool) bool {
	res := false
	seq(func(v T) bool {
		if f(v) {
			res = true
			return false
		}
		return true
	})
	return res
}

// AvgFromSeq return the average value of all elements from seq.
func AvgFromSeq[T constraints.Number](seq Seq[T]) float64 {
	var sum T
	count := 0

	seq(func(t T) bool {
		sum += t
		count++
		return true
	})
	if count == 0 {
		return 0
	}
	return float64(sum) / float64(count)
}

// AvgByFromSeq return the average value of all elements from seq, evaluated by f.
func AvgByFromSeq[V any, T constraints.Number](seq Seq[V], f func(V) T) float64 {
	var sum T
	count := 0

	seq(func(v V) bool {
		sum += f(v)
		count++
		return true
	})
	if count == 0 {
		return 0
	}
	return float64(sum) / float64(count)
}

// Contains return true if v is in seq.
func Contains[T comparable](seq Seq[T], v T) bool {
	res := false
	seq(func(t T) bool {
		if v == t {
			res = true
			return false
		}
		return true
	})
	return res
}

// ContainsBy return true if any element from seq satisfies the condition evaluated by f.
func ContainsBy[T any](seq Seq[T], f func(T) bool) bool {
	res := false
	seq(func(t T) bool {
		if f(t) {
			res = true
			return false
		}
		return true
	})
	return res
}

// ContainsAny return true if any element from seq is in vs.
func ContainsAny[T comparable](seq Seq[T], vs []T) bool {
	if len(vs) == 0 {
		return false
	}
	m := make(map[T]struct{}, len(vs))
	for _, v := range vs {
		m[v] = struct{}{}
	}
	res := false
	seq(func(t T) bool {
		if _, exists := m[t]; exists {
			res = true
			return false
		}
		return true
	})
	return res
}

// ContainsAll return true if all elements from seq is in vs.
func ContainsAll[T comparable](seq Seq[T], vs []T) bool {
	if len(vs) == 0 {
		return true
	}
	m := make(map[T]struct{}, len(vs))
	for _, v := range vs {
		m[v] = struct{}{}
	}
	seq(func(t T) bool {
		if _, exists := m[t]; exists {
			delete(m, t)
			if len(m) == 0 {
				return false
			}
		}
		return true
	})
	return len(m) == 0
}

// Count return the number of elements in seq.
func Count[T any](seq Seq[T]) int {
	var count int
	seq(func(t T) bool {
		count++
		return true
	})
	return count
}

// Find return the first element from seq that satisfies the condition evaluated by f with a boolean representing whether it exists.
func Find[T any](seq Seq[T], f func(T) bool) (val T, found bool) {
	seq(func(t T) bool {
		found = f(t)
		if found {
			val = t
			return false
		}
		return true
	})
	return
}

// FindO return the first element from seq that satisfies the condition evaluated by f.
func FindO[T any](seq Seq[T], f func(T) bool) optional.O[T] {
	var res = optional.Empty[T]()
	seq(func(t T) bool {
		if f(t) {
			res = optional.FromValue(t)
			return false
		}
		return true
	})
	return res
}

// ForEach execute f for each element in seq.
func ForEach[T any](seq Seq[T], f func(T) bool) {
	seq(func(t T) bool {
		return f(t)
	})
}

// ForEachIdx execute f for each element in seq with its index.
func ForEachIdx[T any](seq Seq[T], f func(idx int, v T) bool) {
	i := 0
	seq(func(t T) bool {
		c := f(i, t)
		i++
		return c
	})
}

// HeadO return the first element from seq.
func HeadO[T any](seq Seq[T]) optional.O[T] {
	res := optional.Empty[T]()
	seq(func(t T) bool {
		res = optional.FromValue(t)
		return false
	})
	return res
}

// Head return the first element from seq with a boolean representing whether it is at least one element in seq.
func Head[T any](seq Seq[T]) (v T, hasOne bool) {
	seq(func(t T) bool {
		v = t
		hasOne = true
		return false
	})
	return
}

// Join return the concatenation of all elements in seq with sep.
func Join[T ~string](seq Seq[T], sep T) T {
	//var out T
	//first := false
	//seq(func(t T) bool {
	//	if first {
	//		first = true
	//	} else {
	//		out += sep
	//	}
	//	out += t
	//	return true
	//})
	//return out

	elems := make([]string, 0, 10)
	seq(func(t T) bool {
		elems = append(elems, string(t))
		return true
	})
	return T(strings.Join(elems, string(sep)))
}

// Max returns the maximum element in seq.
func Max[T constraints.Ordered](seq Seq[T]) (r optional.O[T]) {
	first := true
	var _max T
	seq(func(v T) bool {
		if first {
			_max = v
			first = false
		} else if _max < v {
			_max = v
		}
		return true
	})
	if first {
		return
	}
	return optional.FromValue(_max)
}

// MaxBy return the maximum element in seq, evaluated by f.
func MaxBy[T constraints.Ordered](seq Seq[T], less func(T, T) bool) (r optional.O[T]) {
	first := true
	var _max T
	seq(func(v T) bool {
		if first {
			_max = v
			first = false
		} else if less(_max, v) {
			_max = v
		}
		return true
	})
	if first {
		return
	}
	return optional.FromValue(_max)
}

// Min return the minimum element in seq.
func Min[T constraints.Ordered](seq Seq[T]) (r optional.O[T]) {
	first := true
	var _min T
	seq(func(v T) bool {
		if first {
			_min = v
			first = false
		} else if _min > v {
			_min = v
		}
		return true
	})
	if first {
		return
	}
	return optional.FromValue(_min)
}

// MinBy return the minimum element in seq, evaluated by f.
func MinBy[T constraints.Ordered](seq Seq[T], less func(T, T) bool) (r optional.O[T]) {
	first := true
	var _min T
	seq(func(v T) bool {
		if first {
			_min = v
			first = false
		} else if less(v, _min) {
			_min = v
		}
		return true
	})
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
		seq(func(t T) bool {
			if n == 0 {
				return false
			}
			out = append(out, t)
			n--
			return true
		})
	} else { // n < 0
		seq(func(t T) bool {
			out = append(out, t)
			return true
		})
	}
	return out
}

// Skip return a seq that skip n elements from seq.
func Skip[T any](seq Seq[T], n int) Seq[T] {
	return func(yield func(T) bool) {
		seq(func(v T) bool {
			if n == 0 {
				return yield(v)
			} else {
				n--
			}
			return true
		})
	}
}

// Replace return a seq that replace from -> to
func Replace[T comparable](seq Seq[T], from, to T, n int) Seq[T] {
	return func(yield func(T) bool) {
		seq(func(v T) bool {
			// n == 0 means we have no more elements need to be replaced
			if n == 0 {
				return yield(v)
			} else if n > 0 { // we have n elements need to be replaced
				n--
			} else { // n < 0 means we need to replace all elements

			}
			if v == from {
				return yield(to)
			} else {
				return yield(v)
			}
		})
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
