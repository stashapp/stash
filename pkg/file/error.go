package file

import (
	"errors"

	"golang.org/x/sys/windows"
)

// IsBadNetworkError checks if the error is a "bad network name" error, which can occur when a network share is disconnected.
func IsBadNetworkError(err error) bool {
	return errors.Is(err, windows.ERROR_BAD_NETPATH)
}
