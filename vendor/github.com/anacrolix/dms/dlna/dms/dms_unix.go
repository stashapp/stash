//go:build linux || darwin
// +build linux darwin

package dms

import (
	"strings"

	"golang.org/x/sys/unix"
)

func isHiddenPath(path string) (bool, error) {
	return strings.Contains(path, "/."), nil
}

func isReadablePath(path string) (bool, error) {
	err := unix.Access(path, unix.R_OK)
	switch err {
	case nil:
		return true, nil
	case unix.EACCES:
		return false, nil
	default:
		return false, err
	}
}
