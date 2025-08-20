package xmap

import (
	"github.com/dashjay/xiter/pkg/xiter"
	"github.com/dashjay/xiter/pkg/xsync"
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

// MapValues transforms the values of a map using the provided function while keeping keys unchanged.
// This is useful for transforming data structures while preserving the key associations.
//
// Parameters:
//
//	in M: The input map to transform
//	fn func(K, V1) V2: A function that takes a key and its corresponding value, and returns a new value
//
// Returns:
//
//	map[K]V2: A new map with the same keys as the input map but with transformed values
//
// Example:
//
//	m := map[string]int{"a": 1, "b": 2, "c": 3}
//	fn := func(k string, v int) string {
//		return fmt.Sprintf("value_%d", v)
//	}
//	result := MapValues(m, fn)
//	// result will be map[string]string{"a": "value_1", "b": "value_2", "c": "value_3"}
func MapValues[K comparable, V1, V2 any](in map[K]V1, fn func(K, V1) V2) map[K]V2 {
	return xiter.ToMap(xiter.Map2(func(k K, v V1) (K, V2) { return k, fn(k, v) }, xiter.FromMapKeyAndValues(in)))
}

// MapKeys transforms the keys of a map using the provided function while keeping values unchanged.
// This is useful for transforming data structures while preserving the value associations.
//
// Parameters:
//
//	in map[K]V1: The input map to transform
//	fn func(K, V1) K: A function that takes a key and its corresponding value, and returns a new key
//
// Returns:
//
//	map[K]V1: A new map with the same values as the input map but with transformed keys
//
// Example:
//
//	m := map[string]int{"a": 1, "b": 2, "c": 3}
//	fn := func(k string, v int) string {
//		return k + "_key"
//	}
//	result := MapKeys(m, fn)
//	// result will be map[string]int{"a_key": 1, "b_key": 2, "c_key": 3}
func MapKeys[K comparable, V1 any](in map[K]V1, fn func(K, V1) K) map[K]V1 {
	return xiter.ToMap(xiter.Map2(func(k K, v V1) (K, V1) { return fn(k, v), v }, xiter.FromMapKeyAndValues(in)))
}

// ToXSyncMap converts a map to a xsync.SyncMap.
//
// Parameters:
//
//	in map[K]V: The input map to convert
//
// Returns:
//
//	*xsync.SyncMap[K, V]: A new xsync.SyncMap containing the same key-value pairs as the input map
func ToXSyncMap[K comparable, V any](in map[K]V) *xsync.SyncMap[K, V] {
	m := xsync.NewSyncMap[K, V]()
	for k, v := range in {
		m.Store(k, v)
	}
	return m
}
