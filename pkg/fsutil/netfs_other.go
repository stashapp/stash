//go:build !linux && !darwin && !windows

package fsutil

// IsNetworkFS reports whether path resides on a network filesystem.
// On unsupported platforms we cannot classify mounts, so we always return
// true (treat as network) and disable the folder ModTime skip optimisation.
func IsNetworkFS(_ string) (bool, error) {
	return true, nil
}
