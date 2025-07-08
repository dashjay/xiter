//go:build !go1.23
// +build !go1.23

package xiter

import (
	"math/rand"
	"strings"

	"github.com/dashjay/xiter/pkg/cmp"
	"github.com/dashjay/xiter/pkg/internal/constraints"
	"github.com/dashjay/xiter/pkg/optional"
	"github.com/dashjay/xiter/pkg/union"
	"github.com/panjf2000/ants/v2"
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

// Concat returns an iterator over the concatenation of the sequences.
func Concat[V any](seqs ...Seq[V]) Seq[V] {
	return func(yield func(V) bool) {
		contine := false
		for _, seq := range seqs {
			seq(func(v V) bool {
				contine = yield(v)
				return contine
			})
			if !contine {
				break
			}
		}
	}
}

// Concat2 returns an iterator over the concatenation of the sequences.
func Concat2[K, V any](seqs ...Seq2[K, V]) Seq2[K, V] {
	contine := false
	return func(yield func(K, V) bool) {
		for _, seq := range seqs {
			seq(func(k K, v V) bool {
				contine = yield(k, v)
				return contine
			})
			if !contine {
				return
			}
		}
	}
}

// Equal reports whether the two sequences are equal.
func Equal[V comparable](x, y Seq[V]) bool {
	eq := true
	Zip(x, y)(func(z Zipped[V, V]) bool {
		if z.Ok1 != z.Ok2 || z.V1 != z.V2 {
			eq = false
			return false
		}
		return true
	})
	return eq
}

// Equal2 reports whether the two sequences are equal.
func Equal2[K, V comparable](x, y Seq2[K, V]) bool {
	eq := true
	Zip2(x, y)(func(z Zipped2[K, V, K, V]) bool {
		if z.Ok1 != z.Ok2 || z.K1 != z.K2 || z.V1 != z.V2 {
			eq = false
			return false
		}
		return true
	})

	return eq
}

// EqualFunc reports whether the two sequences are equal according to the function f.
func EqualFunc[V1, V2 any](x Seq[V1], y Seq[V2], f func(V1, V2) bool) bool {
	eq := true
	Zip(x, y)(func(z Zipped[V1, V2]) bool {
		if z.Ok1 != z.Ok2 || !f(z.V1, z.V2) {
			eq = false
			return false
		}
		return true
	})
	return eq
}

// EqualFunc2 reports whether the two sequences are equal according to the function f.
func EqualFunc2[K1, V1, K2, V2 any](x Seq2[K1, V1], y Seq2[K2, V2], f func(K1, V1, K2, V2) bool) bool {
	eq := true
	Zip2(x, y)(func(z Zipped2[K1, V1, K2, V2]) bool {
		if z.Ok1 != z.Ok2 || !f(z.K1, z.V1, z.K2, z.V2) {
			eq = false
			return false
		}
		return true
	})
	return eq
}

// Filter returns an iterator over seq that only includes
// the values v for which f(v) is true.
func Filter[V any](f func(V) bool, seq Seq[V]) Seq[V] {
	return func(yield func(V) bool) {
		seq(func(v V) bool {
			if f(v) && !yield(v) {
				return false
			}
			return true
		})
	}
}

// Filter2 returns an iterator over seq that only includes
// the pairs k, v for which f(k, v) is true.
func Filter2[K, V any](f func(K, V) bool, seq Seq2[K, V]) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		seq(func(k K, v V) bool {
			if f(k, v) && !yield(k, v) {
				return false
			}
			return true
		})
	}
}

// Limit returns an iterator over seq that stops after n values.
func Limit[V any](seq Seq[V], n int) Seq[V] {
	return func(yield func(V) bool) {
		if n <= 0 {
			return
		}
		seq(func(v V) bool {
			if n <= 0 {
				return false
			}
			if !yield(v) {
				return false
			}
			if n--; n <= 0 {
				return false
			}
			return true
		})
	}
}

// Limit2 returns an iterator over seq that stops after n key-value pairs.
func Limit2[K, V any](seq Seq2[K, V], n int) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if n <= 0 {
			return
		}
		seq(func(k K, v V) bool {
			if n <= 0 {
				return false
			}
			if !yield(k, v) {
				return false
			}
			if n--; n <= 0 {
				return false
			}
			return true
		})
	}
}

// Map returns an iterator over f applied to seq.
func Map[In, Out any](f func(In) Out, seq Seq[In]) Seq[Out] {
	return func(yield func(Out) bool) {
		seq(func(in In) bool {
			return yield(f(in))
		})
	}
}

// Map2 returns an iterator over f applied to seq.
func Map2[KIn, VIn, KOut, VOut any](f func(KIn, VIn) (KOut, VOut), seq Seq2[KIn, VIn]) Seq2[KOut, VOut] {
	return func(yield func(KOut, VOut) bool) {
		seq(func(k KIn, v VIn) bool {
			return yield(f(k, v))
		})
	}
}

// Merge merges two sequences of ordered values.
// Values appear in the output once for each time they appear in x
// and once for each time they appear in y.
// If the two input sequences are not ordered,
// the output sequence will not be ordered,
// but it will still contain every value from x and y exactly once.
//
// Merge is equivalent to calling MergeFunc with cmp.Compare[V]
// as the ordering function.
func Merge[V cmp.Ordered](x, y Seq[V]) Seq[V] {
	return MergeFunc(x, y, cmp.Compare[V])
}

// MergeFunc merges two sequences of values ordered by the function f.
// Values appear in the output once for each time they appear in x
// and once for each time they appear in y.
// When equal values appear in both sequences,
// the output contains the values from x before the values from y.
// If the two input sequences are not ordered by f,
// the output sequence will not be ordered by f,
// but it will still contain every value from x and y exactly once.
func MergeFunc[V any](x, y Seq[V], f func(V, V) int) Seq[V] {
	return func(yield func(V) bool) {
		next, stop := Pull(y)
		defer stop()
		v2, ok2 := next()
		x(func(v1 V) bool {
			for ok2 && f(v1, v2) > 0 {
				if !yield(v2) {
					return false
				}
				v2, ok2 = next()
			}
			if !yield(v1) {
				return false
			}
			return true
		})
		for ok2 {
			if !yield(v2) {
				return
			}
			v2, ok2 = next()
		}
	}
}

// Merge2 merges two sequences of key-value pairs ordered by their keys.
// Pairs appear in the output once for each time they appear in x
// and once for each time they appear in y.
// If the two input sequences are not ordered by their keys,
// the output sequence will not be ordered by its keys,
// but it will still contain every pair from x and y exactly once.
//
// Merge2 is equivalent to calling MergeFunc2 with cmp.Compare[K]
// as the ordering function.
func Merge2[K cmp.Ordered, V any](x, y Seq2[K, V]) Seq2[K, V] {
	return MergeFunc2(x, y, cmp.Compare[K])
}

// MergeFunc2 merges two sequences of key-value pairs ordered by the function f.
// Pairs appear in the output once for each time they appear in x
// and once for each time they appear in y.
// When pairs with equal keys appear in both sequences,
// the output contains the pairs from x before the pairs from y.
// If the two input sequences are not ordered by f,
// the output sequence will not be ordered by f,
// but it will still contain every pair from x and y exactly once.
func MergeFunc2[K, V any](x, y Seq2[K, V], f func(K, K) int) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		next, stop := Pull2(y)
		defer stop()
		k2, v2, ok2 := next()
		x(func(k1 K, v1 V) bool {
			for ok2 && f(k1, k2) > 0 {
				if !yield(k2, v2) {
					return false
				}
				k2, v2, ok2 = next()
			}
			if !yield(k1, v1) {
				return false
			}
			return true
		})
		for ok2 {
			if !yield(k2, v2) {
				return
			}
			k2, v2, ok2 = next()
		}
	}
}

// Reduce combines the values in seq using f.
// For each value v in seq, it updates sum = f(sum, v)
// and then returns the final sum.
// For example, if iterating over seq yields v1, v2, v3,
// Reduce returns f(f(f(sum, v1), v2), v3).
func Reduce[Sum, V any](f func(Sum, V) Sum, sum Sum, seq Seq[V]) Sum {
	seq(func(v V) bool {
		sum = f(sum, v)
		return true
	})
	return sum
}

// Reduce2 combines the values in seq using f.
// For each pair k, v in seq, it updates sum = f(sum, k, v)
// and then returns the final sum.
// For example, if iterating over seq yields (k1, v1), (k2, v2), (k3, v3)
// Reduce returns f(f(f(sum, k1, v1), k2, v2), k3, v3).
func Reduce2[Sum, K, V any](f func(Sum, K, V) Sum, sum Sum, seq Seq2[K, V]) Sum {
	seq(func(k K, v V) bool {
		sum = f(sum, k, v)
		return true
	})
	return sum
}

// Zip returns an iterator that iterates x and y in parallel,
// yielding Zipped values of successive elements of x and y.
// If one sequence ends before the other, the iteration continues
// with Zipped values in which either Ok1 or Ok2 is false,
// depending on which sequence ended first.
//
// Zip is a useful building block for adapters that process
// pairs of sequences. For example, Equal can be defined as:
//
//	func Equal[V comparable](x, y Seq[V]) bool {
//		for z := range Zip(x, y) {
//			if z.Ok1 != z.Ok2 || z.V1 != z.V2 {
//				return false
//			}
//		}
//		return true
//	}
func Zip[V1, V2 any](x Seq[V1], y Seq[V2]) Seq[Zipped[V1, V2]] {
	return func(yield func(z Zipped[V1, V2]) bool) {
		next, stop := Pull(y)
		defer stop()
		v2, ok2 := next()
		x(func(v1 V1) bool {
			if !yield(Zipped[V1, V2]{v1, true, v2, ok2}) {
				return false
			}
			v2, ok2 = next()
			return true
		})

		var zv1 V1
		for ok2 {
			if !yield(Zipped[V1, V2]{zv1, false, v2, ok2}) {
				return
			}
			v2, ok2 = next()
		}
	}
}

// Zip2 returns an iterator that iterates x and y in parallel,
// yielding Zipped2 values of successive elements of x and y.
// If one sequence ends before the other, the iteration continues
// with Zipped2 values in which either Ok1 or Ok2 is false,
// depending on which sequence ended first.
//
// Zip2 is a useful building block for adapters that process
// pairs of sequences. For example, Equal2 can be defined as:
//
//	func Equal2[K, V comparable](x, y Seq2[K, V]) bool {
//		for z := range Zip2(x, y) {
//			if z.Ok1 != z.Ok2 || z.K1 != z.K2 || z.V1 != z.V2 {
//				return false
//			}
//		}
//		return true
//	}
func Zip2[K1, V1, K2, V2 any](x Seq2[K1, V1], y Seq2[K2, V2]) Seq[Zipped2[K1, V1, K2, V2]] {
	return func(yield func(z Zipped2[K1, V1, K2, V2]) bool) {
		next, stop := Pull2(y)
		defer stop()
		k2, v2, ok2 := next()
		x(func(k1 K1, v1 V1) bool {
			if !yield(Zipped2[K1, V1, K2, V2]{k1, v1, true, k2, v2, ok2}) {
				return false
			}
			k2, v2, ok2 = next()
			return true
		})

		var zk1 K1
		var zv1 V1
		for ok2 {
			if !yield(Zipped2[K1, V1, K2, V2]{zk1, zv1, false, k2, v2, ok2}) {
				return
			}
			k2, v2, ok2 = next()
		}
	}
}

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
		count := n
		seq(func(v T) bool {
			if count != 0 && v == from {
				if !yield(to) {
					return false
				}
				if count > 0 {
					count--
				}
			} else {
				if !yield(v) {
					return false
				}
			}
			return true
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
