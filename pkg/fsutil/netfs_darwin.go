//go:build darwin

package fsutil

import (
	"strings"

	"golang.org/x/sys/unix"
)

// networkFSTypes are fstypenames that indicate a network or FUSE filesystem.
var networkFSTypes = []string{
	"nfs",
	"smbfs",
	"afpfs",
	"webdav",
	"osxfuse",
	"macfuse",
	"fuse",
}

// IsNetworkFS reports whether path resides on a network (or FUSE) filesystem.
// On macOS, it uses statfs(2) f_fstypename. FUSE filesystems are treated as
// network filesystems because they may not provide reliable ModTime semantics.
func IsNetworkFS(path string) (bool, error) {
	var s unix.Statfs_t
	if err := unix.Statfs(path, &s); err != nil {
		return false, err
	}

	// f_fstypename is a [16]int8; convert to string
	b := make([]byte, 0, len(s.Fstypename))
	for _, c := range s.Fstypename {
		if c == 0 {
			break
		}
		b = append(b, byte(c))
	}
	name := strings.ToLower(string(b))

	for _, t := range networkFSTypes {
		if name == t || strings.HasPrefix(name, t) {
			return true, nil
		}
	}

	return false, nil
}
