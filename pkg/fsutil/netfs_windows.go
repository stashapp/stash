//go:build windows

package fsutil

import (
	"path/filepath"

	"golang.org/x/sys/windows"
)

// IsNetworkFS reports whether path resides on a network filesystem.
// On Windows, it uses GetDriveType to detect DRIVE_REMOTE volumes.
func IsNetworkFS(path string) (bool, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return false, err
	}

	// GetDriveType wants a root path like "\\server\share\" or "C:\"
	vol := filepath.VolumeName(abs)
	if vol == "" {
		return false, nil
	}

	root := vol + `\`
	rootPtr, err := windows.UTF16PtrFromString(root)
	if err != nil {
		return false, err
	}

	driveType := windows.GetDriveType(rootPtr)
	return driveType == windows.DRIVE_REMOTE, nil
}
