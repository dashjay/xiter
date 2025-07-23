package gsync_test

import (
	"testing"

	"github.com/dashjay/gog/gsync"
)

func TestSyncPool(t *testing.T) {
	p := gsync.NewSyncPool[[]byte](func() []byte {
		return make([]byte, 4096)
	})

	for i := 0; i < 1000; i++ {
		v := p.Get()
		p.Put(v)
	}
}
