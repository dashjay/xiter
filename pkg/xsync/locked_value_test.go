package xsync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLockedValue(t *testing.T) {
	t.Parallel()

	v := 0
	lv := NewLockedValue[*int](&v)

	lv.LockCB(func(i *int) {
		assert.Equal(t, &v, i)
		assert.Equal(t, v, *i)
	})

	tmp := lv.Lock()
	*tmp += 1
	lv.Unlock()

	val, locked := lv.TryLock()
	assert.True(t, locked)
	assert.Equal(t, v, *val)
	lv.Unlock()

	lv.SetValue(func() *int {
		x := 100
		return &x
	}())

	lv.LockCB(func(i *int) {
		assert.Equal(t, 100, *i)
	})
}

func TestRWLockedValue(t *testing.T) {
	t.Parallel()

	v := 0
	lv := NewRWLockedValue[*int](&v)

	lv.LockCB(func(i *int) {
		assert.Equal(t, &v, i)
		assert.Equal(t, v, *i)
	})

	tmp := lv.Lock()
	*tmp += 1
	lv.Unlock()

	val, locked := lv.TryLock()
	assert.True(t, locked)
	assert.Equal(t, v, *val)
	lv.Unlock()

	lv.SetValue(func() *int {
		x := 100
		return &x
	}())

	lv.LockCB(func(i *int) {
		assert.Equal(t, 100, *i)
	})

	val = lv.RLock()
	assert.Equal(t, 100, *val)
	lv.RUnlock()

	lv.RLockCB(func(i *int) {
		assert.Equal(t, 100, *i)
	})

	val, locked = lv.TryRLock()
	assert.True(t, locked)
	assert.Equal(t, 100, *val)
}
