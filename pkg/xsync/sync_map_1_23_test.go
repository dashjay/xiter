//go:build go1.23
// +build go1.23

package gsync_test

import (
	"strconv"
	"testing"

	"github.com/dashjay/gog/gsync"
	"github.com/stretchr/testify/assert"
)

func TestClear(t *testing.T) {
	t.Run("simple clear", func(t *testing.T) {
		m := gsync.NewSyncMap[string, int]()
		const count = 100
		for i := 0; i < count; i++ {
			m.Store(strconv.Itoa(i), i)
		}
		m.Range(func(key string, value int) bool {
			assert.Equal(t, strconv.Itoa(value), key)
			return true
		})
		m.Clear()
		assert.Equal(t, 0, m.Len())
	})
}
