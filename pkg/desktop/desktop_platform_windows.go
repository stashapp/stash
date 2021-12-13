//go:build windows
// +build windows

package desktop

import (
	"golang.org/x/sys/windows"
	"os/exec"
	"syscall"

	"github.com/go-toast/toast"
	"github.com/stashapp/stash/pkg/logger"
	"golang.org/x/sys/windows/svc"
)

func isService() bool {
	result, err := svc.IsWindowsService()
	if err != nil {
		logger.Errorf("Encountered error checking if running as Windows service: %s", err.Error())
		return false
	}
	return result
}

func isServerDockerized() bool {
	return false
}

// On Windows, calling exec.Cmd.Start() will create a cmd window, even if we live in the taskbar.
// We don't want every ffmpeg / plugin to pop up a window.
func hideExecShell(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: windows.CREATE_NO_WINDOW}
}

func sendNotification(notificationTitle string, notificationText string) {
	notification := toast.Notification{
		AppID:   "Stash",
		Title:   notificationTitle,
		Message: notificationText,
		Icon:    getIconPath(),
		Actions: []toast.Action{
			{"protocol", "Open Stash", getServerURL("")}}}
	err := notification.Push()
	if err != nil {
		logger.Errorf("Error creating Windows notification: %s", err.Error())
	}
}
