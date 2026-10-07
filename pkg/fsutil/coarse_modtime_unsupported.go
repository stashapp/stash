//go:build !linux && !darwin && !windows

package fsutil

// CoarseModtimeFilesystem reports whether path is on FAT, exFAT, or another
// filesystem where directory modification time is too coarse or unreliable for
// Stash's folder modtime scan optimization.
//
// On platforms without detection logic, this always returns false so that the
// folder ModTime skip optimisation is not inhibited by missing classification.
func CoarseModtimeFilesystem(_ string) (bool, error) {
	return false, nil
}
