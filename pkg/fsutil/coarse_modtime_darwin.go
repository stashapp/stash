//go:build darwin

package fsutil

import (
	"strings"

	"golang.org/x/sys/unix"
)

// CoarseModtimeFilesystem reports whether path is on FAT, exFAT, or another
// filesystem where directory modification time is too coarse or unreliable for
// Stash's folder modtime scan optimization.
func CoarseModtimeFilesystem(path string) (bool, error) {
	var s unix.Statfs_t
	if err := unix.Statfs(path, &s); err != nil {
		return false, err
	}

	b := make([]byte, 0, len(s.Fstypename))
	for _, c := range s.Fstypename {
		if c == 0 {
			break
		}
		b = append(b, byte(c))
	}
	name := strings.ToLower(string(b))

	switch name {
	case "msdos", "exfat":
		return true, nil
	default:
		return false, nil
	}
}
