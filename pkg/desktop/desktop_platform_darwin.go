//go:build darwin
// +build darwin

package desktop

import (
	"github.com/deckarep/gosx-notifier"
	"github.com/stashapp/stash/pkg/logger"
)

func isService() bool {
	// MacOS /does/ support services, using launchd, but there is no straightforward way to check if it was used.
	return false
}

func isServerDockerized() bool {
	return false
}

func sendNotification(notificationTitle string, notificationText string) {
	notification := gosxnotifier.NewNotification(notificationText)
	notification.Title = notificationTitle
	err := notification.Push()

	if err != nil {
		logger.Errorf("Could not send MacOS notification: %s", err.Error())
	}
}
