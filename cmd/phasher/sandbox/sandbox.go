//go:build !openbsd && !linux
// +build !openbsd,!linux

package sandbox

import "fmt"

func SandboxPHasher(ffmpegPath string, ffprobePath string, args []string) {
	fmt.Printf("Sandboxing is not yet implemented for your platform")
}
