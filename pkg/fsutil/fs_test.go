package fsutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsFsPathCaseSensitive_UnicodeByteLength(t *testing.T) {
	// Ⱥ (U+023A) is 2 bytes in UTF-8
	// Its lowercase ⱥ (U+2C65) is 3 bytes in UTF-8

	dir := t.TempDir()
	makeDir := func(path string) {
		// Create the directory so os.Stat succeeds
		if err := os.Mkdir(path, 0755); err != nil {
			t.Fatal(err)
		}
	}

	path := filepath.Join(dir, "Ⱥtest")
	makeDir(path)

	// ensure the test does not panic due to byte length differences in the case flipped path
	_, err := IsFsPathCaseSensitive(path)
	if err != nil {
		t.Fatal(err)
	}

	// no guarantee about case sensitivity of the fs running the tests,
	// so we just want to ensure the function works and does not panic
	// assert.True(t, r, "expected fs to be case sensitive")

	// test regular ASCII paths still work
	path2 := filepath.Join(dir, "Test")
	makeDir(path2)

	_, err = IsFsPathCaseSensitive(path2)
	if err != nil {
		t.Fatal(err)
	}

	// assert.True(t, r, "expected fs to be case sensitive")
}
