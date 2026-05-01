//go:build go1.23
// +build go1.23

package xsync

// Clear deletes all the entries, resulting in an empty Map.
// only available in go1.23
func (s *SyncMap[K, V]) Clear() {
	s.m.Clear()
}
