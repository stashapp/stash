package python

import (
	"fmt"
	"os"
	"os/exec"
)

func AppendPythonPath(cmd *exec.Cmd, path string) {
	// Respect the users PYTHONPATH if set
	if currentValue, set := os.LookupEnv("PYTHONPATH"); set {
		path = fmt.Sprintf("%s%c%s", currentValue, os.PathListSeparator, path)
	}
	cmd.Env = append(os.Environ(), fmt.Sprintf("PYTHONPATH=%s", path))
}
