//go:build go1.20
// +build go1.20

package xsync

func (s *SyncMap[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	v, ok := s.m.Swap(key, value)
	if !ok {
		return
	}
	return v.(V), ok
}

func (s *SyncMap[K, V]) CompareAndSwap(key K, old, new V) bool {
	return s.m.CompareAndSwap(key, old, new)
}

func (s *SyncMap[K, V]) CompareAndDelete(key K, old V) bool {
	return s.m.CompareAndDelete(key, old)
}
