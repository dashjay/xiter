package xsync_test

import (
	"testing"

	"github.com/dashjay/xiter/pkg/xsync"
)

func TestSyncPool(t *testing.T) {
	p := xsync.NewSyncPool[[]byte](func() []byte {
		return make([]byte, 4096)
	})

	for i := 0; i < 1000; i++ {
		v := p.Get()
		p.Put(v)
	}
}
