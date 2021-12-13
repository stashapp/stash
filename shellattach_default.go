//go:build !windows
// +build !windows

package main

const (
	ATTACH_PARENT_PROCESS = ^uint32(0) // (DWORD)-1
)

func AttachConsole(dwParentProcess uint32) (ok bool) {
	return true
}
