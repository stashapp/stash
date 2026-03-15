//go:build linux || darwin || !windows
// +build linux darwin !windows

package file

// IsBadNetworkError checks if the error is a "bad network name" error, which can occur when a network share is disconnected.
func IsBadNetworkError(err error) bool {
	return false
}
