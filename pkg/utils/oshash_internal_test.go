package utils

import (
	"testing"
)

// Note that the public API returns "" instead.
func TestOshashEmpty(t *testing.T) {
	var size int64
	head := make([]byte, chunkSize)
	tail := make([]byte, chunkSize)
	want := "0000000000000000"
	got, err := oshash(size, head, tail)
	if err != nil {
		t.Errorf("TestOshashEmpty: Error from oshash: %v", err)
	}
	if got != want {
		t.Errorf("TestOshashEmpty: oshash(0, 0, 0) = %q; want %q", got, want)
	}
}

// As oshash sums byte values, causing collisions is trivial.
func TestOshashCollisions(t *testing.T) {
	buf1 := []byte("this is dumb")
	buf2 := []byte("dumb is this")
	size := int64(len(buf1))
	head := make([]byte, chunkSize)

	tail1 := make([]byte, chunkSize)
	copy(tail1[len(tail1)-len(buf1):], buf1)
	hash1, err := oshash(size, head, tail1)
	if err != nil {
		t.Errorf("TestOshashCollisions: Error from oshash: %v", err)
	}

	tail2 := make([]byte, chunkSize)
	copy(tail2[len(tail2)-len(buf2):], buf2)
	hash2, err := oshash(size, head, tail2)
	if err != nil {
		t.Errorf("TestOshashCollisions: Error from oshash: %v", err)
	}

	if hash1 != hash2 {
		t.Errorf("TestOshashCollisions: oshash(n, k, ... %v) =! oshash(n, k, ... %v)", buf1, buf2)
	}
}
