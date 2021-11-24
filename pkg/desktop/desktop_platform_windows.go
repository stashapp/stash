//go:build windows
// +build windows

package desktop

import (
	"github.com/stashapp/stash/pkg/logger"
	"golang.org/x/sys/windows/svc"
)

func isService() bool {
	result, err := svc.IsWindowsService()
	if err != nil {
		logger.Errorf("Encountered error checking if running as Windows service: %s", err)
		return false
	}
	return result
}

func isServerDockerized() bool {
	return false
}
