//go:build linux
// +build linux

package desktop

import (
	"os/exec"

	"github.com/stashapp/stash/pkg/logger"
)

// isService checks if started by init, e.g. stash is a *nix systemd service
func isService() bool {
	// at the moment there is no desktop integration for FreeBSD
	return true
}

// there is no docker on FreeBSD
func isServerDockerized() bool {
	return false
}

func sendNotification(notificationTitle string, notificationText string) {
	err := exec.Command("notify-send", "-i", getIconPath(), notificationTitle, notificationText, "-a", "Stash").Run()
	if err != nil {
		logger.Errorf("Error sending notification on Linux: %s", err.Error())
	}
}

func revealInFileManager(path string) {

}
