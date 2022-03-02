//go:build windows
// +build windows

package desktop

import (
	"os/exec"

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

func sendNotification(notificationTitle string, notificationText string) {
	notification := toast.Notification{
		AppID:   "Stash",
		Title:   notificationTitle,
		Message: notificationText,
		Icon:    getIconPath(),
		Actions: []toast.Action{{
			Type:      "protocol",
			Label:     "Open Stash",
			Arguments: getServerURL(""),
		}},
	}
	err := notification.Push()
	if err != nil {
		logger.Errorf("Error creating Windows notification: %s", err.Error())
	}
}

func revealInFileManager(path string) {
	exec.Command(`explorer`, `\select`, path)
}
