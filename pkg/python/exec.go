package python

import (
	"context"
	"os/exec"

	stashExec "github.com/stashapp/stash/pkg/exec"
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
func Resolve() (*Python, error) {
	python3, err := exec.LookPath("python3")
	if err == nil {
		return New(python3), nil
	}

	python, err := exec.LookPath("python")
	if err == nil {
		return New(python), nil
	}

	return nil, err
}

// IsPythonCommand returns true if arg is "python" or "python3"
func IsPythonCommand(arg string) bool {
	return arg == "python" || arg == "python3"
}
