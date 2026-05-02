//go:build go1.23

// Package collect provides lazy collection and windowing operations for iter.Seq.
//
// Functions in this package materialize or buffer elements as needed, then
// yield results as a sequence for downstream processing.
package collect

import "github.com/dashjay/xiter/xiter"

// Group represents a group of values sharing the same key, produced by
// GroupBy or SortedGroupBy.
type Group[K comparable, V any] struct {
	Key   K
	Items []V
}

// GroupBy groups elements from seq by the key returned by keyFn.
// The entire sequence is materialized to build groups, then the groups
// are yielded as a Seq[Group[K, V]].
//
// Example:
//
//	words := xiter.FromSlice([]string{"apple", "banana", "avocado", "blueberry"})
//	groups := collect.GroupBy(words, func(s string) string {
//		return string(s[0]) // group by first letter
//	})
//	for g := range groups {
//		fmt.Printf("%s: %v\n", g.Key, g.Items)
//	}
func GroupBy[T any, K comparable](seq xiter.Seq[T], keyFn func(T) K) xiter.Seq[Group[K, T]] {
	return func(yield func(Group[K, T]) bool) {
		groups := make(map[K][]T)
		var keys []K
		for v := range seq {
			k := keyFn(v)
			if _, ok := groups[k]; !ok {
				keys = append(keys, k)
			}
			groups[k] = append(groups[k], v)
		}
		for _, k := range keys {
			if !yield(Group[K, T]{Key: k, Items: groups[k]}) {
				return
			}
		}
	}
}

// SortedGroupBy returns groups from a pre-sorted sequence.
// The input seq must already be sorted by the key returned by keyFn.
// Unlike GroupBy, this function uses O(1) memory per group and yields
// each group as soon as the key changes.
//
// If the input is not sorted by key, groups will be split incorrectly
// (the same key may appear in multiple groups).
//
// Example:
//
//	words := xiter.FromSlice([]string{"apple", "avocado", "banana", "blueberry"})
//	groups := collect.SortedGroupBy(words, func(s string) string {
//		return string(s[0]) // input must be sorted by first letter
//	})
func SortedGroupBy[T any, K comparable](seq xiter.Seq[T], keyFn func(T) K) xiter.Seq[Group[K, T]] {
	return func(yield func(Group[K, T]) bool) {
		var current *Group[K, T]
		for v := range seq {
			k := keyFn(v)
			if current == nil || current.Key != k {
				if current != nil {
					if !yield(*current) {
						return
					}
				}
				current = &Group[K, T]{Key: k}
			}
			current.Items = append(current.Items, v)
		}
		if current != nil && len(current.Items) > 0 {
			yield(*current)
		}
	}
}

// Window yields sliding windows of n elements from seq.
// Each window is a newly allocated slice, overlapping by n-1 elements.
// If seq has fewer than n elements, no windows are yielded.
//
// Example:
//
//	seq := xiter.FromSlice([]int{1, 2, 3, 4, 5, 6})
//	for w := range collect.Window(seq, 3) {
//		fmt.Println(w) // [1 2 3] [2 3 4] [3 4 5] [4 5 6]
//	}
func Window[T any](seq xiter.Seq[T], n int) xiter.Seq[[]T] {
	return func(yield func([]T) bool) {
		if n <= 0 {
			return
		}
		window := make([]T, 0, n)
		for v := range seq {
			window = append(window, v)
			if len(window) == n {
				w := make([]T, n)
				copy(w, window)
				if !yield(w) {
					return
				}
				window = window[1:]
			}
		}
	}
}
