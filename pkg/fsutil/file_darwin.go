//go:build darwin
// +build darwin

package fsutil

import "golang.org/x/text/unicode/norm"

// NormalizePath returns path normalized to Unicode NFC (composed) form.
//
// macOS filesystems report filenames in NFD (decomposed) form but are
// normalization-insensitive, so an NFD file on disk is still found via its
// NFC path.
func NormalizePath(path string) string {
	return norm.NFC.String(path)
}
