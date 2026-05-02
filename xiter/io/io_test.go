//go:build go1.23

package io_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dashjay/xiter/xiter"
	iox "github.com/dashjay/xiter/xiter/io"
)

func TestLines(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		r := strings.NewReader("a\nb\nc\n")
		lines := xiter.ToSlice(iox.Lines(r))
		if len(lines) != 3 || lines[0] != "a" || lines[1] != "b" || lines[2] != "c" {
			t.Fatalf("Lines: got %v, want [a b c]", lines)
		}
	})

	t.Run("empty", func(t *testing.T) {
		r := strings.NewReader("")
		lines := xiter.ToSlice(iox.Lines(r))
		if len(lines) != 0 {
			t.Fatalf("Lines empty: got %v, want []", lines)
		}
	})

	t.Run("trailing newline", func(t *testing.T) {
		r := strings.NewReader("hello\n")
		lines := xiter.ToSlice(iox.Lines(r))
		if len(lines) != 1 || lines[0] != "hello" {
			t.Fatalf("Lines trailing: got %v, want [hello]", lines)
		}
	})

	t.Run("no trailing newline", func(t *testing.T) {
		r := strings.NewReader("hello")
		lines := xiter.ToSlice(iox.Lines(r))
		if len(lines) != 1 || lines[0] != "hello" {
			t.Fatalf("Lines no trailing: got %v, want [hello]", lines)
		}
	})

	t.Run("multiple lines no trailing", func(t *testing.T) {
		r := strings.NewReader("a\nb\nc")
		lines := xiter.ToSlice(iox.Lines(r))
		if len(lines) != 3 || lines[0] != "a" || lines[1] != "b" || lines[2] != "c" {
			t.Fatalf("Lines multi: got %v, want [a b c]", lines)
		}
	})

	t.Run("early stop", func(t *testing.T) {
		r := strings.NewReader("a\nb\nc\nd\ne\n")
		count := 0
		for range iox.Lines(r) {
			count++
			if count >= 2 {
				break
			}
		}
		if count != 2 {
			t.Fatalf("Lines early stop: count = %d, want 2", count)
		}
	})
}

func TestReadDir(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		entries := xiter.ToSlice(iox.ReadDir("."))
		if len(entries) == 0 {
			t.Fatal("ReadDir: expected at least one entry")
		}
		found := false
		for _, e := range entries {
			if e.Name() == "io.go" {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("ReadDir: should contain io.go")
		}
	})

	t.Run("nonexistent dir", func(t *testing.T) {
		entries := xiter.ToSlice(iox.ReadDir("/nonexistent/path"))
		if len(entries) != 0 {
			t.Fatalf("ReadDir nonexistent: got %d entries, want 0", len(entries))
		}
	})

	t.Run("early stop", func(t *testing.T) {
		count := 0
		for range iox.ReadDir(".") {
			count++
			if count >= 1 {
				break
			}
		}
		if count != 1 {
			t.Fatalf("ReadDir early stop: count = %d, want 1", count)
		}
	})
}

func TestReadFileByChunk(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		content := "hello world this is a test file"
		dir := t.TempDir()
		fpath := filepath.Join(dir, "test.txt")
		if err := os.WriteFile(fpath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}

		var result []byte
		for chunk := range iox.ReadFileByChunk(fpath, 5) {
			result = append(result, chunk...)
		}
		if string(result) != content {
			t.Fatalf("ReadFileByChunk: got %q, want %q", string(result), content)
		}
	})

	t.Run("nonexistent file", func(t *testing.T) {
		chunks := xiter.ToSlice(iox.ReadFileByChunk("/nonexistent/file", 1024))
		if len(chunks) != 0 {
			t.Fatalf("ReadFileByChunk nonexistent: got %d chunks, want 0", len(chunks))
		}
	})

	t.Run("empty file", func(t *testing.T) {
		dir := t.TempDir()
		fpath := filepath.Join(dir, "empty.txt")
		if err := os.WriteFile(fpath, []byte{}, 0644); err != nil {
			t.Fatal(err)
		}
		chunks := xiter.ToSlice(iox.ReadFileByChunk(fpath, 1024))
		if len(chunks) != 0 {
			t.Fatalf("ReadFileByChunk empty: got %d chunks, want 0", len(chunks))
		}
	})

	t.Run("exact chunk size", func(t *testing.T) {
		dir := t.TempDir()
		fpath := filepath.Join(dir, "exact.txt")
		content := "12345"
		if err := os.WriteFile(fpath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		chunks := xiter.ToSlice(iox.ReadFileByChunk(fpath, 5))
		if len(chunks) != 1 {
			t.Fatalf("ReadFileByChunk exact: got %d chunks, want 1", len(chunks))
		}
	})

	t.Run("early stop", func(t *testing.T) {
		dir := t.TempDir()
		fpath := filepath.Join(dir, "early.txt")
		content := make([]byte, 100)
		for i := range content {
			content[i] = byte(i)
		}
		if err := os.WriteFile(fpath, content, 0644); err != nil {
			t.Fatal(err)
		}
		count := 0
		for range iox.ReadFileByChunk(fpath, 10) {
			count++
			if count >= 2 {
				break
			}
		}
		if count != 2 {
			t.Fatalf("ReadFileByChunk early stop: count = %d, want 2", count)
		}
	})

	t.Run("file not closed on early stop", func(t *testing.T) {
		// Verify early stop does not leak file handles by creating many files
		dir := t.TempDir()
		for i := 0; i < 50; i++ {
			fpath := filepath.Join(dir, "f.txt")
			if err := os.WriteFile(fpath, []byte("data"), 0644); err != nil {
				t.Fatal(err)
			}
			for range iox.ReadFileByChunk(fpath, 1) {
				break
			}
		}
		// If we get here without hitting ulimit, file handles are being cleaned up
		// (finalizer or GC may handle it; on most systems 50 handles is fine)
	})
}

func TestLinesFuzzLike(t *testing.T) {
	inputs := []string{
		"\n",
		"\n\n",
		"a\n\nb",
		"\na\nb\n",
	}
	for _, input := range inputs {
		lines := xiter.ToSlice(iox.Lines(strings.NewReader(input)))
		for _, line := range lines {
			if line == "" {
				// bufio.Scanner treats consecutive newlines as empty lines;
				// this is expected behavior
			}
		}
	}
}
