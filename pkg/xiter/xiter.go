//go:build go1.23
// +build go1.23

package xiter

import (
	"iter"
	"maps"
	"math/rand/v2"
	"strings"

	"github.com/dashjay/xiter/pkg/cmp"
	"github.com/dashjay/xiter/pkg/internal/constraints"
	"github.com/dashjay/xiter/pkg/optional"
)

// Seq is a sequence of elements provided by an iterator-like function.
// We made this alias Seq to iter.Seq for providing a compatible interface in lower go versions.
type Seq[V any] iter.Seq[V]

// Seq2 is a sequence of key/value pair provided by an iterator-like function.
// We made this alias Seq2 to iter.Seq2 for providing a compatible interface in lower go versions.
type Seq2[K, V any] iter.Seq2[K, V]

// Concat returns a Seq over the concatenation of the sequences.
// It combines multiple Seqs into a single Seq by iterating each Seq one by one
// in order.
//
// Example:
//
//	seq1 := xiter.FromSlice([]int{1, 2})
//	seq2 := xiter.FromSlice([]int{3, 4})
//	seq3 := xiter.FromSlice([]int{5, 6})
//	combined := xiter.Concat(seq1, seq2, seq3)
//	fmt.Println(xiter.ToSlice(combined))
//	// output:
//	// [1 2 3 4 5 6]
func Concat[V any](seqs ...Seq[V]) Seq[V] {
	return func(yield func(V) bool) {
		for _, seq := range seqs {
			for e := range seq {
				if !yield(e) {
					return
				}
			}
		}
	}
}

// Concat2 returns an Seq2 over the concatenation of the given Seq2s.
// Like Concat but run with Seq2
func Concat2[K, V any](seqs ...Seq2[K, V]) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, seq := range seqs {
			for k, v := range seq {
				if !yield(k, v) {
					return
				}
			}
		}
	}
}

// Equal returns whether the two sequences are equal.
// It compares elements from both Seq in parallel. If the Seq have different lengths
// or if any corresponding elements are not equal, it returns false.
//
// Example:
//
//	seq1 := xiter.FromSlice([]int{1, 2, 3})
//	seq2 := xiter.FromSlice([]int{1, 2, 3})
//	seq3 := xiter.FromSlice([]int{1, 2, 4})
//	fmt.Println(xiter.Equal(seq1, seq2))
//	fmt.Println(xiter.Equal(seq1, seq3))
//	// output:
//	// true
//	// false
func Equal[V comparable](x, y Seq[V]) bool {
	for z := range Zip(x, y) {
		if z.Ok1 != z.Ok2 || z.V1 != z.V2 {
			return false
		}
	}
	return true
}

// Equal2 returns whether the two Seq2 are equal.
// Like Equal but run with Seq2
func Equal2[K, V comparable](x, y Seq2[K, V]) bool {
	for z := range Zip2(x, y) {
		if z.Ok1 != z.Ok2 || z.K1 != z.K2 || z.V1 != z.V2 {
			return false
		}
	}
	return true
}

// EqualFunc returns whether the two sequences are equal according to the function f.
// Example:
//
//	seq1 := xiter.FromSlice([]int{6, 11, 16})
//	seq2 := xiter.FromSlice([]int{26, 36, 41})
//	seq3 := xiter.FromSlice([]int{1, 2, 4})
//	mod5Eq := func(a int, b int) bool {
//		return math.Mod(float64(a), 5) == math.Mod(float64(b), 5)
//	}
//	fmt.Println(xiter.EqualFunc(seq1, seq2, mod5Eq))
//	fmt.Println(xiter.EqualFunc(seq1, seq3, mod5Eq))
//	// output:
//	// true
//	// false
func EqualFunc[V1, V2 any](x Seq[V1], y Seq[V2], f func(V1, V2) bool) bool {
	for z := range Zip(x, y) {
		if z.Ok1 != z.Ok2 || !f(z.V1, z.V2) {
			return false
		}
	}
	return true
}

// EqualFunc2 returns whether the two sequences are equal according to the function f.
// Like EqualFunc but run with Seq2
func EqualFunc2[K1, V1, K2, V2 any](x Seq2[K1, V1], y Seq2[K2, V2], f func(K1, V1, K2, V2) bool) bool {
	for z := range Zip2(x, y) {
		if z.Ok1 != z.Ok2 || !f(z.K1, z.V1, z.K2, z.V2) {
			return false
		}
	}
	return true
}

// Filter returns a Seq over seq that only includes
// the values v for which f(v) is true.
//
// Example:
//
//	seq := FromSlice([]int{1, 2, 3, 4, 5})
//	evenNumbers := Filter(func(v int) bool { return v%2 == 0 }, seq)
//	// evenNumbers will yield: 2, 4
func Filter[V any](f func(V) bool, seq Seq[V]) Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			if f(v) && !yield(v) {
				return
			}
		}
	}
}

// Filter2 returns an Seq over seq that only includes
// the key-value pairs k, v for which f(k, v) is true.
// Like Filter but run with Seq2
func Filter2[K, V any](f func(K, V) bool, seq Seq2[K, V]) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range seq {
			if f(k, v) && !yield(k, v) {
				return
			}
		}
	}
}

// Limit returns an iterator over the first n values of seq.
// If n is less than or equal to 0, an empty sequence is returned.
//
// Example:
//
//	seq := xiter.FromSlice([]int{1, 2, 3, 4, 5})
//	limitedSeq := xiter.Limit(seq, 3)
//	fmt.Println(xiter.ToSlice(limitedSeq))
//	// output:
//	// [1 2 3]
func Limit[V any](seq Seq[V], n int) Seq[V] {
	return func(yield func(V) bool) {
		if n <= 0 {
			return
		}
		for v := range seq {
			if !yield(v) {
				return
			}
			n--
			if n <= 0 {
				break
			}
		}
	}
}

// Limit2 returns a Seq over Seq2 that stops after n key-value pairs.
// Like Limit but run with Seq2
func Limit2[K, V any](seq Seq2[K, V], n int) Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if n <= 0 {
			return
		}
		for k, v := range seq {
			if !yield(k, v) {
				return
			}
			n--
			if n <= 0 {
				break
			}
		}
	}
}

// Map returns a Seq over the results of applying f to each value in seq.
//
// Example:
//
//	seq := xiter.FromSlice([]int{1, 2, 3})
//	doubled := xiter.Map(func(v int) int { return v * 2 }, seq)
//	fmt.Println(xiter.ToSlice(doubled))
//	// output:
//	// [2 4 6]
func Map[In, Out any](f func(In) Out, seq Seq[In]) Seq[Out] {
	return func(yield func(Out) bool) {
		for in := range seq {
			if !yield(f(in)) {
				return
			}
		}
	}
}

// Map2 returns a Seq2 over the results of applying f to each key-value pair in seq.
// Like Map but run with Seq2
func Map2[KIn, VIn, KOut, VOut any](
	f func(KIn, VIn) (KOut, VOut),
	seq Seq2[KIn, VIn]) Seq2[KOut, VOut] {
	return func(yield func(KOut, VOut) bool) {
		for k, v := range seq {
			if !yield(f(k, v)) {
				return
			}
		}
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
		for v1 := range x {
			for ok2 && f(v1, v2) > 0 {
				if !yield(v2) {
					return
				}
				v2, ok2 = next()
			}
			if !yield(v1) {
				return
			}
		}
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
		for k1, v1 := range x {
			for ok2 && f(k1, k2) > 0 {
				if !yield(k2, v2) {
					return
				}
				k2, v2, ok2 = next()
			}
			if !yield(k1, v1) {
				return
			}
		}
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
// Example:
//
//	seq1To100 := xiter.FromSlice(_range(1, 101))
//	sum := xiter.Reduce(func(sum int, v int) int {
//		return sum + v
//	}, 0, seq1To100)
//	fmt.Println(sum)
//	// output:
//	// 5050
func Reduce[Sum, V any](f func(Sum, V) Sum, sum Sum, seq Seq[V]) Sum {
	for v := range seq {
		sum = f(sum, v)
	}
	return sum
}

// Reduce2 combines the values in seq using f.
// For each pair k, v in seq, it updates sum = f(sum, k, v)
// and then returns the final sum.
// For example, if iterating over seq yields (k1, v1), (k2, v2), (k3, v3)
// Reduce returns f(f(f(sum, k1, v1), k2, v2), k3, v3).
func Reduce2[Sum, K, V any](f func(Sum, K, V) Sum, sum Sum, seq Seq2[K, V]) Sum {
	for k, v := range seq {
		sum = f(sum, k, v)
	}
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
		for v1 := range x {
			if !yield(Zipped[V1, V2]{v1, true, v2, ok2}) {
				return
			}
			v2, ok2 = next()
		}
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
		for k1, v1 := range x {
			if !yield(Zipped2[K1, V1, K2, V2]{k1, v1, true, k2, v2, ok2}) {
				return
			}
			k2, v2, ok2 = next()
		}
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

// ToSlice returns the elements in seq as a slice.
func ToSlice[T any](seq Seq[T]) (out []T) {
	for v := range seq {
		out = append(out, v)
	}
	return out
}

// ToSliceSeq2Key returns the keys in seq as a slice.
//
// Example:
//
//	seq := FromMap(map[string]int{"a": 1, "b": 2})
//	keys := ToSliceSeq2Key(seq)
//	// keys will contain: []string{"a", "b"} (order may vary)
func ToSliceSeq2Key[K, V any](seq Seq2[K, V]) (out []K) {
	for k := range seq {
		out = append(out, k)
	}
	return
}

// ToSliceSeq2Value returns the values in seq as a slice.
//
// Example:
//
//	seq := FromMap(map[string]int{"a": 1, "b": 2})
//	values := ToSliceSeq2Value(seq)
//	// values will contain: []int{1, 2} (order may vary)
func ToSliceSeq2Value[K, V any](seq Seq2[K, V]) (out []V) {
	for _, v := range seq {
		out = append(out, v)
	}
	return
}

func Seq2KeyToSeq[K, V any](in Seq2[K, V]) Seq[K] {
	return func(yield func(K) bool) {
		for k := range in {
			if !yield(k) {
				break
			}
		}
	}
}

func Seq2ValueToSeq[K, V any](in Seq2[K, V]) Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range in {
			if !yield(v) {
				break
			}
		}
	}
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

// Pull wrapped iter.Pull create an iterator from seq.
// Example:
//
//	seq := xiter.FromSlice([]int{1, 2, 3})
//	next, stop := xiter.Pull(seq)
//	defer stop()
//	x, ok := next()
//	fmt.Println(x, ok)
//	// output:
//	// 1 true
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

// ToSliceN pull out n elements from seq.
func ToSliceN[T any](seq Seq[T], n int) (out []T) {
	switch {
	case n == 0:
		return
	case n > 0:
		out = make([]T, 0, n)
		for v := range seq {
			if n == 0 {
				break
			}
			out = append(out, v)
			n--
		}
		return out
	default:
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
//
// Example:
//
//	seq := FromSlice([]int{1, 2, 3, 2, 4})
//	replacedSeq := Replace(seq, 2, 99, -1) // Replace all 2s with 99
//	// replacedSeq will yield: 1, 99, 3, 99, 4
func Replace[T comparable](seq Seq[T], from, to T, n int) Seq[T] {
	return func(yield func(T) bool) {
		for v := range seq {
			if n != 0 && v == from {
				if !yield(to) {
					break
				}
				if n > 0 {
					n--
				}
			} else if !yield(v) {
				break
			}
		}
	}
}

// ReplaceAll return a seq that replace all from -> to
//
// Example:
//
//	seq := FromSlice([]int{1, 2, 3, 2, 4})
//	replacedSeq := ReplaceAll(seq, 2, 99)
//	// replacedSeq will yield: 1, 99, 3, 99, 4
func ReplaceAll[T comparable](seq Seq[T], from, to T) Seq[T] {
	return Replace(seq, from, to, -1)
}

// FromSliceShuffle return a seq that shuffle the elements in the input slice.
//
// Example:
//
//	seq := FromSlice([]int{1, 2, 3, 4, 5})
//	shuffledSeq := FromSliceShuffle(ToSlice(seq))
//	// shuffledSeq will yield a shuffled sequence of 1, 2, 3, 4, 5
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
