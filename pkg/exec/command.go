// Package exec provides functions that wrap os/exec functions. These functions prevent external commands from opening windows on the Windows platform.
package exec

import (
	"context"
	"os/exec"
)

// Command wraps the exec.Command function, preventing Windows from opening a window when starting.
func Command(name string, arg ...string) *exec.Cmd {
	ret := exec.Command(name, arg...)
	hideExecShell(ret)
	return ret
}

// CommandContext wraps the exec.CommandContext function, preventing Windows from opening a window when starting.
func CommandContext(ctx context.Context, name string, arg ...string) *exec.Cmd {
	ret := exec.CommandContext(ctx, name, arg...)
	hideExecShell(ret)
	return ret
}
