//go:build linux

package fsutil

import (
	"golang.org/x/sys/unix"
)

// Network FS magic numbers from statfs(2)
const (
	// nfsSuper and related remote FS magic numbers
	nfsMagic     = 0x6969
	smbMagic     = 0xFF534D42 // SMBv1 / CIFS
	smb2Magic    = 0xFE534D42 // SMBv2
	cifsMagic    = 0x517B
	afsNetMagic  = 0x5346414F // AFS
	ncpMagic     = 0x564C     // Novell NCP
	cephMagic    = 0x00C36400
	glusterMagic = 0x65735546 // also FUSE magic — treat conservatively as unknown (network)
	// FUSE magic covers e.g. sshfs, gvfs — conservatively treated as network
	fuseMagic = 0x65735546
)

// IsNetworkFS reports whether path resides on a network (or FUSE) filesystem.
// On Linux, it uses statfs(2) magic numbers. FUSE filesystems are treated as
// network filesystems because they may not provide reliable ModTime semantics.
func IsNetworkFS(path string) (bool, error) {
	var s unix.Statfs_t
	if err := unix.Statfs(path, &s); err != nil {
		return false, err
	}

	// Type is int32 on e.g. linux/arm; SMB/CIFS f_type magics exceed MaxInt32 — compare as uint32.
	switch uint32(s.Type) {
	case nfsMagic, smbMagic, smb2Magic, cifsMagic, afsNetMagic, ncpMagic, cephMagic, fuseMagic:
		return true, nil
	}

	return false, nil
}
