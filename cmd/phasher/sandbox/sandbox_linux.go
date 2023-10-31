//go:build linux
// +build linux

package sandbox

import "fmt"

func SandboxPHasher(ffmpegPath string, ffprobePath string, args []string) {
	// TODO: SECCOMP and Landlock
	fmt.Printf("Sandboxing is not yet implemented for your platform")
}
