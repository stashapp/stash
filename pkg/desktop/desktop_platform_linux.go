//go:build linux
// +build linux

package desktop

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/0xAX/notificator"
)

// isService checks if started by init, e.g. stash is a *nix systemd service
func isService() bool {
	return os.Getppid() == 1
}

func isServerDockerized() bool {
	_, dockerEnvErr := os.Stat("/.dockerenv")
	cgroups, _ := ioutil.ReadFile("/proc/self/cgroup")
	if os.IsExist(dockerEnvErr) || strings.Contains(string(cgroups), "docker") {
		return true
	}

	return false
}

func sendNotification(notificationTitle string, notificationText string) {
	notificator.New(notificator.Options{
		AppName: "Stash",
	}).Push(notificationTitle, notificationText, "", notificator.UR_NORMAL)
}
