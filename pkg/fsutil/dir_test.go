package fsutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
