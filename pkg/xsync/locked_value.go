package xsync

import "sync"

type noCopy struct {
}

func (n noCopy) Lock() {}

func (n noCopy) Unlock() {}

// LockedValue is a wrapper wrapping a value protect by a mutex
type LockedValue[T any] struct {
	value T
	mu    sync.Mutex
	_     noCopy
}

// NewLockedValue returns a new LockedValue
func NewLockedValue[T any](value T) *LockedValue[T] {
	return &LockedValue[T]{value: value}
}

// LockCB is a shortcut for l.mu.Lock() and defer l.mu.Unlock()
func (l *LockedValue[T]) LockCB(cb func(T)) {
	l.mu.Lock()
	cb(l.value)
	l.mu.Unlock()
}

// Lock and get the value
func (l *LockedValue[T]) Lock() T {
	l.mu.Lock()
	return l.value
}

// Unlock then the value is unprotected
func (l *LockedValue[T]) Unlock() {
	l.mu.Unlock()
}

// TryLock return true with value if lock successfully, return false with zero value if lock failed
func (l *LockedValue[T]) TryLock() (val T, locked bool) {
	locked = l.mu.TryLock()
	if locked {
		val = l.value
	}
	return
}

// SetValue can modify the underlying value with protection
func (l *LockedValue[T]) SetValue(value T) {
	l.mu.Lock()
	l.value = value
	l.mu.Unlock()
}

// RWLockedValue is a wrapper wrapping a value protect by a RWMutex
type RWLockedValue[T any] struct {
	value T
	mu    sync.RWMutex
	_     noCopy
}

// NewRWLockedValue returns a new RWLockedValue
func NewRWLockedValue[T any](value T) *RWLockedValue[T] {
	return &RWLockedValue[T]{value: value}
}

// Lock and get the value
func (l *RWLockedValue[T]) Lock() T {
	l.mu.Lock()
	return l.value
}

// Unlock then the value is unprotected
func (l *RWLockedValue[T]) Unlock() {
	l.mu.Unlock()
}

// LockCB is a shortcut for l.mu.RLock() and defer l.mu.RUnlock()
func (l *RWLockedValue[T]) LockCB(cb func(T)) {
	l.mu.Lock()
	cb(l.value)
	l.mu.Unlock()
}

// RLockCB is a shortcut for l.mu.RLock() and defer l.mu.RUnlock()
func (l *RWLockedValue[T]) RLockCB(cb func(T)) {
	l.mu.RLock()
	cb(l.value)
	l.mu.RUnlock()
}

// TryLock return true with value if lock successfully, return false with zero value if lock failed
func (l *RWLockedValue[T]) TryLock() (val T, locked bool) {
	locked = l.mu.TryLock()
	if locked {
		val = l.value
	}
	return
}

// RLock and get the value
// Gentleman's agreement: RLock means you should not modify the value,
// We cannot force a declaration that the return value cannot be modified
func (l *RWLockedValue[T]) RLock() T {
	l.mu.RLock()
	return l.value
}

// RUnlock then the value is unprotected
func (l *RWLockedValue[T]) RUnlock() {
	l.mu.RUnlock()
}

// TryRLock return true with value if lock successfully, return false with zero value if lock failed
func (l *RWLockedValue[T]) TryRLock() (val T, locked bool) {
	locked = l.mu.TryRLock()
	if locked {
		val = l.value
	}
	return
}

// SetValue can modify the underlying value with protection
func (l *RWLockedValue[T]) SetValue(value T) {
	l.mu.Lock()
	l.value = value
	l.mu.Unlock()
}
