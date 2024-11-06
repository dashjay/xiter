//go:build go1.18
// +build go1.18

package gsync

import "sync"

// SyncMap is a wrapper for sync.Map.
type SyncMap[K comparable, V any] struct {
	m sync.Map
}

// NewSyncMap creates a new SyncMap.
func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{}
}

// Load wraps sync.Map.Load.
func (s *SyncMap[K, V]) Load(key K) (value V, ok bool) {
	v, ok := s.m.Load(key)
	if !ok {
		return
	}
	return v.(V), ok
}

// Store wraps sync.Map.Store.
func (s *SyncMap[K, V]) Store(key K, value V) {
	s.m.Store(key, value)
}

// LoadOrStore wraps sync.Map.LoadOrStore.
func (s *SyncMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	v, ok := s.m.LoadOrStore(key, value)
	return v.(V), ok
}

// LoadAndDelete wraps sync.Map.LoadAndDelete.
func (s *SyncMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	v, ok := s.m.LoadAndDelete(key)
	if !ok {
		return
	}
	return v.(V), ok
}

// Delete wraps sync.Map.Delete.
func (s *SyncMap[K, V]) Delete(key K) {
	s.m.Delete(key)
}

// Range wraps sync.Map.Range.
func (s *SyncMap[K, V]) Range(f func(key K, value V) bool) {
	s.m.Range(func(key, value any) bool {
		return f(key.(K), value.(V))
	})
}

// Len returns the number of elements in the map.
// The complexity is O(n).
// Not provided in stdlib but by our own
func (s *SyncMap[K, V]) Len() int {
	l := 0
	s.m.Range(func(key, value any) bool {
		l++
		return true
	})
	return l
}

// ToMap returns a copy of the map as a regular map.
func (s *SyncMap[K, V]) ToMap() map[K]V {
	out := make(map[K]V)
	s.Range(func(key K, value V) bool {
		out[key] = value
		return true
	})
	return out
}
