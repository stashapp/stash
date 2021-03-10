package utils

import (
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
