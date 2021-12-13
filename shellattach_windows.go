//go:build windows
// +build windows

package main

import (
	"syscall"
)

const (
	ATTACH_PARENT_PROCESS = ^uint32(0) // (DWORD)-1
)

var (
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")

	procAttachConsole = modkernel32.NewProc("AttachConsole")
)

func AttachConsole(dwParentProcess uint32) (ok bool) {
	r0, _, _ := syscall.Syscall(procAttachConsole.Addr(), 1, uintptr(dwParentProcess), 0, 0)
	ok = bool(r0 != 0)
	return
}
