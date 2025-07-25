package xmap

import (
	"github.com/dashjay/xiter/pkg/xiter"
)

// CoalesceMaps combines multiple maps into a single map. When duplicate keys are encountered,
// the value from the rightmost (last) map in the input slice takes precedence.
// It iterates through the input maps, converts them to sequences of key-value pairs,
// concatenates these sequences, and then converts the combined sequence back into a new map.
//
// Parameters:
//
//	maps ...M: A variadic slice of maps of type M. Each map must have comparable keys K and values V.
//
// Returns:
//
//	M: A new map containing all key-value pairs from the input maps, with later maps overriding
//	   values for duplicate keys.
//
// Example:
//
//	map1 := map[string]int{"a": 1, "b": 2}
//	map2 := map[string]int{"b": 3, "c": 4}
//	map3 := map[string]int{"d": 5}
//
//	// CoalesceMaps will combine map1, map2, and map3.
//	// For key "b", the value 3 from map2 will override 2 from map1.
//	result := CoalesceMaps(map1, map2, map3)
//	// result will be map[string]int{"a": 1, "b": 3, "c": 4, "d": 5}
func CoalesceMaps[M ~map[K]V, K comparable, V any](maps ...M) M {
	seqs := make([]xiter.Seq2[K, V], 0, len(maps))
	for _, m := range maps {
		seqs = append(seqs, xiter.FromMapKeyAndValues(m))
	}
	return xiter.ToMap(xiter.Concat2(seqs...))
}

// Filter filters the map by the given function.
// Example:
//
//	m := map[string]int{"a": 1, "b": 2, "c": 3}
//	fn := func(k string, v int) bool {
//		return v > 1
//	}
//	result := Filter(m, fn)
//	// result will be map[string]int{"b": 2, "c": 3}
func Filter[M ~map[K]V, K comparable, V any](in M, fn func(K, V) bool) M {
	return xiter.ToMap(xiter.Filter2(fn, xiter.FromMapKeyAndValues(in)))
}
