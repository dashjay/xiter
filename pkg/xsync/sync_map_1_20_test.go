//go:build go1.20
// +build go1.20

package xsync_test

import (
	"testing"

	"github.com/dashjay/xiter/pkg/xsync"
	"github.com/stretchr/testify/assert"
)

func TestSyncMap120(t *testing.T) {
	t.Parallel()

	t.Run("simple swap", func(t *testing.T) {
		m := xsync.NewSyncMap[string, int]()
		m.Store("1", 1)

		v, exists := m.Load("1")
		assert.True(t, exists)
		assert.Equal(t, 1, v)

		// swap key exists
		prev, loaded := m.Swap("1", 2)
		assert.True(t, loaded)
		assert.Equal(t, 1, prev)

		v, exists = m.Load("1")
		assert.True(t, exists)
		assert.Equal(t, 2, v)

		// swap key does not exist
		prev, loaded = m.Swap("2", 2)
		assert.False(t, loaded)
	})

	t.Run("simple compare and swap", func(t *testing.T) {
		m := xsync.NewSyncMap[string, int]()
		m.Store("1", 1)

		v, exists := m.Load("1")
		assert.True(t, exists)
		assert.Equal(t, 1, v)

		swapped := m.CompareAndSwap("1", 1, 2)
		assert.True(t, swapped)

		v, exists = m.Load("1")
		assert.True(t, exists)
		assert.Equal(t, 2, v)

		swapped = m.CompareAndSwap("1", 1, 3)
		assert.False(t, swapped)
	})

	t.Run("simple compare and delete", func(t *testing.T) {
		m := xsync.NewSyncMap[string, int]()
		m.Store("1", 1)

		v, exists := m.Load("1")
		assert.True(t, exists)
		assert.Equal(t, 1, v)

		deleted := m.CompareAndDelete("1", 1)
		assert.True(t, deleted)

		v, exists = m.Load("1")
		assert.False(t, exists)
	})
}
