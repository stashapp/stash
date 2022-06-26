//go:build windows
// +build windows

package exec

import (
	"os/exec"
	"syscall"

	"golang.org/x/sys/windows"
)

// hideExecShell hides the windows when executing on Windows.
func hideExecShell(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: windows.DETACHED_PROCESS & windows.CREATE_NO_WINDOW}
}
