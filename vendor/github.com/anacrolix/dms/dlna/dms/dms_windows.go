//go:build windows
// +build windows

package dms

import (
	"path/filepath"

	"golang.org/x/sys/windows"
)

const hiddenAttributes = windows.FILE_ATTRIBUTE_HIDDEN | windows.FILE_ATTRIBUTE_SYSTEM

func isHiddenPath(path string) (hidden bool, err error) {
	if path == filepath.VolumeName(path)+"\\" {
		// Volumes always have the "SYSTEM" flag, so do not even test them
		return false, nil
	}
	winPath, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return
	}
	attrs, err := windows.GetFileAttributes(winPath)
	if err != nil {
		return
	}
	if attrs&hiddenAttributes != 0 {
		hidden = true
		return
	}
	return isHiddenPath(filepath.Dir(path))
}

func isReadablePath(path string) (bool, error) {
	return tryToOpenPath(path)
}
