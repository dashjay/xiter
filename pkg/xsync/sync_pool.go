package xsync

import "sync"

type SyncPool[T any] struct {
	New  func() T
	pool sync.Pool
	once sync.Once
}

// NewSyncPool creates a new SyncPool with specified init function
func NewSyncPool[T any](new func() T) *SyncPool[T] {
	sp := &SyncPool[T]{
		New: new,
	}
	sp.init()
	return sp
}

// init initializes the SyncPool
func (s *SyncPool[T]) init() {
	s.once.Do(func() {
		s.pool.New = func() interface{} {
			return s.New()
		}
	})
}

// Get wraps sync.Pool.Get.
func (s *SyncPool[T]) Get() T {
	s.init()
	return s.pool.Get().(T)
}

// Put wraps sync.Pool.Put.
func (s *SyncPool[T]) Put(x T) {
	s.init()
	s.pool.Put(x)
}
