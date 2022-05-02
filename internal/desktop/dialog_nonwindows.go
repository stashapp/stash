//go:build !windows
// +build !windows

package desktop

func FatalError(err error) int {
	// nothing to do
	return 0
}
