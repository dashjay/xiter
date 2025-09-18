//go:build !go1.23
// +build !go1.23

package xiter

import (
	"math/rand"
	"runtime"
	"strconv"
	"strings"

	"github.com/dashjay/xiter/pkg/cmp"
	"github.com/dashjay/xiter/pkg/internal/constraints"
	"github.com/dashjay/xiter/pkg/internal/utils"
	"github.com/dashjay/xiter/pkg/optional"
	"github.com/dashjay/xiter/pkg/union"
)

type Seq[V any] func(yield func(V) bool)

type Seq2[K, V any] func(yield func(K, V) bool)

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

func Map[In, Out any](f func(In) Out, seq Seq[In]) Seq[Out] {
	return func(yield func(Out) bool) {
		seq(func(in In) bool {
			return yield(f(in))
		})
	}
}

func Map2[KIn, VIn, KOut, VOut any](f func(KIn, VIn) (KOut, VOut), seq Seq2[KIn, VIn]) Seq2[KOut, VOut] {
	return func(yield func(KOut, VOut) bool) {
		seq(func(k KIn, v VIn) bool {
			return yield(f(k, v))
		})
	}
}

func Merge[V cmp.Ordered](x, y Seq[V]) Seq[V] {
	return MergeFunc(x, y, cmp.Compare[V])
}

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

func Merge2[K cmp.Ordered, V any](x, y Seq2[K, V]) Seq2[K, V] {
	return MergeFunc2(x, y, cmp.Compare[K])
}

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

func Reduce[Sum, V any](f func(Sum, V) Sum, sum Sum, seq Seq[V]) Sum {
	seq(func(v V) bool {
		sum = f(sum, v)
		return true
	})
	return sum
}

func Reduce2[Sum, K, V any](f func(Sum, K, V) Sum, sum Sum, seq Seq2[K, V]) Sum {
	seq(func(k K, v V) bool {
		sum = f(sum, k, v)
		return true
	})
	return sum
}

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

func Seq2KeyToSeq[K, V any](in Seq2[K, V]) Seq[K] {
	return func(yield func(K) bool) {
		in(func(k K, v V) bool {
			return yield(k)
		})
	}
}

func Seq2ValueToSeq[K, V any](in Seq2[K, V]) Seq[V] {
	return func(yield func(V) bool) {
		in(func(k K, v V) bool {
			return yield(v)
		})
	}
}

func ToMap[K comparable, V any](seq Seq2[K, V]) (out map[K]V) {
	out = make(map[K]V)
	seq(func(k K, v V) bool {
		out[k] = v
		return true
	})
	return out
}

func ToMapFromSeq[K comparable, V any](seq Seq[K], fn func(k K) V) (out map[K]V) {
	out = make(map[K]V)
	seq(func(k K) bool {
		out[k] = fn(k)
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

// Pull has async operation, under 1.23 we provide goroutine version which will create a lot goroutine,
// the behavior is not same as the 1.23 version
// so we strongly recommend you to use go new version like 1.23 to use this function.
// Deprecated: Upgrade to go 1.23 to use the internal implement
func Pull[V any](seq Seq[V]) (next func() (V, bool), stop func()) {
	ch := make(chan V)
	done := make(chan struct{})
	quit := make(chan struct{})
	panicked := make(chan interface{})
	var seqGoroutineID int64

	go func() {
		seqGoroutineID = getGID()
		defer func() {
			if r := recover(); r != nil {
				panicked <- r
			}
			close(ch)
		}()
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

	}()
	next = func() (v V, ok bool) {
		if getGID() == seqGoroutineID {
			panic("xiter: next called re-entrantly")
		}
		select {
		case err, goRoutinePanicked := <-panicked:
			if goRoutinePanicked {
				panic(err)
			}
			return v, ok
		case v, ok = <-ch:
			return v, ok
		case <-done:
			return v, false
		}
	}

	stop = func() {
		select {
		case <-quit:

		default:
			close(quit)
		}
	}

	return next, stop
}

// Pull2 has async operation, under 1.23 we provide goroutine version which will create a lot goroutine,
// the behavior is not same as the 1.23 version
// so we strongly recommend you to use go new version like 1.23 to use this function.
// Deprecated: Upgrade to go 1.23 to use the internal implement
func Pull2[K, V any](seq Seq2[K, V]) (next func() (K, V, bool), stop func()) {
	ch := make(chan union.U2[K, V])
	done := make(chan struct{})
	quit := make(chan struct{})
	panicked := make(chan interface{})
	var seqGoroutineID int64

	go func() {
		seqGoroutineID = getGID()
		defer func() {
			if r := recover(); r != nil {
				panicked <- r
			}
			close(ch)
		}()
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
	}()
	next = func() (k K, v V, ok bool) {
		if getGID() == seqGoroutineID {
			panic("xiter: next called re-entrantly")
		}
		select {
		case err, goRoutinePanicked := <-panicked:
			if goRoutinePanicked {
				panic(err)
			}
			return k, v, ok
		case u2, ok := <-ch:
			return u2.T1, u2.T2, ok
		case <-done:
			return k, v, false
		}
	}

	stop = func() {
		select {
		case <-quit:

		default:
			close(quit)
		}
	}

	return next, stop
}

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

func Count[T any](seq Seq[T]) int {
	var count int
	seq(func(t T) bool {
		count++
		return true
	})
	return count
}

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

func ForEach[T any](seq Seq[T], f func(T) bool) {
	seq(func(t T) bool {
		return f(t)
	})
}

func ForEachIdx[T any](seq Seq[T], f func(idx int, v T) bool) {
	i := 0
	seq(func(t T) bool {
		c := f(i, t)
		i++
		return c
	})
}

func HeadO[T any](seq Seq[T]) optional.O[T] {
	res := optional.Empty[T]()
	seq(func(t T) bool {
		res = optional.FromValue(t)
		return false
	})
	return res
}

func Head[T any](seq Seq[T]) (v T, hasOne bool) {
	seq(func(t T) bool {
		v = t
		hasOne = true
		return false
	})
	return
}

func Join[T ~string](seq Seq[T], sep T) T {

	elems := make([]string, 0, 10)
	seq(func(t T) bool {
		elems = append(elems, string(t))
		return true
	})
	return T(strings.Join(elems, string(sep)))
}

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

func MaxBy[T any](seq Seq[T], less func(T, T) bool) (r optional.O[T]) {
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

func MinBy[T any](seq Seq[T], less func(T, T) bool) (r optional.O[T]) {
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

func ToSliceN[T any](seq Seq[T], n int) (out []T) {
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
	} else {
		seq(func(t T) bool {
			out = append(out, t)
			return true
		})
	}
	return out
}

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

func ReplaceAll[T comparable](seq Seq[T], from, to T) Seq[T] {
	return Replace(seq, from, to, -1)
}

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

// getGID returns the current goroutine ID.
// This is an internal function and should not be used in production code.
// It's used here for testing purposes to detect re-entrant calls.
func getGID() int64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	s := strings.TrimPrefix(string(b), "goroutine ")
	s = s[:strings.Index(s, " ")]
	gid, _ := strconv.ParseInt(s, 10, 64)
	return gid
}

func Chunk[T any](seq Seq[T], n int) Seq[[]T] {
	return func(yield func([]T) bool) {
		tmp := make([]T, 0, n)
		seq(func(v T) bool {
			tmp = append(tmp, v)
			if len(tmp) == n {
				con := yield(tmp)
				tmp = make([]T, 0, n)
				return con
			}
			return true
		})
		if len(tmp) > 0 {
			yield(tmp)
		}
	}
}

func Seq2ToSeqUnion[K, V any](seq Seq2[K, V]) Seq[union.U2[K, V]] {
	return func(yield func(union.U2[K, V]) bool) {
		seq(func(k K, v V) bool {
			return yield(union.U2[K, V]{T1: k, T2: v})
		})
	}
}

func Sum[T constraints.Number](seq Seq[T]) T {
	var sum T
	seq(func(t T) bool {
		sum += t
		return true
	})
	return sum
}

func Index[T comparable](seq Seq[T], v T) int {
	found := false
	idx := 0
	seq(func(t T) bool {
		if t == v {
			found = true
			return false
		}
		idx++
		return true
	})
	if found {
		return idx
	} else {
		return -1
	}
}

func Uniq[T comparable](seq Seq[T]) Seq[T] {
	return func(yield func(T) bool) {
		m := make(map[T]struct{})
		seq(func(v T) bool {
			if _, ok := m[v]; !ok {
				m[v] = struct{}{}
				return yield(v)
			}
			return true
		})
	}
}

func MapToSeq2[T any, K comparable](in Seq[T], mapFn func(T) K) Seq2[K, T] {
	return func(yield func(K, T) bool) {
		in(func(v T) bool {
			k := mapFn(v)
			return yield(k, v)
		})
	}
}

func MapToSeq2Value[T any, K comparable, V any](in Seq[T], mapFn func(T) (K, V)) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		in(func(ele T) bool {
			k, v := mapFn(ele)
			return yield(k, v)
		})
	}
}

func First[T any](in Seq[T]) (T, bool) {
	var v T
	var ok = false
	in(func(t T) bool {
		v = t
		ok = true
		return false
	})
	return v, ok
}

func FirstO[T any](in Seq[T]) optional.O[T] {
	return optional.FromValue2(First(in))
}

func Last[T any](in Seq[T]) (T, bool) {
	var v T
	var ok = false
	in(func(t T) bool {
		v = t
		ok = true
		return true
	})
	return v, ok
}

func LastO[T any](in Seq[T]) optional.O[T] {
	return optional.FromValue2(Last(in))
}

func Compact[T comparable](in Seq[T]) Seq[T] {
	return func(yield func(T) bool) {
		in(func(t T) bool {
			if utils.IsZero(t) {
				return true
			}
			return yield(t)
		})
	}
}
