//go:build linux
// +build linux

package sandbox

import (
	"fmt"
	"github.com/landlock-lsm/go-landlock/landlock"
	"os"
)

func SandboxPHasher(ffmpegPath string, ffprobePath string, args []string) {
	err := landlock.V3.BestEffort().RestrictPaths(
		landlock.ROFiles("/dev/null"),
		landlock.ROFiles("/usr/lib"),
		landlock.ROFiles(args...),
		landlock.ROFiles(ffmpegPath),
		landlock.ROFiles(ffprobePath),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ffmpeg landlock failed with %s\n", err.Error())
	}

}
