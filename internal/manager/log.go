package manager

import (
	"errors"
	"os/exec"

	"github.com/stashapp/stash/pkg/logger"
)

func logErrorOutput(err error) {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		logger.Errorf("command stderr: %v", string(exitErr.Stderr))
	}
}
