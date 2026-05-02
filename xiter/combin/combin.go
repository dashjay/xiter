//go:build go1.23

// Package combin provides combinatorial generators for iter.Seq.
//
// These functions generate combinations, permutations, and cartesian products
// as lazy sequences. Input sequences are materialized internally since
// combinatorial operations require random access to elements.
package combin

import "github.com/dashjay/xiter/xiter"

// Combinations yields all k-length combinations of elements from seq.
// Results are yielded in lexicographic order of indices.
// If seq has fewer than k elements, an empty sequence is returned.
//
// Example:
//
//	seq := xiter.FromSlice([]string{"a", "b", "c"})
//	for comb := range combin.Combinations(seq, 2) {
//		fmt.Println(comb) // [a b] [a c] [b c]
//	}
func Combinations[T any](seq xiter.Seq[T], k int) xiter.Seq[[]T] {
	return func(yield func([]T) bool) {
		elements := toSlice(seq)
		n := len(elements)
		if k <= 0 || k > n {
			return
		}
		indices := make([]int, k)
		for i := range indices {
			indices[i] = i
		}
		// First combination
		comb := make([]T, k)
		for i, idx := range indices {
			comb[i] = elements[idx]
		}
		if !yield(comb) {
			return
		}
		for {
			i := k - 1
			for i >= 0 && indices[i] == i+n-k {
				i--
			}
			if i < 0 {
				return
			}
			indices[i]++
			for j := i + 1; j < k; j++ {
				indices[j] = indices[j-1] + 1
			}
			comb := make([]T, k)
			for j, idx := range indices {
				comb[j] = elements[idx]
			}
			if !yield(comb) {
				return
			}
		}
	}
}

// Permutations yields all k-length permutations of elements from seq.
// Results are yielded in lexicographic order of indices.
// If seq has fewer than k elements, an empty sequence is returned.
//
// Example:
//
//	seq := xiter.FromSlice([]string{"a", "b", "c"})
//	for perm := range combin.Permutations(seq, 2) {
//		fmt.Println(perm) // [a b] [a c] [b a] [b c] [c a] [c b]
//	}
func Permutations[T any](seq xiter.Seq[T], k int) xiter.Seq[[]T] {
	return func(yield func([]T) bool) {
		elements := toSlice(seq)
		n := len(elements)
		if k <= 0 || k > n {
			return
		}
		selected := make([]int, k)
		used := make([]bool, n)
		var backtrack func(pos int) bool
		backtrack = func(pos int) bool {
			if pos == k {
				p := make([]T, k)
				for i, idx := range selected {
					p[i] = elements[idx]
				}
				return yield(p)
			}
			for i := 0; i < n; i++ {
				if used[i] {
					continue
				}
				used[i] = true
				selected[pos] = i
				if !backtrack(pos + 1) {
					return false
				}
				used[i] = false
			}
			return true
		}
		backtrack(0)
	}
}

// Product yields the cartesian product of the input sequences.
// All input sequences are materialized internally.
// If any input sequence is empty, an empty sequence is yielded.
//
// Example:
//
//	colors := xiter.FromSlice([]string{"red", "blue"})
//	sizes := xiter.FromSlice([]string{"S", "M", "L"})
//	for p := range combin.Product(colors, sizes) {
//		fmt.Println(p) // [red S] [red M] [red L] [blue S] [blue M] [blue L]
//	}
func Product[T any](seqs ...xiter.Seq[T]) xiter.Seq[[]T] {
	return func(yield func([]T) bool) {
		if len(seqs) == 0 {
			return
		}
		materialized := make([][]T, len(seqs))
		for i, seq := range seqs {
			for v := range seq {
				materialized[i] = append(materialized[i], v)
			}
			if len(materialized[i]) == 0 {
				return
			}
		}
		indices := make([]int, len(seqs))
		n := len(seqs)
		for {
			p := make([]T, n)
			for i, idx := range indices {
				p[i] = materialized[i][idx]
			}
			if !yield(p) {
				return
			}
			i := n - 1
			for i >= 0 {
				indices[i]++
				if indices[i] < len(materialized[i]) {
					break
				}
				indices[i] = 0
				i--
			}
			if i < 0 {
				return
			}
		}
	}
}

//nolint:prealloc // size unknown upfront for sequences
func toSlice[T any](seq xiter.Seq[T]) []T {
	var result []T
	for v := range seq {
		result = append(result, v)
	}
	return result
}
