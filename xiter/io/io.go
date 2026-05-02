//go:build go1.23

// Package io provides lazy I/O operations that produce iter.Seq sequences.
//
// Unlike standard I/O functions that read entire files into memory,
// these functions yield elements on demand, enabling processing of
// large or streaming data with bounded memory.
package io

import (
	"bufio"
	"io"
	"io/fs"
	"os"

	"github.com/dashjay/xiter/xiter"
)

// Lines reads from r line by line, yielding each line as a string.
// The iteration stops when the reader is exhausted or the consumer
// stops iterating.
//
// Example:
//
//	f, _ := os.Open("file.txt")
//	defer f.Close()
//	for line := range iox.Lines(f) {
//		fmt.Println(line)
//	}
func Lines(r io.Reader) xiter.Seq[string] {
	return func(yield func(string) bool) {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			if !yield(scanner.Text()) {
				return
			}
		}
	}
}

// ReadDir returns a Seq of directory entries for the specified directory.
// If the directory cannot be read, an empty sequence is returned.
//
// Example:
//
//	for entry := range iox.ReadDir(".") {
//		fmt.Println(entry.Name())
//	}
func ReadDir(dir string) xiter.Seq[fs.DirEntry] {
	return func(yield func(fs.DirEntry) bool) {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return
		}
		for _, entry := range entries {
			if !yield(entry) {
				return
			}
		}
	}
}

// ReadFileByChunk reads the file at filename in chunks of the specified size.
// Each chunk is a newly allocated slice of bytes.
// If the file cannot be opened, an empty sequence is returned.
//
// Example:
//
//	for chunk := range iox.ReadFileByChunk("large.bin", 4096) {
//		process(chunk)
//	}
//nolint:gosec // G304: file path from caller is intended API
func ReadFileByChunk(filename string, size int) xiter.Seq[[]byte] {
	return func(yield func([]byte) bool) {
		f, err := os.Open(filename)
		if err != nil {
			return
		}
		defer func() { _ = f.Close() }()
		buf := make([]byte, size)
		for {
			n, err := f.Read(buf)
			if n > 0 {
				chunk := make([]byte, n)
				copy(chunk, buf[:n])
				if !yield(chunk) {
					return
				}
			}
			if err != nil {
				break
			}
		}
	}
}
