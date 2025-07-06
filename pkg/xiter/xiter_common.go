package xiter

import (
	gassert "github.com/dashjay/xiter/pkg/internal/xassert"
	"github.com/dashjay/xiter/pkg/optional"
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
	elements := PullOut(seq, index+1)
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
	all := PullOut(seq, -1)
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
