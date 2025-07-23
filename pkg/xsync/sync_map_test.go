//go:build go1.18
// +build go1.18

package xsync_test

import (
	"strconv"
	"sync"
	"testing"

	"github.com/dashjay/xiter/pkg/xsync"
	"github.com/stretchr/testify/assert"
)

func TestSyncMap(t *testing.T) {
	t.Parallel()

	t.Run("simple store and load", func(t *testing.T) {
		m := xsync.NewSyncMap[string, int]()
		v, exists := m.Load("1")
		assert.False(t, exists)
		m.Store("1", 1)
		v, exists = m.Load("1")
		assert.True(t, exists)
		assert.Equal(t, 1, v)
	})

	t.Run("simple load or store", func(t *testing.T) {
		m := xsync.NewSyncMap[string, int]()
		v, loaded := m.LoadOrStore("1", 1)
		assert.False(t, loaded)
		assert.Equal(t, 1, v)
		v, loaded = m.LoadOrStore("1", 2)
		assert.True(t, loaded)
		assert.Equal(t, 1, v)
	})

	t.Run("simple load and delete", func(t *testing.T) {
		m := xsync.NewSyncMap[string, int]()
		v, loaded := m.LoadAndDelete("1")
		assert.False(t, loaded)
		// default zero value is 0
		assert.Equal(t, 0, v)
		m.Store("1", 1)
		v, loaded = m.LoadAndDelete("1")
		assert.True(t, loaded)
		assert.Equal(t, 1, v)
	})

	t.Run("simple delete", func(t *testing.T) {
		m := xsync.NewSyncMap[string, int]()
		v, exists := m.Load("1")
		assert.False(t, exists)
		m.Store("1", 1)
		v, exists = m.Load("1")
		assert.True(t, exists)
		assert.Equal(t, 1, v)
		m.Delete("1")
		v, exists = m.Load("1")
		assert.False(t, exists)
	})

	t.Run("simple range_len_tomap", func(t *testing.T) {
		m := xsync.NewSyncMap[string, int]()
		const count = 100
		for i := 0; i < count; i++ {
			m.Store(strconv.Itoa(i), i)
		}
		m.Range(func(key string, value int) bool {
			assert.Equal(t, strconv.Itoa(value), key)
			return true
		})
		assert.Equal(t, count, m.Len())
		om := m.ToMap()
		assert.Len(t, om, count)
		for key, value := range om {
			assert.Equal(t, strconv.Itoa(value), key)
		}
	})

	t.Run("concurrent simple test", func(t *testing.T) {
		m := xsync.NewSyncMap[string, int]()
		var wg sync.WaitGroup
		const count = 10_000
		concurrency := 10
		ch := make(chan struct{}, concurrency)
		for i := 0; i < count; i++ {
			ch <- struct{}{}
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				<-ch
				m.Store(strconv.Itoa(idx), idx)
			}(i)
		}
		wg.Wait()
		assert.Equal(t, count, m.Len())
	})
}
