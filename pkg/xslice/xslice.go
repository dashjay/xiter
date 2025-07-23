package xslice

import (
	"math/rand"

	"github.com/dashjay/xiter/pkg/internal/constraints"
	"github.com/dashjay/xiter/pkg/internal/xassert"
	"github.com/dashjay/xiter/pkg/optional"
	"github.com/dashjay/xiter/pkg/xiter"
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
	return xiter.AvgFromSeq(xiter.FromSlice(in))
}

// AvgN returns the average value of the items
//
// EXAMPLE:
//
//	xslice.AvgN(1, 2, 3) 👉 float(2)
//	xslice.AvgN() 👉 float(0)
func AvgN[T constraints.Number](inputs ...T) float64 {
	return xiter.AvgFromSeq(xiter.FromSlice(inputs))
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
	return xiter.AvgByFromSeq(xiter.FromSlice(in), f)
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
	return xiter.Count(xiter.FromSlice(in))
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
	return xiter.Join(xiter.FromSlice(in), sep)
}

// Min returns the minimum value in the slice.
//
// EXAMPLE:
//
//	xslice.Min([]int{1, 2, 3}) 👉 1
//	xslice.Min([]int{}) 👉 0
func Min[T constraints.Ordered](in []T) optional.O[T] {
	return xiter.Min(xiter.FromSlice(in))
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
	return xiter.MinBy(xiter.FromSlice(in), f)
}

// Max returns the maximum value in the slice.
//
// EXAMPLE:
//
//	xslice.Max([]int{1, 2, 3}) 👉 3
//	xslice.Max([]int{}) 👉 0
func Max[T constraints.Ordered](in []T) optional.O[T] {
	return xiter.Max(xiter.FromSlice(in))
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
	return xiter.MaxBy(xiter.FromSlice(in), f)
}

// Map returns a new slice with the results of applying the given function to every element in this slice.
//
// EXAMPLE:
//
//	xslice.Map([]int{1, 2, 3}, func(x int) int { return x * 2 }) 👉 [2, 4, 6]
//	xslice.Map([]int{1, 2, 3}, strconv.Itoa) 👉 ["1", "2", "3"]
func Map[T any, U any](in []T, f func(T) U) []U {
	out := make([]U, len(in))
	for i := range in {
		out[i] = f(in[i])
	}
	return out
}

// Clone returns a copy of the slice.
//
// EXAMPLE:
//
//	xslice.Clone([]int{1, 2, 3}) 👉 [1, 2, 3]
func Clone[T any](in []T) []T {
	if in == nil {
		return nil
	}
	return xiter.ToSlice(xiter.FromSlice(in))
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
	var seqs = make([]xiter.Seq[T], 0, len(vs))
	for _, v := range vs {
		seqs = append(seqs, xiter.FromSlice(v))
	}
	return xiter.ToSlice(xiter.Concat(seqs...))
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
	if start >= len(in) || -start > len(in) {
		return nil
	}
	if start >= 0 {
		return xiter.ToSlice(xiter.Limit(xiter.Skip(xiter.FromSlice(in), start), count))
	} else {
		return xiter.ToSlice(xiter.Limit(xiter.Skip(xiter.FromSlice(in), len(in)+start), count))
	}
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
	return xiter.ToSlice(xiter.Replace(xiter.FromSlice(in), from, to, count))
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
	// why we do not use slices.Reverse() directly ?
	// because lower version golang may has not package "slices"
	return xiter.ToSlice(xiter.FromSliceReverse(in))
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
	return xiter.ToSlice(xiter.Repeat(xiter.FromSlice(in), count))
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
	return xiter.ToSlice(xiter.FromSliceShuffle(in))
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
	out := make([]Slice, 0, len(in)/chunkSize+1)
	seq := xiter.FromSlice(in)
	for {
		res := xiter.ToSlice(xiter.Limit(xiter.Skip(seq, len(out)*chunkSize), chunkSize))
		if len(res) == 0 {
			break
		}
		out = append(out, res)
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
	return xiter.Sum(xiter.FromSlice(in))
}

// SumN returns the sum of all input arguments.
//
// EXAMPLE:
//
//	xslice.SumN(1, 2, 3) 👉 6
//	xslice.SumN() 👉 0
func SumN[T constraints.Number](in ...T) T {
	return xiter.Sum(xiter.FromSlice(in))
}

//
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
	return xiter.Sum(xiter.Map(f, xiter.FromSlice(in)))
}

// Uniq returns a new slice with the duplicate elements removed.
//
// EXAMPLE:
//
//	xslice.Uniq([]int{1, 2, 3, 2, 4}) 👉 [1, 2, 3, 4]
func Uniq[T comparable, Slice ~[]T](in Slice) Slice {
	return xiter.ToSlice(xiter.Uniq(xiter.FromSlice(in)))
}

// GroupBy returns a map of the slice elements grouped by the given function f.
//
// EXAMPLE:
//
//	xslice.GroupBy([]int{1, 2, 3, 2, 4}, func(x int) int { return x % 2 }) 👉 map[0:[2 4] 1:[1 3]]
func GroupBy[T any, K comparable, Slice ~[]T](in Slice, f func(T) K) map[K]Slice {
	seq2 := xiter.MapToSeq2(xiter.FromSlice(in), f)
	return xiter.Reduce2(func(sum map[K]Slice, k K, v T) map[K]Slice {
		sum[k] = append(sum[k], v)
		return sum
	}, map[K]Slice{}, seq2)
}

// GroupByMap returns a map of the slice elements grouped by the given function f.
//
// EXAMPLE:
//
//	xslice.GroupByMap([]int{1, 2, 3, 2, 4}, func(x int) (int, int) { return x % 2, x }) 👉 map[0:[2 4] 1:[1 3]]
func GroupByMap[T any, Slice ~[]T, K comparable, V any](in Slice, f func(T) (K, V)) map[K][]V {
	seq2 := xiter.MapToSeq2Value(xiter.FromSlice(in), f)
	return xiter.Reduce2(func(sum map[K][]V, k K, v V) map[K][]V {
		sum[k] = append(sum[k], v)
		return sum
	}, map[K][]V{}, seq2)
}
