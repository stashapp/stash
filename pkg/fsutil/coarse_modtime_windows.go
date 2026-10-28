//go:build windows

package fsutil

import (
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows"
)

// CoarseModtimeFilesystem reports whether path is on FAT, exFAT, or another
// filesystem where directory modification time is too coarse or unreliable for
// Stash's folder modtime scan optimization.
func CoarseModtimeFilesystem(path string) (bool, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return false, err
	}

	vol := filepath.VolumeName(abs)
	if vol == "" {
		return false, nil
	}

	root := vol
	if !strings.HasSuffix(root, `\`) {
		root += `\`
	}

	rootUTF16, err := windows.UTF16PtrFromString(root)
	if err != nil {
		return false, err
	}

	var fsName [windows.MAX_PATH + 1]uint16
	err = windows.GetVolumeInformation(rootUTF16, nil, 0, nil, nil, nil, &fsName[0], windows.MAX_PATH+1)
	if err != nil {
		return false, err
	}

	name := strings.ToLower(strings.TrimSpace(windows.UTF16ToString(fsName[:])))
	switch name {
	case "fat", "fat32", "exfat":
		return true, nil
	default:
		if strings.HasPrefix(name, "fat") {
			return true, nil
		}
		return false, nil
	}
}
