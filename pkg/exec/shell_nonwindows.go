//go:build linux || darwin || !windows
// +build linux darwin !windows

package exec

import "os/exec"

// hideExecShell does nothing on non-Windows platforms.
func hideExecShell(cmd *exec.Cmd) {
}
