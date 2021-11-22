//go:build darwin || freebsd || linux
// +build darwin freebsd linux

package desktop

import (
	"io/ioutil"
	"os"
	"strings"
)

// isService checks if started by init, e.g. stash is a *nix systemd service / MacOS launchd service
func isService() bool {
	if os.Getppid() == 1 {
		return true
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
