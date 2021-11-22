//go:build windows
// +build windows

package desktop

import "golang.org/x/sys/windows/svc"

func isService() bool {
	return svc.IsWindowsService()
}

func isServerDockerized() bool {
	return false
}
