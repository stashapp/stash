//go:build windows
// +build windows

package desktop

import (
	"github.com/stashapp/stash/pkg/logger"
	"golang.org/x/sys/windows/svc"

	"github.com/go-toast/toast"
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
		Message: notificationText}
	err := notification.Push()
	if err != nil {
		logger.Errorf("Error creating Windows notification: %s", err.Error())
	}
}
