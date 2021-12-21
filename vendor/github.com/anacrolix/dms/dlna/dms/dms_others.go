//go:build !linux && !darwin && !windows
// +build !linux,!darwin,!windows

package dms

func isHiddenPath(path string) (bool, error) {
	return false, nil
}

func isReadablePath(path string) (bool, error) {
	return tryToOpenPath(path)
}
