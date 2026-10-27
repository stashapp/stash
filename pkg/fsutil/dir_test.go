package fsutil

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIsPathInDirNormalization verifies that containment holds when the dir and
// the path use different Unicode normalization forms, as happens on macOS where
// the filesystem may report NFD while the configured/stored path is NFC.
func TestIsPathInDirNormalization(t *testing.T) {
	dirNFC := filepath.Join("/library", gaNFC)
	fileNFD := filepath.Join("/library", gaNFD, "video.mp4")

	// Only macOS treats NFC and NFD as equivalent; elsewhere they differ
	// byte-wise and the path is considered outside the dir.
	want := runtime.GOOS == "darwin"
	assert.Equal(t, want, IsPathInDir(dirNFC, fileNFD))

	// An exact (same-form) match must always hold.
	assert.True(t, IsPathInDir(filepath.Join("/library", gaNFC), filepath.Join("/library", gaNFC, "video.mp4")))
}

func TestIsPathInDir(t *testing.T) {
	type test struct {
		dir         string
		pathToCheck string
		expected    bool
	}

	const parentDirName = "parentDir"
	const subDirName = "subDir"
	const filename = "filename"
	subDir := filepath.Join(parentDirName, subDirName)
	fileInSubDir := filepath.Join(subDir, filename)
	fileInParentDir := filepath.Join(parentDirName, filename)
	subSubSubDir := filepath.Join(parentDirName, subDirName, subDirName, subDirName)

	tests := []test{
		{dir: parentDirName, pathToCheck: subDir, expected: true},
		{dir: subDir, pathToCheck: subDir, expected: true},
		{dir: subDir, pathToCheck: parentDirName, expected: false},
		{dir: subDir, pathToCheck: fileInSubDir, expected: true},
		{dir: parentDirName, pathToCheck: fileInSubDir, expected: true},
		{dir: subDir, pathToCheck: fileInParentDir, expected: false},
		{dir: parentDirName, pathToCheck: fileInParentDir, expected: true},
		{dir: parentDirName, pathToCheck: filename, expected: false},
		{dir: parentDirName, pathToCheck: subSubSubDir, expected: true},
		{dir: subSubSubDir, pathToCheck: parentDirName, expected: false},
	}

	assert := assert.New(t)
	for i, tc := range tests {
		result := IsPathInDir(tc.dir, tc.pathToCheck)
		assert.Equal(tc.expected, result, "[%d] expected: %t for dir: %s; pathToCheck: %s", i, tc.expected, tc.dir, tc.pathToCheck)
	}
}

func TestDirExists(t *testing.T) {
	type test struct {
		dir      string
		expected bool
	}

	const st = "stash_tmp"

	tmp := os.TempDir()
	tmpDir, err := os.MkdirTemp(tmp, st) // create a tmp dir in the system's tmp folder
	if err == nil {
		defer os.RemoveAll(tmpDir)

		tmpFile, err := os.CreateTemp(tmpDir, st)
		if err != nil {
			return
		}
		tmpFile.Close()

		tests := []test{
			{dir: tmpDir, expected: true},                     // exists
			{dir: tmpFile.Name(), expected: false},            // not a directory
			{dir: filepath.Join(tmpDir, st), expected: false}, // doesn't exist
			{dir: "\000x", expected: false},                   // stat error  \000 (ASCII: NUL) is an invalid character in unix,ntfs file names.
		}

		assert := assert.New(t)

		for i, tc := range tests {
			result, _ := DirExists(tc.dir)
			assert.Equal(tc.expected, result, "[%d] expected: %t for dir: %s;", i, tc.expected, tc.dir)
		}
	}
}
