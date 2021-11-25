//go:build darwin || freebsd || linux
// +build darwin freebsd linux

package desktop

import (
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/0xAX/notificator"
)

// isService checks if started by init, e.g. stash is a *nix systemd service
func isService() bool {
	if runtime.GOOS != "darwin" {
		return os.Getppid() == 1
	}
	return false
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
