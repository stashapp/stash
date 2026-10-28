//go:build linux

package fsutil

import (
	"golang.org/x/sys/unix"
)

// Magic numbers from Linux uapi/linux/magic.h (statfs.f_type).
const (
	msdosSuperMagic = 0x4d44     // FAT / VFAT (msdos)
	fatSuperMagic   = 0x4006     // legacy FAT magic (rare)
	exfatSuperMagic = 0x2011BAB0 // exFAT
)

// CoarseModtimeFilesystem reports whether path is on FAT, exFAT, or another
// filesystem where directory modification time is too coarse or unreliable for
// Stash's folder modtime scan optimization.
func CoarseModtimeFilesystem(path string) (bool, error) {
	var s unix.Statfs_t
	if err := unix.Statfs(path, &s); err != nil {
		return false, err
	}

	switch s.Type {
	case msdosSuperMagic, fatSuperMagic, exfatSuperMagic:
		return true, nil
	default:
		return false, nil
	}
}
