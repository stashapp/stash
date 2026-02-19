//go:build darwin
// +build darwin

package desktop

import (
	"os"
	"os/exec"

	"github.com/kermieisinthehouse/gosx-notifier"
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
	notification.AppIcon = getIconPath()
	notification.Open = getServerURL("")
	notification.Sender = "cc.stashapp.stash"
	err := notification.Push()

	if err != nil {
		logger.Errorf("Could not send MacOS notification: %s", err.Error())
	}
}

func revealInFileManager(path string, _ os.FileInfo) {
	if err := exec.Command(`open`, `-R`, path).Run(); err != nil {
		logger.Errorf("Error revealing path in Finder: %s", err.Error())
	}
}

func isDoubleClickLaunched() bool {
	return false
}

func hideConsole() {

}
