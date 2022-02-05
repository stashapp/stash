package image

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/stashapp/stash/pkg/desktop"
	"github.com/stashapp/stash/pkg/logger"
)

type vipsEncoder string

func (e *vipsEncoder) ImageThumbnail(image *bytes.Buffer, maxSize int) ([]byte, error) {
	args := []string{
		"thumbnail_source",
		"[descriptor=0]",
		".jpg[Q=70,strip]",
		fmt.Sprint(maxSize),
		"--size", "down",
	}
	data, err := e.run(args, image)

	return []byte(data), err
}

func (e *vipsEncoder) run(args []string, stdin *bytes.Buffer) (string, error) {
	cmd := exec.Command(string(*e), args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = stdin

	desktop.HideExecShell(cmd)
	if err := cmd.Start(); err != nil {
		return "", err
	}

	err := cmd.Wait()

	if err != nil {
		// error message should be in the stderr stream
		logger.Errorf("image encoder error when running command <%s>: %s", strings.Join(cmd.Args, " "), stderr.String())
		return stdout.String(), err
	}

	return stdout.String(), nil
}
