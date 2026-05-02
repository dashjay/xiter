// MIT License
//
// Copyright (c) 2025 DashJay
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package lockedmap provides a concurrency-safe generic map backed by
// sync.RWMutex and a regular Go map.
//
// Unlike xsync.SyncMap (which wraps sync.Map), LockedMap ensures that all
// methods including Len, ToMap, and Range provide consistent snapshots
// under concurrent access.
package lockedmap

import (
	"reflect"
	"sync"
)

// LockedMap is a concurrency-safe generic map backed by sync.RWMutex.
// The zero value is not ready to use; use NewLockedMap to create one.
type LockedMap[K comparable, V any] struct {
	mu sync.RWMutex
	m  map[K]V
}

// NewLockedMap creates a new empty LockedMap.
func NewLockedMap[K comparable, V any]() *LockedMap[K, V] {
	return &LockedMap[K, V]{m: make(map[K]V)}
}

// Load returns the value stored for the key, or false if no value exists.
func (m *LockedMap[K, V]) Load(key K) (value V, ok bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok = m.m[key]
	return
}

// Store sets the value for a key.
func (m *LockedMap[K, V]) Store(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.m[key] = value
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
func (m *LockedMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if v, ok := m.m[key]; ok {
		return v, true
	}
	m.m[key] = value
	return value, false
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
func (m *LockedMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if v, ok := m.m[key]; ok {
		delete(m.m, key)
		return v, true
	}
	return value, false
}

// Delete deletes the value for a key.
func (m *LockedMap[K, V]) Delete(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.m, key)
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
//
// NOTE: Do not call Store, Delete, or other mutating methods on m from within f,
// or a deadlock will occur.
func (m *LockedMap[K, V]) Range(f func(key K, value V) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for k, v := range m.m {
		if !f(k, v) {
			break
		}
	}
}

// Len returns the number of elements in the map. O(1) complexity.
func (m *LockedMap[K, V]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.m)
}

// ToMap returns a copy of the map as a regular map.
// The returned map is a consistent snapshot at the time of the call.
func (m *LockedMap[K, V]) ToMap() map[K]V {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make(map[K]V, len(m.m))
	for k, v := range m.m {
		out[k] = v
	}
	return out
}

// Swap swaps the value for a key and returns the previous value if any.
func (m *LockedMap[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	previous, loaded = m.m[key]
	m.m[key] = value
	return
}

// CompareAndSwap swaps the value for a key if the stored value equals old.
// Comparison uses reflect.DeepEqual.
func (m *LockedMap[K, V]) CompareAndSwap(key K, old, new V) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if v, ok := m.m[key]; ok && reflect.DeepEqual(v, old) {
		m.m[key] = new
		return true
	}
	return false
}

// CompareAndDelete deletes the entry for a key if the stored value equals old.
// Comparison uses reflect.DeepEqual.
func (m *LockedMap[K, V]) CompareAndDelete(key K, old V) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if v, ok := m.m[key]; ok && reflect.DeepEqual(v, old) {
		delete(m.m, key)
		return true
	}
	return false
}

// Clear deletes all entries, resulting in an empty map.
func (m *LockedMap[K, V]) Clear() {
	m.mu.Lock()
	m.m = make(map[K]V)
	m.mu.Unlock()
}
