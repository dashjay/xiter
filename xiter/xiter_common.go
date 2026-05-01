package xiter

import (
	"github.com/dashjay/xiter/internal/constraints"
	gassert "github.com/dashjay/xiter/internal/xassert"
	"github.com/dashjay/xiter/optional"
)

// A Zipped2 is a pair of zipped key-value pairs,
// one of which may be missing, drawn from two different sequences.
type Zipped2[K1, V1, K2, V2 any] struct {
	K1  K1
	V1  V1
	Ok1 bool // whether K1, V1 are present (if not, they will be zero)
	K2  K2
	V2  V2
	Ok2 bool // whether K2, V2 are present (if not, they will be zero)
}

// A Zipped is a pair of zipped values, one of which may be missing,
// drawn from two different sequences.
type Zipped[V1, V2 any] struct {
	V1  V1
	Ok1 bool // whether V1 is present (if not, it will be zero)
	V2  V2
	Ok2 bool // whether V2 is present (if not, it will be zero)
}

// FromSlice received a slice and returned a Seq for this slice.
func FromSlice[T any](in []T) Seq[T] {
	return func(yield func(T) bool) {
		for i := 0; i < len(in); i++ {
			if !yield(in[i]) {
				break
			}
		}
	}
}

// FromSliceIdx received a slice and returned a Seq2 for this slice, key is index.
func FromSliceIdx[T any](in []T) Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i := 0; i < len(in); i++ {
			if !yield(i, in[i]) {
				break
			}
		}
	}
}

// At return the element at index from seq.
func At[T any](seq Seq[T], index int) optional.O[T] {
	gassert.MustBePositive(index)
	elements := ToSliceN(seq, index+1)
	if index >= len(elements) {
		return optional.Empty[T]()
	}
	return optional.FromValue(elements[index])
}

func FromSliceReverse[T any, Slice ~[]T](in Slice) Seq[T] {
	return func(yield func(T) bool) {
		for i := len(in) - 1; i >= 0; i-- {
			if !yield(in[i]) {
				break
			}
		}
	}
}

// Reverse return a reversed seq.
func Reverse[T any](seq Seq[T]) Seq[T] {
	all := ToSliceN(seq, -1)
	return func(yield func(T) bool) {
		for i := len(all) - 1; i >= 0; i-- {
			if !yield(all[i]) {
				break
			}
		}
	}
}

// Repeat return a seq that repeat seq for count times.
func Repeat[T any](seq Seq[T], count int) Seq[T] {
	seqs := make([]Seq[T], 0, count)
	for i := 0; i < count; i++ {
		seqs = append(seqs, seq)
	}
	return Concat(seqs...)
}

// FromChan creates a Seq from a Go channel. It yields elements from the channel
// until the channel is closed or the consumer stops iterating.
//
// Example:
//
//	ch := make(chan int, 3)
//	ch <- 1
//	ch <- 2
//	close(ch)
//
//	seq := FromChan(ch)
//
//	// Iterate over the sequence
//	_ = ToSlice(seq) // Returns []int{1, 2}
func FromChan[T any](in <-chan T) Seq[T] {
	return func(yield func(T) bool) {
		for elem := range in {
			if !yield(elem) {
				return
			}
		}
	}
}

// Difference returns two sequences: the first sequence contains elements that are in the left sequence but not in the right sequence,
// and the second sequence contains elements that are in the right sequence but not in the left sequence.
//
// EXAMPLE:
//
//	left := []int{1, 2, 3, 4}
//	right := []int{3, 4, 5, 6}
//	onlyLeft, onlyRight := Difference(FromSlice(left), FromSlice(right))
//	// onlyLeft 👉 [1 2]
//	// onlyRight 👉 [5 6]
func Difference[T comparable](left Seq[T], right Seq[T]) (onlyLeft Seq[T], onlyRight Seq[T]) {
	leftMap := ToMapFromSeq(left, func(k T) struct{} {
		return struct{}{}
	})
	rightMap := ToMapFromSeq(right, func(k T) struct{} {
		return struct{}{}
	})

	return Filter(func(v T) bool {
			_, ok := rightMap[v]
			return !ok
		}, left),
		Filter(func(v T) bool {
			_, ok := leftMap[v]
			return !ok
		}, right)
}

// Intersect return a seq that only contain elements in both left and right.
//
// EXAMPLE:
//
//	left := []int{1, 2, 3, 4}
//	right := []int{3, 4, 5, 6}
//	intersect := Intersect(FromSlice(left), FromSlice(right))
//	// intersect 👉 [3 4]
func Intersect[T comparable](left Seq[T], right Seq[T]) Seq[T] {
	leftMap := ToMapFromSeq(left, func(k T) struct{} {
		return struct{}{}
	})
	return Filter(func(v T) bool {
		_, exists := leftMap[v]
		return exists
	}, right)
}

// Union return a seq that contain all elements in left and right.
//
// EXAMPLE:
//
//	left := []int{1, 2, 3, 4}
//	right := []int{3, 4, 5, 6}
//	union := Union(FromSlice(left), FromSlice(right))
//	// union 👉 [1 2 3 4 5 6]
func Union[T comparable](left, right Seq[T]) Seq[T] {
	leftMap := ToMapFromSeq(left, func(k T) struct{} {
		return struct{}{}
	})
	return Concat(left, Filter(func(v T) bool {
		_, exists := leftMap[v]
		return !exists
	}, right))
}

// Mean return the mean of seq.
//
// EXAMPLE:
//
//	mean := Mean(FromSlice([]int{1, 2, 3, 4, 5}))
//	// mean 👉 3
func Mean[T constraints.Number](in Seq[T]) T {
	var count T = 0
	s := Reduce(func(sum T, v T) T {
		count++
		return sum + v
	}, 0, in)
	return s / count
}

// MeanBy return the mean of seq by fn.
//
// EXAMPLE:
//
//	mean := MeanBy(FromSlice([]int{1, 2, 3, 4, 5}), func(v int) int {
//		return v * 2
//	})
//	// mean 👉 6
func MeanBy[T any, R constraints.Number](in Seq[T], fn func(T) R) R {
	var count R = 0
	s := Reduce(func(sum R, v T) R {
		count++
		return sum + fn(v)
	}, 0, in)
	return s / count
}

// Moderate return the most common element in seq.
//
// EXAMPLE:
//
//	moderate := Moderate(FromSlice([]int{1, 2, 3, 4, 5, 5, 5, 6, 6, 6, 6}))
//	// moderate 👉 6
func Moderate[T comparable](in Seq[T]) (T, bool) {
	var maxTimes int
	var result T
	_ = Reduce(func(sum map[T]int, v T) map[T]int {
		sum[v]++
		if sum[v] > maxTimes {
			maxTimes = sum[v]
			result = v
		}
		return sum
	}, make(map[T]int), in)
	return result, maxTimes > 0
}

// ModerateO return the most common element in seq.
//
// EXAMPLE:
//
//	moderate := ModerateO(FromSlice([]int{1, 2, 3, 4, 5, 5, 5, 6, 6, 6, 6}))
//	// moderate 👉 6
func ModerateO[T constraints.Number](in Seq[T]) optional.O[T] {
	return optional.FromValue2(Moderate(in))
}

// Cycle returns a Seq that infinitely repeats the elements of seq.
// The input seq is materialized once, then cycled in memory.
//
// EXAMPLE:
//
//	seq := xiter.Cycle(xiter.FromSlice([]int{1, 2, 3}))
//	// seq will yield: 1, 2, 3, 1, 2, 3, 1, 2, 3, ...
func Cycle[T any](seq Seq[T]) Seq[T] {
	var elems []T
	seq(func(v T) bool {
		elems = append(elems, v)
		return true
	})
	if len(elems) == 0 {
		return func(yield func(T) bool) {}
	}
	return func(yield func(T) bool) {
		for {
			for _, v := range elems {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Generate returns an infinite Seq where each element is produced by calling fn.
// The sequence is unbounded; use with Limit, TakeWhile, etc. to constrain.
//
// EXAMPLE:
//
//	seq := xiter.Generate(func() int { return rand.Intn(100) })
//	first5 := xiter.ToSlice(xiter.Limit(seq, 5))
//	// first5 contains 5 random numbers
func Generate[T any](fn func() T) Seq[T] {
	return func(yield func(T) bool) {
		for yield(fn()) {
		}
	}
}

// ToChan sends all elements of seq to a returned channel and closes it when the seq is exhausted.
//
// EXAMPLE:
//
//	seq := xiter.FromSlice([]int{1, 2, 3})
//	ch := xiter.ToChan(seq)
//	for v := range ch {
//		fmt.Println(v)
//	}
func ToChan[T any](seq Seq[T]) <-chan T {
	ch := make(chan T)
	go func() {
		seq(func(v T) bool {
			ch <- v
			return true
		})
		close(ch)
	}()
	return ch
}

// Range returns a Seq of integers from start to end, stepping by step.
// If step > 0, elements are yielded while i < end.
// If step < 0, elements are yielded while i > end.
// If step == 0, an empty sequence is returned.
//
// EXAMPLE:
//
//	seq := xiter.Range(0, 10, 2)
//	// seq will yield: 0, 2, 4, 6, 8
//
//	seq = xiter.Range(10, 0, -3)
//	// seq will yield: 10, 7, 4, 1
func Range[T constraints.Integer](start, end, step T) Seq[T] {
	return func(yield func(T) bool) {
		if step > 0 {
			for i := start; i < end; i += step {
				if !yield(i) {
					return
				}
			}
		} else if step < 0 {
			for i := start; i > end; i += step {
				if !yield(i) {
					return
				}
			}
		}
	}
}

// WithIndex returns a Seq2 that pairs each element from seq with its 0-based index.
//
// EXAMPLE:
//
//	seq := xiter.FromSlice([]string{"a", "b", "c"})
//	for idx, v := range xiter.WithIndex(seq) {
//		fmt.Println(idx, v)
//	}
//	// output:
//	// 0 a
//	// 1 b
//	// 2 c
func WithIndex[T any](seq Seq[T]) Seq2[int, T] {
	return func(yield func(int, T) bool) {
		i := 0
		seq(func(v T) bool {
			if !yield(i, v) {
				return false
			}
			i++
			return true
		})
	}
}
