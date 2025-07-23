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
}
