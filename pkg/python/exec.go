// Package python provides utilities for working with the python executable.
package python

import (
	"context"
	"fmt"
	"os/exec"

	stashExec "github.com/stashapp/stash/pkg/exec"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
)

type Python string

func (p *Python) Command(ctx context.Context, args []string) *exec.Cmd {
	return stashExec.CommandContext(ctx, string(*p), args...)
}

// New returns a new Python instance at the given path.
func New(path string) *Python {
	ret := Python(path)
	return &ret
}

// Resolve tries to find the python executable in the system.
// It first checks for python3, then python.
// Returns nil and an exec.ErrNotFound error if not found.
func Resolve(configuredPythonPath string) (*Python, error) {
	if configuredPythonPath != "" {
		isFile, err := fsutil.FileExists(configuredPythonPath)
		switch {
		case err == nil && isFile:
			logger.Tracef("using configured python path: %s", configuredPythonPath)
			return New(configuredPythonPath), nil
		case err == nil && !isFile:
			logger.Warnf("configured python path is not a file: %s", configuredPythonPath)
		case err != nil:
			logger.Warnf("unable to use configured python path: %v", err)
		}
	}

	python3, err := exec.LookPath("python3")

	if err != nil {
		python, err := exec.LookPath("python")
		if err != nil {
			return nil, fmt.Errorf("python executable not in PATH: %w", err)
		}
		ret := Python(python)
		return &ret, nil
	}

	ret := Python(python3)
	return &ret, nil
}

// IsPythonCommand returns true if arg is "python" or "python3"
func IsPythonCommand(arg string) bool {
	return arg == "python" || arg == "python3"
}
