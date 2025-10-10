package xsync

import "sync"

type SyncPool[T any] struct {
	pool sync.Pool
}

// NewSyncPool creates a new SyncPool with specified init function
func NewSyncPool[T any](new func() T) *SyncPool[T] {
	sp := &SyncPool[T]{
		pool: sync.Pool{New: func() interface{} {
			return new()
		}},
	}
	return sp
}

// Get wraps sync.Pool.Get.
func (s *SyncPool[T]) Get() T {
	return s.pool.Get().(T)
}

// Put wraps sync.Pool.Put.
func (s *SyncPool[T]) Put(x T) {
	s.pool.Put(x)
}
