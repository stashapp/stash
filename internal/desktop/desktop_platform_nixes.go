//go:build unix && !darwin
// +build unix,!darwin

package desktop

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/stashapp/stash/pkg/logger"
)

// isService checks if started by init, e.g. stash is a *nix systemd service
func isService() bool {
	return os.Getppid() == 1
}

func isServerDockerized() bool {
	_, dockerEnvErr := os.Stat("/.dockerenv")
	cgroups, _ := os.ReadFile("/proc/self/cgroup")
	if !os.IsNotExist(dockerEnvErr) || strings.Contains(string(cgroups), "docker") {
		return true
	}

	return false
}

func sendNotification(notificationTitle string, notificationText string) {
	err := exec.Command("notify-send", "-i", getIconPath(), notificationTitle, notificationText, "-a", "Stash").Run()
	if err != nil {
		logger.Errorf("Error sending notification on Linux: %s", err.Error())
	}
}

func revealInFileManager(path string, info os.FileInfo) error {
	dir := path
	if !info.IsDir() {
		dir = filepath.Dir(path)
	}
	if err := exec.Command("xdg-open", dir).Run(); err != nil {
		return fmt.Errorf("error opening directory in file manager: %w", err)
	}
	return nil
}

func isDoubleClickLaunched() bool {
	return false
}

func hideConsole() {

}
