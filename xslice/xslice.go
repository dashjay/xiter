package xslice

import (
	"math/rand"

	"github.com/dashjay/xiter/internal/constraints"
	"github.com/dashjay/xiter/internal/xassert"
	"github.com/dashjay/xiter/optional"
	"github.com/dashjay/xiter/xiter"
)

// All returns true if all elements in the slice satisfy the condition provided by f.
// return false if any element in the slice does not satisfy the condition provided by f.
//
// EXAMPLE:
//
//	xslice.All([]int{1, 2, 3}, func(x int) bool { return x > 0 }) 👉 true
//	xslice.All([]int{-1, 1, 2, 3}, func(x int) bool { return x > 0 }) 👉 false
func All[T any](in []T, f func(T) bool) bool {
	return xiter.AllFromSeq(xiter.FromSlice(in), f)
}

// Any returns true if any element in the slice satisfy the condition provided by f.
// return false if none of  element in the slice satisfy the condition provided by f.
//
// EXAMPLE:
//
//	xslice.Any([]int{0, 1, 2, 3}, func(x int) bool { return x == 0 }) 👉 true
//	xslice.Any([]int{0, 1, 2, 3}, func(x int) bool { return x == -1 }) 👉 false
func Any[T any](in []T, f func(T) bool) bool {
	return xiter.AnyFromSeq(xiter.FromSlice(in), f)
}

// Avg returns the average value of the items in slice (float64).
//
// EXAMPLE:
//
//	xslice.Avg([]int{1, 2, 3}) 👉 float(2)
//	xslice.Avg([]int{}) 👉 float(0)
func Avg[T constraints.Number](in []T) float64 {
	if len(in) == 0 {
		return 0
	}
	var sum T
	for _, v := range in {
		sum += v
	}
	return float64(sum) / float64(len(in))
}

// AvgN returns the average value of the items
//
// EXAMPLE:
//
//	xslice.AvgN(1, 2, 3) 👉 float(2)
//	xslice.AvgN() 👉 float(0)
func AvgN[T constraints.Number](inputs ...T) float64 {
	return Avg(inputs)
}

// AvgBy returns the averaged of each item's value evaluated by f.
//
// EXAMPLE:
//
//	xslice.AvgBy([]string{"1", "2", "3"}, func(x string) int {
//		i, _ := strconv.Atoi(x)
//		return i
//	}) 👉 float(2)
func AvgBy[V any, T constraints.Number](in []V, f func(V) T) float64 {
	if len(in) == 0 {
		return 0
	}
	var sum T
	for _, v := range in {
		sum += f(v)
	}
	return float64(sum) / float64(len(in))
}

// Contains returns true if the slice contains the value v.
//
// EXAMPLE:
//
//	xslice.Contains([]int{1, 2, 3}, 1) 👉 true
//	xslice.Contains([]int{-1, 2, 3}, 1) 👉 false
func Contains[T comparable](in []T, v T) bool {
	return xiter.Contains(xiter.FromSlice(in), v)
}

// ContainsBy returns true if the slice contains the value v evaluated by f.
//
// EXAMPLE:
//
//	xslice.ContainsBy([]string{"1", "2", "3"}, func(x string) bool {
//		i, _ := strconv.Atoi(x)
//		return i == 1
//	}) 👉 true
//
//	xslice.ContainsBy([]string{"1", "2", "3"}, func(x string) bool {
//		i, _ := strconv.Atoi(x)
//		return i == -1
//	}) 👉 false
func ContainsBy[T any](in []T, f func(T) bool) bool {
	return xiter.ContainsBy(xiter.FromSlice(in), f)
}

// ContainsAny returns true if the slice contains any value in v.
//
// EXAMPLE:
//
//	xslice.ContainsAny([]string{"1", "2", "3"}, []string{"1", "99", "1000"}) 👉 true
//	xslice.ContainsAny([]string{"1", "2", "3"}, []string{"-1"}) 👉 false
//	xslice.ContainsAny([]string{"1", "2", "3"}, []string{}) 👉 false
func ContainsAny[T comparable](in []T, v []T) bool {
	return xiter.ContainsAny(xiter.FromSlice(in), v)
}

// ContainsAll returns true if the slice contains all values in v.
//
// EXAMPLE:
//
//	xslice.ContainsAll([]string{"1", "2", "3"}, []string{"1", "2", "3"})  👉 true
//	xslice.ContainsAll([]string{"1", "2", "3"}, []string{"1", "99", "1000"}) 👉 false
//	xslice.ContainsAll([]string{"1", "2", "3"}, []string{}) 👉 true
func ContainsAll[T comparable](in []T, v []T) bool {
	return xiter.ContainsAll(xiter.FromSlice(in), v)
}

// Count returns the number of items in the slice.
//
// EXAMPLE:
//
//	xslice.Count([]int{1, 2, 3}) 👉 3
//	xslice.Count([]int{}) 👉 0
func Count[T any](in []T) int {
	return len(in)
}

// Find returns the first item in the slice that satisfies the condition provided by f.
//
// EXAMPLE:
//
//	xslice.Find([]int{1, 2, 3}, func(x int) bool { return x == 1 })  👉 1, true
//	xslice.Find([]int{1, 2, 3}, func(x int) bool { return x == -1 }) 👉 0, false
func Find[T any](in []T, f func(T) bool) (val T, found bool) {
	return xiter.Find(xiter.FromSlice(in), f)
}

// FindO returns the first item in the slice that satisfies the condition provided by f.
//
// EXAMPLE:
//
//	xslice.FindO(_range(0, 10), func(x int) bool { return x == 1 }).Must() 👉 1
//	xslice.FindO(_range(0, 10), func(x int) bool { return x == -1 }).Ok() 👉 false
func FindO[T any](in []T, f func(T) bool) optional.O[T] {
	return xiter.FindO(xiter.FromSlice(in), f)
}

// ForEach iterates over each item in the slice, stop if f returns false.
//
// EXAMPLE:
//
//	ForEach([]int{1, 2, 3}, func(x int) bool {
//		fmt.Println(x)
//		return true
//	}
//	Output:
//	1
//	2
//	3
func ForEach[T any](in []T, f func(T) bool) {
	xiter.ForEach(xiter.FromSlice(in), f)
}

// ForEachIdx iterates over each item in the slice, stop if f returns false.
//
// EXAMPLE:
//
//	ForEach([]int{1, 2, 3}, func(idx, x int) bool {
//		fmt.Println(idx, x)
//		return true
//	}
//	Output:
//	0 1
//	1 2
//	2 3
func ForEachIdx[T any](in []T, f func(idx int, v T) bool) {
	xiter.ForEachIdx(xiter.FromSlice(in), f)
}

// HeadO returns the first item in the slice.
//
// EXAMPLE:
//
//	xslice.HeadO(_range(0, 10)).Must() 👉 0
//	xslice.HeadO(_range(0, 0)).Ok() 👉 false
func HeadO[T any](in []T) optional.O[T] {
	return xiter.HeadO(xiter.FromSlice(in))
}

// Head returns the first item in the slice.
//
// EXAMPLE:
//
//	optional.FromValue2(xslice.Head(_range(0, 10))).Must() 👉 0
//	optional.FromValue2(xslice.Head(_range(0, 0))).Ok() 👉 false
func Head[T any](in []T) (v T, hasOne bool) {
	return xiter.Head(xiter.FromSlice(in))
}

// Join joins the slice with sep.
//
// EXAMPLE:
//
//	xslice.Join([]string{"1", "2", "3"}, ".") 👉 "1.2.3"
//	xslice.Join([]string{}, ".") 👉 ""
func Join[T ~string](in []T, sep T) T {
	if len(in) == 0 {
		return ""
	}
	if len(in) == 1 {
		return in[0]
	}
	n := len(sep) * (len(in) - 1)
	for _, s := range in {
		n += len(s)
	}
	b := make([]byte, n)
	bp := copy(b, string(in[0]))
	for _, s := range in[1:] {
		bp += copy(b[bp:], string(sep))
		bp += copy(b[bp:], string(s))
	}
	return T(string(b))
}

// Min returns the minimum value in the slice.
//
// EXAMPLE:
//
//	xslice.Min([]int{1, 2, 3}) 👉 1
//	xslice.Min([]int{}) 👉 0
func Min[T constraints.Ordered](in []T) optional.O[T] {
	if len(in) == 0 {
		return optional.Empty[T]()
	}
	m := in[0]
	for _, v := range in[1:] {
		if v < m {
			m = v
		}
	}
	return optional.FromValue(m)
}

// MinN returns the minimum value in the slice.
//
// EXAMPLE:
//
//	xslice.MinN(1, 2, 3) 👉 1
func MinN[T constraints.Ordered](in ...T) optional.O[T] {
	return Min(in)
}

// MinBy returns the minimum value evaluated by f in the slice.
//
// EXAMPLE:
//
//	xslice.MinBy([]int{3, 2, 1} /*less = */, func(a, b int) bool { return a > b }).Must() 👉 3
func MinBy[T constraints.Ordered](in []T, f func(T, T) bool) optional.O[T] {
	if len(in) == 0 {
		return optional.Empty[T]()
	}
	m := in[0]
	for _, v := range in[1:] {
		if f(v, m) {
			m = v
		}
	}
	return optional.FromValue(m)
}

// Max returns the maximum value in the slice.
//
// EXAMPLE:
//
//	xslice.Max([]int{1, 2, 3}) 👉 3
//	xslice.Max([]int{}) 👉 0
func Max[T constraints.Ordered](in []T) optional.O[T] {
	if len(in) == 0 {
		return optional.Empty[T]()
	}
	m := in[0]
	for _, v := range in[1:] {
		if v > m {
			m = v
		}
	}
	return optional.FromValue(m)
}

// MaxN returns the maximum value in the slice.
//
// EXAMPLE:
//
//	xslice.MaxN(1, 2, 3) 👉 3
func MaxN[T constraints.Ordered](in ...T) optional.O[T] {
	return Max(in)
}

// MaxBy returns the maximum value evaluated by f in the slice.
//
// EXAMPLE:
//
//	xslice.MaxBy([]int{1, 2, 3} /*less = */, func(a, b int) bool { return a > b }).Must() 👉 1
func MaxBy[T constraints.Ordered](in []T, f func(T, T) bool) optional.O[T] {
	if len(in) == 0 {
		return optional.Empty[T]()
	}
	m := in[0]
	for _, v := range in[1:] {
		if f(m, v) {
			m = v
		}
	}
	return optional.FromValue(m)
}

// Map returns a new slice with the results of applying the given function to every element in this slice.
//
// EXAMPLE:
//
//	xslice.Map([]int{1, 2, 3}, func(x int) int { return x * 2 }) 👉 [2, 4, 6]
//	xslice.Map([]int{1, 2, 3}, strconv.Itoa) 👉 ["1", "2", "3"]
func Map[T any, U any](in []T, f func(T) U) []U {
	out := make([]U, len(in))
	for i, v := range in {
		out[i] = f(v)
	}
	return out
}

// Clone returns a copy of the slice.
//
// EXAMPLE:
//
//	xslice.Clone([]int{1, 2, 3}) 👉 [1, 2, 3]
func Clone[T any](in []T) []T {
	return append([]T(nil), in...)
}

// CloneBy returns a copy of the slice with the results of applying the given function to every element in this slice.
//
// EXAMPLE:
//
//	xslice.CloneBy([]int{1, 2, 3}, func(x int) int { return x * 2 }) 👉 [2, 4, 6]
//	xslice.CloneBy([]int{1, 2, 3}, strconv.Itoa) 👉 ["1", "2", "3"]
func CloneBy[T any, U any](in []T, f func(T) U) []U {
	if in == nil {
		return nil
	}
	return Map(in, f)
}

// Concat concatenates the slices.
//
// EXAMPLE:
//
//	xslice.Concat([]int{1, 2, 3}, []int{4, 5, 6}) 👉 [1, 2, 3, 4, 5, 6]
//	xslice.Concat([]int{1, 2, 3}, []int{}) 👉 [1, 2, 3]
func Concat[T any](vs ...[]T) []T {
	n := 0
	for _, v := range vs {
		n += len(v)
	}
	out := make([]T, 0, n)
	for _, v := range vs {
		out = append(out, v...)
	}
	return out
}

// Subset returns a subset slice from the slice.
// if start < -1 means that we take subset from right-to-left
//
// EXAMPLE:
//
//	xslice.Subset([]int{1, 2, 3}, 0, 2) 👉 [1, 2]
//	xslice.Subset([]int{1, 2, 3}, -1, 2) 👉 [2, 3]
func Subset[T any, Slice ~[]T](in Slice, start, count int) Slice {
	if count < 0 {
		count = 0
	}
	n := len(in)
	if start >= n || -start > n {
		return nil
	}
	if start < 0 {
		start = n + start
	}
	if start+count > n {
		count = n - start
	}
	out := make(Slice, count)
	for i := 0; i < count; i++ {
		out[i] = in[start+i]
	}
	return out
}

// SubsetInPlace returns a subset slice copied from the slice.
// if start < -1 means that we take subset from right-to-left
// EXAMPLE:
//
//	xslice.SubsetInPlace([]int{1, 2, 3}, 0, 2) 👉 [1, 2]
//	xslice.SubsetInPlace([]int{1, 2, 3}, -1, 2) 👉 [2, 3]
func SubsetInPlace[T any, Slice ~[]T](in Slice, start int, count int) Slice {
	size := len(in)

	if start < 0 {
		start = size + start
		if start < 0 {
			return Slice{}
		}
	}
	if start > size {
		return Slice{}
	}

	if count > size-start {
		count = size - start
	}
	return in[start : start+count]
}

// Replace replaces the count elements in the slice from 'from' to 'to'.
//
// EXAMPLE:
//
//	xslice.Replace([]int{1, 2, 3}, 2, 4, 1) 👉 [1, 4, 3]
//	xslice.Replace([]int{1, 2, 2}, 2, 4, -1) 👉 [1, 4, 4]
func Replace[T comparable, Slice ~[]T](in Slice, from, to T, count int) []T {
	if count == 0 {
		return Clone(in)
	}
	out := Clone(in)
	replaced := 0
	for i, v := range out {
		if v == from {
			out[i] = to
			replaced++
			if count > 0 && replaced >= count {
				break
			}
		}
	}
	return out
}

// ReplaceAll replaces all elements in the slice from 'from' to 'to'.
//
// EXAMPLE:
//
//	xslice.ReplaceAll([]int{1, 2, 3}, 2, 4) 👉 [1, 4, 3]
//	xslice.ReplaceAll([]int{1, 2, 2}, 2, 4) 👉 [1, 4, 4]
func ReplaceAll[T comparable, Slice ~[]T](in Slice, from, to T) []T {
	return Replace(in, from, to, -1)
}

// ReverseClone reverses the slice.
//
// EXAMPLE:
//
//	xslice.ReverseClone([]int{1, 2, 3}) 👉 [3, 2, 1]
//	xslice.ReverseClone([]int{}) 👉 []int{}
//	xslice.ReverseClone([]int{3, 2, 1}) 👉 [1, 2, 3]
func ReverseClone[T any, Slice ~[]T](in Slice) Slice {
	out := make(Slice, len(in))
	for i, v := range in {
		out[len(in)-1-i] = v
	}
	return out
}

// Reverse reverses the slice.
//
// EXAMPLE:
//
//	xslice.Reverse([]int{1, 2, 3}) 👉 [3, 2, 1]
//	xslice.Reverse([]int{}) 👉 []int{}
func Reverse[T any, Slice ~[]T](in Slice) {
	for i, j := 0, len(in)-1; i < j; i, j = i+1, j-1 {
		in[i], in[j] = in[j], in[i]
	}
}

// Repeat returns a new slice with the elements repeated 'count' times.
//
// EXAMPLE:
//
//	xslice.Repeat([]int{1, 2, 3}, 3) 👉 [1, 2, 3, 1, 2, 3, 1, 2, 3]
//	xslice.Repeat([]int{1, 2, 3}, 0) 👉 []int{}
func Repeat[T any, Slice ~[]T](in Slice, count int) Slice {
	if count <= 0 {
		return nil
	}
	out := make(Slice, 0, len(in)*count)
	for i := 0; i < count; i++ {
		out = append(out, in...)
	}
	return out
}

// RepeatBy returns a new slice with the elements return by f repeated 'count' times.
//
// EXAMPLE:
//
//	xslice.RepeatBy(3, func(i int) int { return i }) 👉 [0, 1, 2]
//	xslice.RepeatBy(3, func(i int) string { return strconv.Itoa(i) }) 👉 []string{"1", "2", "3"}
func RepeatBy[T any](n int, f func(i int) T) []T {
	out := make([]T, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, f(i))
	}
	return out
}

// Shuffle shuffles the slice.
//
// EXAMPLE:
//
//	xslice.Shuffle([]int{1, 2, 3}) 👉 [2, 1, 3] (random)
//	xslice.Shuffle([]int{}) 👉 []int{}
func Shuffle[T any, Slice ~[]T](in Slice) Slice {
	out := make(Slice, len(in))
	copy(out, in)
	rand.Shuffle(len(out), func(i, j int) {
		out[i], out[j] = out[j], out[i]
	})
	return out
}

// ShuffleInPlace shuffles the slice.
//
// EXAMPLE:
//
//	array := []int{1, 2, 3}
//	xslice.ShuffleInPlace(array) 👉 [2, 1, 3] (random)
func ShuffleInPlace[T any, Slice ~[]T](in Slice) {
	// why we do not use slices.Shuffle() directly?
	// because lower version golang may has not package "slices"
	rand.Shuffle(len(in), func(i, j int) {
		in[i], in[j] = in[j], in[i]
	})
}

// Chunk returns a new slice with the elements in the slice chunked into smaller slices of the specified size.
//
// EXAMPLE:
//
//	xslice.Chunk([]int{1, 2, 3, 4, 5}, 2) 👉 [[1, 2], [3, 4], [5]]
//	xslice.Chunk([]int{1, 2, 3, 4, 5}, 10) 👉 [[1, 2, 3, 4, 5]]
//	xslice.Chunk([]int{1, 2, 3, 4, 5}, 0) 👉 []int{}
func Chunk[T any, Slice ~[]T](in Slice, chunkSize int) []Slice {
	xassert.MustBePositive(chunkSize)
	n := len(in)
	out := make([]Slice, 0, n/chunkSize+1)
	for i := 0; i < n; i += chunkSize {
		end := i + chunkSize
		if end > n {
			end = n
		}
		out = append(out, append(Slice(nil), in[i:end]...))
	}
	return out
}

// ChunkInPlace returns a new slice with the elements in the slice chunked into smaller slices of the specified size.
// This function will not copy the elements, has no extra costs.
// EXAMPLE:
//
//	xslice.Chunk([]int{1, 2, 3, 4, 5}, 2) 👉 [[1, 2], [3, 4], [5]]
//	xslice.Chunk([]int{1, 2, 3, 4, 5}, 10) 👉 [[1, 2, 3, 4, 5]]
//	xslice.Chunk([]int{1, 2, 3, 4, 5}, 0) 👉 []int{}
func ChunkInPlace[T any, Slice ~[]T](in Slice, chunkSize int) []Slice {
	xassert.MustBePositive(chunkSize)
	out := make([]Slice, 0, len(in)/chunkSize+1)
	for i := 0; i < len(in); i += chunkSize {
		end := i + chunkSize
		if end > len(in) {
			end = len(in)
		}
		out = append(out, in[i:end])
	}
	return out
}

// Index returns the index of the first element in the slice that is equal to v.
// If no such element is found, -1 is returned.
// EXAMPLE:
//
//	xslice.Index([]int{1, 2, 3, 4, 5}, 1) 👉 0
//	xslice.Index([]int{1, 2, 3, 4, 5}, 3) 👉 2
//	xslice.Index([]int{1, 2, 3, 4, 5}, 666) 👉 -1
func Index[T comparable, Slice ~[]T](in Slice, v T) int {
	return xiter.Index(xiter.FromSlice(in), v)
}

// Sum returns the sum of all elements in the slice.
//
// EXAMPLE:
//
//	xslice.Sum([]int{1, 2, 3}) 👉 6
//	xslice.Sum([]int{}) 👉 0
func Sum[T constraints.Number, Slice ~[]T](in Slice) T {
	var sum T
	for _, v := range in {
		sum += v
	}
	return sum
}

// SumN returns the sum of all input arguments.
//
// EXAMPLE:
//
//	xslice.SumN(1, 2, 3) 👉 6
//	xslice.SumN() 👉 0
func SumN[T constraints.Number](in ...T) T {
	return Sum(in)
}

// SumBy returns the sum of all elements in the slice after applying the given function f to each element.
//
// EXAMPLE:
//
//	xslice.SumBy([]string{"1", "2", "3"}, func(x string) int {
//		i, _ := strconv.Atoi(x)
//		return i
//	}) 👉 6
//	xslice.SumBy([]string{}, func(x string) int { return 0 }) 👉 0
func SumBy[T any, R constraints.Number, Slice ~[]T](in Slice, f func(T) R) R {
	var sum R
	for _, v := range in {
		sum += f(v)
	}
	return sum
}

// Uniq returns a new slice with the duplicate elements removed.
//
// EXAMPLE:
//
//	xslice.Uniq([]int{1, 2, 3, 2, 4}) 👉 [1, 2, 3, 4]
func Uniq[T comparable, Slice ~[]T](in Slice) Slice {
	seen := make(map[T]struct{}, len(in)/2)
	out := make(Slice, 0, len(in))
	for _, v := range in {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}

// GroupBy returns a map of the slice elements grouped by the given function f.
//
// EXAMPLE:
//
//	xslice.GroupBy([]int{1, 2, 3, 2, 4}, func(x int) int { return x % 2 }) 👉 map[0:[2 4] 1:[1 3]]
func GroupBy[T any, K comparable, Slice ~[]T](in Slice, f func(T) K) map[K]Slice {
	out := make(map[K]Slice, len(in)/2)
	for _, v := range in {
		k := f(v)
		out[k] = append(out[k], v)
	}
	return out
}

// GroupByMap returns a map of the slice elements grouped by the given function f.
//
// EXAMPLE:
//
//	xslice.GroupByMap([]int{1, 2, 3, 2, 4}, func(x int) (int, int) { return x % 2, x }) 👉 map[0:[2 4] 1:[1 3]]
func GroupByMap[T any, Slice ~[]T, K comparable, V any](in Slice, f func(T) (K, V)) map[K][]V {
	out := make(map[K][]V, len(in)/2)
	for _, v := range in {
		k, mv := f(v)
		out[k] = append(out[k], mv)
	}
	return out
}

// Filter returns a new slice with the elements that satisfy the given function f.
//
// EXAMPLE:
//
//	xslice.Filter([]int{1, 2, 3, 2, 4}, func(x int) bool { return x%2 == 0 }) 👉 [2 4]
func Filter[T any, Slice ~[]T](in Slice, f func(T) bool) Slice {
	out := make(Slice, 0, len(in)/2)
	for _, v := range in {
		if f(v) {
			out = append(out, v)
		}
	}
	return out
}

// Compact returns a new slice with the zero elements removed.
//
// EXAMPLE:
//
//	xslice.Compact([]int{0, 1, 2, 3, 4}) 👉 [1 2 3 4]
func Compact[T comparable, Slice ~[]T](in Slice) Slice {
	var zero T
	out := make(Slice, 0, len(in))
	for _, v := range in {
		if v != zero {
			out = append(out, v)
		}
	}
	return out
}

// First returns the first element in the slice.
// If the slice is empty, the zero value of T is returned.
// EXAMPLE:
//
//	xslice.First([]int{1, 2, 3}) 👉 1
//	xslice.First([]int{}) 👉 0
func First[T any, Slice ~[]T](in Slice) (T, bool) {
	return xiter.First(xiter.FromSlice(in))
}

// FirstO returns the first element in the slice as an optional.O[T].
// If the slice is empty, the zero value of T is returned.
// EXAMPLE:
//
//	xslice.FirstO([]int{1, 2, 3}) 👉 1
//	xslice.FirstO([]int{}) 👉 0
func FirstO[T any, Slice ~[]T](in Slice) optional.O[T] {
	return optional.FromValue2(First(in))
}

// Last returns the last element in the slice.
// If the slice is empty, the zero value of T is returned.
// EXAMPLE:
//
//	xslice.Last([]int{1, 2, 3}) 👉 3
//	xslice.Last([]int{}) 👉 0
func Last[T any, Slice ~[]T](in Slice) (T, bool) {
	if len(in) == 0 {
		var zero T
		return zero, false
	}
	return in[len(in)-1], true
}

// LastO returns the last element in the slice as an optional.O[T].
// If the slice is empty, the zero value of T is returned.
// EXAMPLE:
//
//	xslice.LastO([]int{1, 2, 3}) 👉 3
//	xslice.LastO([]int{}) 👉 0
func LastO[T any, Slice ~[]T](in Slice) optional.O[T] {
	return optional.FromValue2(Last(in))
}

// Difference returns two slices: the first slice contains the elements that are in the left slice but not in the right slice,
// and the second slice contains the elements that are in the right slice but not in the left slice.
//
// EXAMPLE:
//
//	left := []int{1, 2, 3, 4, 5}
//	right := []int{4, 5, 6, 7, 8}
//	onlyLeft, onlyRight := xslice.Difference(left, right)
//	fmt.Println(onlyLeft)  // [1 2 3]
//	fmt.Println(onlyRight) // [6 7 8]
func Difference[T comparable, Slice ~[]T](left, right Slice) (onlyLeft, onlyRight Slice) {
	rightSet := make(map[T]struct{}, len(right))
	for _, v := range right {
		rightSet[v] = struct{}{}
	}
	leftSet := make(map[T]struct{}, len(left))
	for _, v := range left {
		leftSet[v] = struct{}{}
	}
	onlyLeft = make(Slice, 0, len(left)/2)
	for _, v := range left {
		if _, ok := rightSet[v]; !ok {
			onlyLeft = append(onlyLeft, v)
		}
	}
	onlyRight = make(Slice, 0, len(right)/2)
	for _, v := range right {
		if _, ok := leftSet[v]; !ok {
			onlyRight = append(onlyRight, v)
		}
	}
	return
}

// Intersect returns a slice that contains the elements that are in both left and right slices.
//
// EXAMPLE:
//
//	left := []int{1, 2, 3, 4, 5}
//	right := []int{4, 5, 6, 7, 8}
//	intersect := xslice.Intersect(left, right)
//	fmt.Println(intersect) // [4 5]
func Intersect[T comparable, Slice ~[]T](left, right Slice) Slice {
	var smaller, larger Slice
	if len(left) > len(right) {
		smaller, larger = right, left
	} else {
		smaller, larger = left, right
	}
	set := make(map[T]struct{}, len(smaller))
	for _, v := range smaller {
		set[v] = struct{}{}
	}
	out := make(Slice, 0, len(smaller))
	for _, v := range larger {
		if _, ok := set[v]; ok {
			out = append(out, v)
		}
	}
	return out
}

// Union returns a slice that contains all elements in left and right slices.
//
// EXAMPLE:
//
//	left := []int{1, 2, 3, 4}
//	right := []int{3, 4, 5, 6}
//	union := xslice.Union(left, right)
//	fmt.Println(union) // [1 2 3 4 5 6]
func Union[T comparable, Slice ~[]T](left, right Slice) Slice {
	var smaller, larger Slice
	if len(left) <= len(right) {
		smaller, larger = left, right
	} else {
		smaller, larger = right, left
	}
	set := make(map[T]struct{}, len(smaller))
	for _, v := range smaller {
		set[v] = struct{}{}
	}
	out := make(Slice, 0, len(smaller)+len(larger)/2)
	out = append(out, smaller...)
	for _, v := range larger {
		if _, ok := set[v]; !ok {
			out = append(out, v)
		}
	}
	return out
}

// Remove returns a slice that remove all elements in wantToRemove
//
// EXAMPLE:
//
//	arr := []int{1, 2, 3, 4}
//	arr1 := xslice.Remove(arr, 1)
//	fmt.Println(arr1) // [2, 3, 4]
func Remove[T comparable, Slice ~[]T](in Slice, wantToRemove ...T) Slice {
	if len(wantToRemove) == 0 {
		return Clone(in)
	}
	removeSet := make(map[T]struct{}, len(wantToRemove))
	for _, v := range wantToRemove {
		removeSet[v] = struct{}{}
	}
	out := make(Slice, 0, len(in))
	for _, v := range in {
		if _, ok := removeSet[v]; !ok {
			out = append(out, v)
		}
	}
	return out
}

// Flatten returns a new slice with all nested slices flattened into a single slice.
//
// EXAMPLE:
//
//	xslice.Flatten([][]int{{1, 2}, {3, 4}, {5}}) 👉 [1, 2, 3, 4, 5]
//	xslice.Flatten([][]int{{1, 2}, {}, {3, 4}}) 👉 [1, 2, 3, 4]
//	xslice.Flatten([][]int{}) 👉 []int{}
//	xslice.Flatten([][]int{{}, {}, {}}) 👉 []int{}
func Flatten[T any](in [][]T) []T {
	return Concat(in...)
}

// ToMap returns a map where keys are elements from the slice and values are the result of applying f to each element.
// If there are duplicate keys in the slice, only the last element with that key will be present in the map.
//
// EXAMPLE:
//
//	xslice.ToMap([]string{"a", "b", "c"}, func(s string) int { return len(s) }) 👉 map[a:1 b:1 c:1]
//	xslice.ToMap([]int{1, 2, 3}, func(i int) string { return fmt.Sprintf("num_%d", i) }) 👉 map[1:num_1 2:num_2 3:num_3]
//	xslice.ToMap([]int{1, 2, 1, 3}, func(i int) string { return fmt.Sprintf("val_%d", i) }) 👉 map[1:val_1 2:val_2 3:val_3] (note: key 1 has "val_1" from the last occurrence)
//	xslice.ToMap([]int{}, func(i int) string { return "" }) 👉 map[int]string{}
func ToMap[T comparable, U any](in []T, f func(T) U) map[T]U {
	out := make(map[T]U, len(in))
	for _, v := range in {
		out[v] = f(v)
	}
	return out
}

// Sample returns a new slice with n randomly selected elements from the input slice.
// If n is greater than the length of the slice, it returns all elements in random order.
// If n is less than or equal to 0, it returns an empty slice.
//
// EXAMPLE:
//
//	xslice.Sample([]int{1, 2, 3, 4, 5}, 3) 👉 [3, 1, 5] (random order, 3 elements)
//	xslice.Sample([]int{1, 2, 3}, 5) 👉 [2, 1, 3] (random order, all elements)
//	xslice.Sample([]int{1, 2, 3}, 0) 👉 []int{}
//	xslice.Sample([]int{}, 3) 👉 []int{}
func Sample[T any, Slice ~[]T](in Slice, n int) Slice {
	if n <= 0 || len(in) == 0 {
		return nil
	}
	if n >= len(in) {
		return Shuffle(in)
	}
	out := make(Slice, len(in))
	copy(out, in)
	for i := 0; i < n; i++ {
		j := rand.Intn(len(in)-i) + i //nolint:gosec
		out[i], out[j] = out[j], out[i]
	}
	return out[:n]
}

// RandomElement returns a random element from the slice as an optional.O[T].
// If the slice is empty, it returns an optional.O[T] with Ok() == false.
//
// EXAMPLE:
//
//	xslice.RandomElement([]int{1, 2, 3, 4, 5}) 👉 3 (random element)
//	xslice.RandomElement([]int{42}) 👉 42 (always returns the only element)
//	xslice.RandomElement([]int{}).Ok() 👉 false
func RandomElement[T any, Slice ~[]T](in Slice) optional.O[T] {
	if len(in) == 0 {
		return optional.Empty[T]()
	}
	return optional.FromValue(in[rand.Intn(len(in))]) //nolint:gosec
}

// CountBy counts occurrences of each key in the slice, returning a map of keys to counts.
//
// EXAMPLE:
//
//	xslice.CountBy([]int{1, 2, 3, 2, 1, 2}, func(x int) int { return x })
//	// 👉 map[int]int{1: 2, 2: 3, 3: 1}
func CountBy[T any, K comparable](in []T, fn func(T) K) map[K]int {
	result := make(map[K]int)
	for _, v := range in {
		result[fn(v)]++
	}
	return result
}

// KeyBy creates a map from the slice using the key function.
// Later elements overwrite earlier ones for duplicate keys.
//
// EXAMPLE:
//
//	xslice.KeyBy([]int{1, 2, 3}, func(x int) int { return x * 10 })
//	// 👉 map[int]int{10: 1, 20: 2, 30: 3}
func KeyBy[T any, K comparable](in []T, fn func(T) K) map[K]T {
	result := make(map[K]T, len(in))
	for _, v := range in {
		result[fn(v)] = v
	}
	return result
}

// Partition splits a slice into two slices based on a predicate.
// The first return slice contains elements where fn returns true.
//
// EXAMPLE:
//
//	yes, no := xslice.Partition([]int{1, 2, 3, 4, 5}, func(x int) bool { return x%2 == 0 })
//	// yes 👉 []int{2, 4}, no 👉 []int{1, 3, 5}
func Partition[T any, Slice ~[]T](in Slice, fn func(T) bool) (yes, no Slice) {
	yes = make(Slice, 0)
	no = make(Slice, 0)
	for _, v := range in {
		if fn(v) {
			yes = append(yes, v)
		} else {
			no = append(no, v)
		}
	}
	return
}

// FlatMap maps each element to a slice and flattens the results into a single slice.
//
// EXAMPLE:
//
//	xslice.FlatMap([]int{1, 2, 3}, func(x int) []int { return []int{x, x * 10} })
//	// 👉 []int{1, 10, 2, 20, 3, 30}
func FlatMap[T any, U any](in []T, fn func(T) []U) []U {
	result := make([]U, 0)
	for _, v := range in {
		result = append(result, fn(v)...)
	}
	return result
}

// IsSorted checks if the slice is sorted in ascending order.
// Empty and single-element slices are considered sorted.
//
// EXAMPLE:
//
//	xslice.IsSorted([]int{1, 2, 3, 4}) 👉 true
//	xslice.IsSorted([]int{1, 3, 2, 4}) 👉 false
func IsSorted[T constraints.Ordered](in []T) bool {
	for i := 0; i < len(in)-1; i++ {
		if in[i] > in[i+1] {
			return false
		}
	}
	return true
}

// AllEqual checks if all elements in the slice are equal.
// Empty and single-element slices are considered to have all equal elements.
//
// EXAMPLE:
//
//	xslice.AllEqual([]int{1, 1, 1, 1}) 👉 true
//	xslice.AllEqual([]int{1, 2, 1, 1}) 👉 false
func AllEqual[T comparable](in []T) bool {
	for i := 1; i < len(in); i++ {
		if in[i] != in[0] {
			return false
		}
	}
	return true
}

// MinMax returns the minimum and maximum elements in the slice in a single pass.
//
// EXAMPLE:
//
//	xslice.MinMax([]int{3, 1, 4, 1, 5, 9}) 👉 (1, 9, true)
//	xslice.MinMax([]int{}) 👉 (0, 0, false)
func MinMax[T constraints.Ordered](in []T) (min T, max T, ok bool) {
	if len(in) == 0 {
		return
	}
	min, max = in[0], in[0]
	for _, v := range in[1:] {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return min, max, true
}

// Mode returns the most frequently occurring element in the slice.
// If the slice is empty, it returns an empty optional.
// If there are multiple modes (tie), the first one to reach the maximum count is returned.
//
// EXAMPLE:
//
//	xslice.Mode([]int{1, 2, 3, 2, 1, 2}) 👉 2
//	xslice.Mode([]int{}) 👉 optional.Empty[int]()
func Mode[T comparable](in []T) optional.O[T] {
	if len(in) == 0 {
		return optional.Empty[T]()
	}
	counts := make(map[T]int)
	maxCount := 0
	var mode T
	for _, v := range in {
		counts[v]++
		if counts[v] > maxCount {
			maxCount = counts[v]
			mode = v
		}
	}
	return optional.FromValue(mode)
}
