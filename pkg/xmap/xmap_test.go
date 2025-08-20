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
}
