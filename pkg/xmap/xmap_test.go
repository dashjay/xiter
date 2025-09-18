package xmap_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/dashjay/xiter/pkg/union"
	"github.com/dashjay/xiter/pkg/xiter"
	"github.com/dashjay/xiter/pkg/xmap"
	"github.com/stretchr/testify/assert"
)

func _map(start, end int) (m map[int]string) {
	m = make(map[int]string)
	for i := start; i < end; i++ {
		m[i] = fmt.Sprintf("%09d", i)
	}
	return m
}

func TestMap(t *testing.T) {
	t.Run("clone", func(t *testing.T) {
		assert.True(t, xmap.Equal(_map(0, 100), xmap.Clone(_map(0, 100))))
		assert.True(t, xmap.EqualFunc(
			_map(0, 100),
			xmap.Clone(_map(0, 100)),
			func(a string, b string) bool {
				return a == b
			},
		))
	})

	t.Run("copy", func(t *testing.T) {
		m := make(map[int]string)
		xmap.Copy(m, _map(0, 100))
		assert.True(t, xmap.Equal(_map(0, 100), m))
	})

	t.Run("keys values", func(t *testing.T) {
		keys := xmap.Keys(_map(0, 100))
		sort.Ints(keys)

		values := xmap.Values(_map(0, 100))
		sort.Strings(values)
		for i := 0; i < 100; i++ {
			assert.Equal(t, i, keys[i])
			assert.Equal(t, fmt.Sprintf("%09d", i), values[i])
		}
	})

	t.Run("key value union", func(t *testing.T) {
		keyValues := xmap.ToUnionSlice(_map(0, 101))
		result := xiter.Sum(xiter.Map(func(t union.U2[int, string]) int {
			return t.T1
		}, xiter.FromSlice(keyValues)))
		assert.Equal(t, 5050, result)
	})

	t.Run("map equal", func(t *testing.T) {
		assert.False(t, xmap.Equal(_map(0, 10), _map(0, 11)))
		assert.False(t, xmap.EqualFunc(_map(0, 10), _map(0, 11), func(a string, b string) bool {
			return a == b
		}))

		aMap := _map(0, 10)
		bMap := xiter.ToMap(xiter.Concat2(xiter.FromMapKeyAndValues(_map(0, 10)), xiter.FromMapKeyAndValues(map[int]string{1: "11"})))
		assert.False(t, xmap.Equal(aMap, bMap))
		assert.False(t, xmap.EqualFunc(aMap, bMap, func(a string, b string) bool {
			return a == b
		}))
	})

	t.Run("coalesce maps", func(t *testing.T) {
		var maps []map[int]string
		for i := 0; i < 100; i += 10 {
			maps = append(maps, _map(i, i+10))
		}
		assert.True(t, xmap.Equal(xmap.CoalesceMaps(maps...), _map(0, 100)))
	})

	t.Run("filter map", func(t *testing.T) {
		m := _map(0, 100)
		fn := func(k int, v string) bool {
			return k > 50
		}
		result := xmap.Filter(m, fn)
		assert.True(t, xmap.Equal(result, _map(51, 100)))
	})

	t.Run("map values", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2, "c": 3}
		fn := func(k string, v int) string {
			return fmt.Sprintf("value_%d", v)
		}
		var result = xmap.MapValues(m, fn)
		expected := map[string]string{"a": "value_1", "b": "value_2", "c": "value_3"}
		assert.Equal(t, expected, result)

		// Test with key usage in transformation
		fnWithKey := func(k string, v int) string {
			return fmt.Sprintf("%s_%d", k, v)
		}
		resultWithKey := xmap.MapValues(m, fnWithKey)
		expectedWithKey := map[string]string{"a": "a_1", "b": "b_2", "c": "c_3"}
		assert.Equal(t, expectedWithKey, resultWithKey)

		// Test with empty map
		emptyMap := map[string]int{}
		emptyResult := xmap.MapValues(emptyMap, fn)
		assert.Equal(t, map[string]string{}, emptyResult)
	})

	t.Run("map keys", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2, "c": 3}
		fn := func(k string, v int) string {
			return k + "_key"
		}
		var result = xmap.MapKeys(m, fn)
		expected := map[string]int{"a_key": 1, "b_key": 2, "c_key": 3}
		assert.Equal(t, expected, result)

		// Test with value usage in transformation
		fnWithValue := func(k string, v int) string {
			return fmt.Sprintf("%s_%d", k, v)
		}
		resultWithValue := xmap.MapKeys(m, fnWithValue)
		expectedWithValue := map[string]int{"a_1": 1, "b_2": 2, "c_3": 3}
		assert.Equal(t, expectedWithValue, resultWithValue)

		// Test with empty map
		emptyMap := map[string]int{}
		emptyResult := xmap.MapKeys(emptyMap, fn)
		assert.Equal(t, map[string]int{}, emptyResult)
	})

	t.Run("to xsync-map", func(t *testing.T) {
		m := xmap.ToXSyncMap(_map(0, 100))
		assert.Equal(t, 100, m.Len())
	})

	t.Run("xmap find", func(t *testing.T) {
		k, _, ok := xmap.Find(_map(0, 100), func(k int, v string) bool {
			return k == 50
		})
		assert.Equal(t, 50, k)
		assert.True(t, ok)

		res := xmap.FindO(_map(0, 100), func(k int, v string) bool {
			return k == 50
		})
		assert.Equal(t, 50, res.Must().T1)
		assert.True(t, res.Ok())

		k, _, ok = xmap.FindKey(_map(0, 100), 50)
		assert.Equal(t, 50, k)
		assert.True(t, ok)

		res = xmap.FindKeyO(_map(0, 100), 50)
		assert.Equal(t, 50, res.Must().T1)
		assert.True(t, res.Ok())

		// can not find
		k, _, ok = xmap.Find(_map(0, 100), func(k int, v string) bool {
			return k == 100
		})
		assert.Equal(t, 0, k)
		assert.False(t, ok)

		res = xmap.FindO(_map(0, 100), func(k int, v string) bool {
			return k == 100
		})
		assert.False(t, res.Ok())

		k, _, ok = xmap.FindKey(_map(0, 100), 100)
		assert.Equal(t, 0, k)
		assert.False(t, ok)

		res = xmap.FindKeyO(_map(0, 100), 100)
		assert.False(t, res.Ok())
	})

	t.Run("update", func(t *testing.T) {
		// Test updating existing key
		m := map[string]int{"a": 1, "b": 2, "c": 3}
		old, replaced := xmap.Update(m, "b", 20)
		assert.Equal(t, 2, old)
		assert.True(t, replaced)
		assert.Equal(t, 20, m["b"])
		assert.Equal(t, 3, len(m)) // Map size should not change

		// Test updating non-existing key
		old, replaced = xmap.Update(m, "d", 4)
		assert.Equal(t, 0, old)
		assert.False(t, replaced)
		assert.Equal(t, 4, m["d"])
		assert.Equal(t, 4, len(m)) // Map size should increase

		// Test updating empty map
		emptyMap := map[string]int{}
		old, replaced = xmap.Update(emptyMap, "key", 100)
		assert.Equal(t, 0, old)
		assert.False(t, replaced)
		assert.Equal(t, 100, emptyMap["key"])
		assert.Equal(t, 1, len(emptyMap))
	})

	t.Run("update if", func(t *testing.T) {
		// Test condition is true - should update
		m := map[string]int{"a": 1, "b": 2, "c": 3}
		old, replaced := xmap.UpdateIf(m, "b", 20, func(k string, v int) bool { return v > 1 })
		assert.Equal(t, 2, old)
		assert.True(t, replaced)
		assert.Equal(t, 20, m["b"])

		// Test condition is false - should not update
		old, replaced = xmap.UpdateIf(m, "b", 200, func(k string, v int) bool { return v > 50 })
		assert.Equal(t, 0, old)
		assert.False(t, replaced)
		assert.Equal(t, 20, m["b"]) // Value should remain unchanged

		// Test key doesn't exist - should not update
		old, replaced = xmap.UpdateIf(m, "d", 4, func(k string, v int) bool { return v > 0 })
		assert.Equal(t, 0, old)
		assert.False(t, replaced)
		_, exists := m["d"]
		assert.False(t, exists) // Key should not be added

		// Test complex condition
		old, replaced = xmap.UpdateIf(m, "a", 10, func(k string, v int) bool { return k == "a" && v < 5 })
		assert.Equal(t, 1, old)
		assert.True(t, replaced)
		assert.Equal(t, 10, m["a"])
	})

	t.Run("delete", func(t *testing.T) {
		// Test deleting existing key
		m := map[string]int{"a": 1, "b": 2, "c": 3}
		old, deleted := xmap.Delete(m, "b")
		assert.Equal(t, 2, old)
		assert.True(t, deleted)
		_, exists := m["b"]
		assert.False(t, exists)
		assert.Equal(t, 2, len(m))

		// Test deleting non-existing key
		old, deleted = xmap.Delete(m, "d")
		assert.Equal(t, 0, old)
		assert.False(t, deleted)
		assert.Equal(t, 2, len(m))

		// Test deleting from empty map
		emptyMap := map[string]int{}
		old, deleted = xmap.Delete(emptyMap, "key")
		assert.Equal(t, 0, old)
		assert.False(t, deleted)
		assert.Equal(t, 0, len(emptyMap))

		// Test deleting all elements
		m2 := map[string]int{"x": 10, "y": 20}
		xmap.Delete(m2, "x")
		xmap.Delete(m2, "y")
		assert.Equal(t, 0, len(m2))
	})

	t.Run("delete if", func(t *testing.T) {
		// Test condition is true - should delete
		m := map[string]int{"a": 1, "b": 2, "c": 3}
		old, deleted := xmap.DeleteIf(m, "b", func(k string, v int) bool { return v > 1 })
		assert.Equal(t, 2, old)
		assert.True(t, deleted)
		_, exists := m["b"]
		assert.False(t, exists)
		assert.Equal(t, 2, len(m))

		// Test condition is false - should not delete
		old, deleted = xmap.DeleteIf(m, "a", func(k string, v int) bool { return v > 10 })
		assert.Equal(t, 0, old)
		assert.False(t, deleted)
		assert.Equal(t, 1, m["a"]) // Value should remain unchanged
		assert.Equal(t, 2, len(m))

		// Test key doesn't exist - should not delete
		old, deleted = xmap.DeleteIf(m, "d", func(k string, v int) bool { return v > 0 })
		assert.Equal(t, 0, old)
		assert.False(t, deleted)
		assert.Equal(t, 2, len(m))

		// Test complex condition
		m2 := map[string]int{"x": 10, "y": 5, "z": 15}
		old, deleted = xmap.DeleteIf(m2, "y", func(k string, v int) bool { return k == "y" && v < 10 })
		assert.Equal(t, 5, old)
		assert.True(t, deleted)
		assert.Equal(t, 2, len(m2))
	})

	t.Run("max key", func(t *testing.T) {
		// Test with string keys
		m := map[string]int{"b": 2, "a": 1, "c": 3}
		result := xmap.MaxKey(m)
		assert.True(t, result.Ok())
		assert.Equal(t, "c", result.Must().T1)
		assert.Equal(t, 3, result.Must().T2)

		// Test with int keys
		m2 := map[int]string{3: "three", 1: "one", 5: "five"}
		result2 := xmap.MaxKey(m2)
		assert.True(t, result2.Ok())
		assert.Equal(t, 5, result2.Must().T1)
		assert.Equal(t, "five", result2.Must().T2)

		// Test empty map
		emptyMap := map[string]int{}
		result3 := xmap.MaxKey(emptyMap)
		assert.False(t, result3.Ok())

		// Test single element
		singleMap := map[string]int{"only": 42}
		result4 := xmap.MaxKey(singleMap)
		assert.True(t, result4.Ok())
		assert.Equal(t, "only", result4.Must().T1)
		assert.Equal(t, 42, result4.Must().T2)
	})

	t.Run("min key", func(t *testing.T) {
		// Test with string keys
		m := map[string]int{"b": 2, "a": 1, "c": 3}
		result := xmap.MinKey(m)
		assert.True(t, result.Ok())
		assert.Equal(t, "a", result.Must().T1)
		assert.Equal(t, 1, result.Must().T2)

		// Test with int keys
		m2 := map[int]string{3: "three", 1: "one", 5: "five"}
		result2 := xmap.MinKey(m2)
		assert.True(t, result2.Ok())
		assert.Equal(t, 1, result2.Must().T1)
		assert.Equal(t, "one", result2.Must().T2)

		// Test empty map
		emptyMap := map[string]int{}
		result3 := xmap.MinKey(emptyMap)
		assert.False(t, result3.Ok())

		// Test single element
		singleMap := map[string]int{"only": 42}
		result4 := xmap.MinKey(singleMap)
		assert.True(t, result4.Ok())
		assert.Equal(t, "only", result4.Must().T1)
		assert.Equal(t, 42, result4.Must().T2)
	})

	t.Run("max value", func(t *testing.T) {
		// Test with int values
		m := map[string]int{"a": 3, "b": 1, "c": 2}
		result := xmap.MaxValue(m)
		assert.True(t, result.Ok())
		assert.Equal(t, "a", result.Must().T1)
		assert.Equal(t, 3, result.Must().T2)

		// Test with string values
		m2 := map[int]string{1: "zebra", 2: "apple", 3: "banana"}
		result2 := xmap.MaxValue(m2)
		assert.True(t, result2.Ok())
		assert.Equal(t, 1, result2.Must().T1)
		assert.Equal(t, "zebra", result2.Must().T2)

		// Test empty map
		emptyMap := map[string]int{}
		result3 := xmap.MaxValue(emptyMap)
		assert.False(t, result3.Ok())

		// Test single element
		singleMap := map[string]int{"only": 42}
		result4 := xmap.MaxValue(singleMap)
		assert.True(t, result4.Ok())
		assert.Equal(t, "only", result4.Must().T1)
		assert.Equal(t, 42, result4.Must().T2)

		// Test with duplicate max values - should return one of them
		dupMap := map[string]int{"x": 5, "y": 5, "z": 3}
		result5 := xmap.MaxValue(dupMap)
		assert.True(t, result5.Ok())
		assert.Equal(t, 5, result5.Must().T2)
		assert.Contains(t, []string{"x", "y"}, result5.Must().T1)
	})

	t.Run("min value", func(t *testing.T) {
		// Test with int values
		m := map[string]int{"a": 3, "b": 1, "c": 2}
		result := xmap.MinValue(m)
		assert.True(t, result.Ok())
		assert.Equal(t, "b", result.Must().T1)
		assert.Equal(t, 1, result.Must().T2)

		// Test with string values
		m2 := map[int]string{1: "zebra", 2: "apple", 3: "banana"}
		result2 := xmap.MinValue(m2)
		assert.True(t, result2.Ok())
		assert.Equal(t, 2, result2.Must().T1)
		assert.Equal(t, "apple", result2.Must().T2)

		// Test empty map
		emptyMap := map[string]int{}
		result3 := xmap.MinValue(emptyMap)
		assert.False(t, result3.Ok())

		// Test single element
		singleMap := map[string]int{"only": 42}
		result4 := xmap.MinValue(singleMap)
		assert.True(t, result4.Ok())
		assert.Equal(t, "only", result4.Must().T1)
		assert.Equal(t, 42, result4.Must().T2)

		// Test with duplicate min values - should return one of them
		dupMap := map[string]int{"x": 1, "y": 1, "z": 3}
		result5 := xmap.MinValue(dupMap)
		assert.True(t, result5.Ok())
		assert.Equal(t, 1, result5.Must().T2)
		assert.Contains(t, []string{"x", "y"}, result5.Must().T1)
	})

}
