//go:build !darwin
// +build !darwin

package fsutil

// NormalizePath returns path unchanged.
//
// Only macOS is normalization-insensitive; normalizing on a
// normalization-sensitive filesystem (e.g. Linux) would prevent files with
// NFD names on disk from being found.
func NormalizePath(path string) string {
	return path
}
