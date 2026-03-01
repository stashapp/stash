//go:build darwin
// +build darwin

package desktop

import (
	"fmt"
	"os"
	"os/exec"

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
	script := fmt.Sprintf(`display notification %q with title %q`, notificationText, notificationTitle)
	if err := exec.Command("osascript", "-e", script).Run(); err != nil {
		logger.Errorf("Could not send MacOS notification: %s", err.Error())
	}
}

func revealInFileManager(path string, _ os.FileInfo) error {
	if err := exec.Command(`open`, `-R`, path).Run(); err != nil {
		return fmt.Errorf("error revealing path in Finder: %w", err)
	}
	return nil
}

func isDoubleClickLaunched() bool {
	return false
}

func hideConsole() {

}
