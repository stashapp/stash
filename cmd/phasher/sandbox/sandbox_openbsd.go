//go:build openbsd
// +build openbsd

package sandbox

import (
	"fmt"
	"golang.org/x/sys/unix"
	"os"
)

func SandboxPHasher(ffmpegPath string, ffprobePath string, args []string) {
	err := unix.Unveil("/dev/null", "rw")
	err = unix.Unveil(ffmpegPath, "x")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ffmpeg unveil failed with %s\n", err.Error())
		os.Exit(2)
	}
	err = unix.Unveil(ffprobePath, "x")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ffprobe unveil failed with %s\n", err.Error())
		os.Exit(2)
	}
	for _, item := range args {
		err = unix.Unveil(item, "r")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s unveil failed with %s\n", item, err.Error())
			os.Exit(2)
		}
	}
	err = unix.UnveilBlock()

	err = unix.PledgePromises("stdio rpath exec proc error")
	if err != nil {
		fmt.Fprintf(os.Stderr, "pledge failed with %s\n", err.Error())
		os.Exit(2)
	}
}
