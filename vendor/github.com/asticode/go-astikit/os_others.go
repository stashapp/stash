// +build !js !wasm

package astikit

import (
	"os"
	"syscall"
)

func isTermSignal(s os.Signal) bool {
	return s == syscall.SIGABRT || s == syscall.SIGKILL || s == syscall.SIGINT || s == syscall.SIGQUIT || s == syscall.SIGTERM
}
