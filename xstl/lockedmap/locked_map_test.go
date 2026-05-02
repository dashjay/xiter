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

package lockedmap

import (
	"sync"
	"testing"
)

func TestLockedMap(t *testing.T) {
	t.Run("load missing", func(t *testing.T) {
		m := NewLockedMap[string, int]()
		v, ok := m.Load("a")
		if ok || v != 0 {
			t.Errorf("Load missing key: got (%d, %t), want (0, false)", v, ok)
		}
	})

	t.Run("store and load", func(t *testing.T) {
		m := NewLockedMap[string, int]()
		m.Store("a", 1)
		v, ok := m.Load("a")
		if !ok || v != 1 {
			t.Errorf("Store+Load: got (%d, %t), want (1, true)", v, ok)
		}
	})

	t.Run("load or store existing", func(t *testing.T) {
		m := NewLockedMap[string, int]()
		m.Store("a", 1)
		v, loaded := m.LoadOrStore("a", 2)
		if !loaded || v != 1 {
			t.Errorf("LoadOrStore existing: got (%d, %t), want (1, true)", v, loaded)
		}
	})

	t.Run("load or store missing", func(t *testing.T) {
		m := NewLockedMap[string, int]()
		v, loaded := m.LoadOrStore("a", 1)
		if loaded || v != 1 {
			t.Errorf("LoadOrStore missing: got (%d, %t), want (1, false)", v, loaded)
		}
	})

	t.Run("load and delete", func(t *testing.T) {
		m := NewLockedMap[string, int]()
		m.Store("a", 1)
		v, loaded := m.LoadAndDelete("a")
		if !loaded || v != 1 {
			t.Errorf("LoadAndDelete: got (%d, %t), want (1, true)", v, loaded)
		}
		if _, ok := m.Load("a"); ok {
			t.Error("LoadAndDelete: key still exists")
		}
	})

	t.Run("delete", func(t *testing.T) {
		m := NewLockedMap[string, int]()
		m.Store("a", 1)
		m.Delete("a")
		if _, ok := m.Load("a"); ok {
			t.Error("Delete: key still exists")
		}
	})

	t.Run("len", func(t *testing.T) {
		m := NewLockedMap[string, int]()
		if n := m.Len(); n != 0 {
			t.Errorf("Len empty: got %d, want 0", n)
		}
		m.Store("a", 1)
		m.Store("b", 2)
		if n := m.Len(); n != 2 {
			t.Errorf("Len: got %d, want 2", n)
		}
		m.Delete("a")
		if n := m.Len(); n != 1 {
			t.Errorf("Len after delete: got %d, want 1", n)
		}
	})

	t.Run("to map", func(t *testing.T) {
		m := NewLockedMap[string, int]()
		m.Store("a", 1)
		m.Store("b", 2)
		out := m.ToMap()
		if len(out) != 2 {
			t.Fatalf("ToMap len: got %d, want 2", len(out))
		}
		if out["a"] != 1 || out["b"] != 2 {
			t.Errorf("ToMap content: got %v", out)
		}
		// verify it's a copy
		out["c"] = 3
		if m.Len() != 2 {
			t.Error("ToMap: modification of returned map affected original")
		}
	})

	t.Run("range", func(t *testing.T) {
		m := NewLockedMap[int, int]()
		for i := 0; i < 10; i++ {
			m.Store(i, i*10)
		}
		seen := make(map[int]bool)
		m.Range(func(k int, v int) bool {
			if v != k*10 {
				t.Errorf("Range: unexpected value %d for key %d", v, k)
			}
			seen[k] = true
			return true
		})
		if len(seen) != 10 {
			t.Errorf("Range visited %d keys, want 10", len(seen))
		}
	})

	t.Run("range early stop", func(t *testing.T) {
		m := NewLockedMap[int, int]()
		for i := 0; i < 10; i++ {
			m.Store(i, i)
		}
		count := 0
		m.Range(func(k int, v int) bool {
			count++
			return count < 3
		})
		if count != 3 {
			t.Errorf("Range early stop: visited %d, want 3", count)
		}
	})

	t.Run("swap", func(t *testing.T) {
		m := NewLockedMap[string, int]()
		prev, loaded := m.Swap("a", 1)
		if loaded || prev != 0 {
			t.Errorf("Swap missing: got (%d, %t), want (0, false)", prev, loaded)
		}
		v, _ := m.Load("a")
		if v != 1 {
			t.Errorf("Swap missing: stored value = %d, want 1", v)
		}

		prev, loaded = m.Swap("a", 2)
		if !loaded || prev != 1 {
			t.Errorf("Swap existing: got (%d, %t), want (1, true)", prev, loaded)
		}
		v, _ = m.Load("a")
		if v != 2 {
			t.Errorf("Swap existing: stored value = %d, want 2", v)
		}
	})

	t.Run("compare and swap", func(t *testing.T) {
		m := NewLockedMap[string, int]()
		m.Store("a", 1)

		if m.CompareAndSwap("a", 99, 2) {
			t.Error("CompareAndSwap with wrong old value: should have failed")
		}
		if !m.CompareAndSwap("a", 1, 2) {
			t.Error("CompareAndSwap with correct old value: should have succeeded")
		}
		v, _ := m.Load("a")
		if v != 2 {
			t.Errorf("CompareAndSwap: stored value = %d, want 2", v)
		}
	})

	t.Run("compare and delete", func(t *testing.T) {
		m := NewLockedMap[string, int]()
		m.Store("a", 1)

		if m.CompareAndDelete("a", 99) {
			t.Error("CompareAndDelete with wrong old value: should have failed")
		}
		if _, ok := m.Load("a"); !ok {
			t.Error("CompareAndDelete: key should still exist")
		}
		if !m.CompareAndDelete("a", 1) {
			t.Error("CompareAndDelete with correct old value: should have succeeded")
		}
		if _, ok := m.Load("a"); ok {
			t.Error("CompareAndDelete: key should have been deleted")
		}
	})

	t.Run("clear", func(t *testing.T) {
		m := NewLockedMap[string, int]()
		m.Store("a", 1)
		m.Store("b", 2)
		m.Clear()
		if n := m.Len(); n != 0 {
			t.Errorf("Clear: len = %d, want 0", n)
		}
		if _, ok := m.Load("a"); ok {
			t.Error("Clear: key 'a' still exists")
		}
	})

	t.Run("concurrent store and load", func(t *testing.T) {
		m := NewLockedMap[int, int]()
		var wg sync.WaitGroup
		n := 1000
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				m.Store(i, i)
			}(i)
		}
		wg.Wait()

		if m.Len() != n {
			t.Errorf("Concurrent store: len = %d, want %d", m.Len(), n)
		}
	})

	t.Run("concurrent store and len", func(t *testing.T) {
		m := NewLockedMap[int, int]()
		var wg sync.WaitGroup
		// writer
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				m.Store(i, i)
			}
		}()
		// reader calling Len concurrently
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				_ = m.Len()
			}
		}()
		wg.Wait()
	})

	t.Run("concurrent store and to map", func(t *testing.T) {
		m := NewLockedMap[int, int]()
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				m.Store(i, i)
			}
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				_ = m.ToMap()
			}
		}()
		wg.Wait()
	})

	t.Run("concurrent store and range", func(t *testing.T) {
		m := NewLockedMap[int, int]()
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				m.Store(i, i)
			}
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 100; i++ {
				m.Range(func(k int, v int) bool { return true })
			}
		}()
		wg.Wait()
	})
}

func TestLockedMapZeroValue(t *testing.T) {
	// Zero value should panic; the only way to use it is via NewLockedMap
	defer func() {
		if r := recover(); r == nil {
			t.Error("Using zero-value LockedMap should panic (nil map)")
		}
	}()
	var m LockedMap[string, int]
	m.Store("a", 1) // should panic: assignment to entry in nil map
}
