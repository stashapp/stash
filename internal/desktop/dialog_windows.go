//go:build windows
// +build windows

package desktop

import (
	"fmt"
	"syscall"
	"unsafe"
)

func FatalError(err error) int {
	const (
		NULL         = 0
		MB_OK        = 0
		MB_ICONERROR = 0x10
	)

	return messageBox(NULL, fmt.Sprintf("Error: %v", err), "Stash - Fatal Error", MB_OK|MB_ICONERROR)
}

func messageBox(hwnd uintptr, caption, title string, flags uint) int {
	lpText, _ := syscall.UTF16PtrFromString(caption)
	lpCaption, _ := syscall.UTF16PtrFromString(title)

	ret, _, _ := syscall.NewLazyDLL("user32.dll").NewProc("MessageBoxW").Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(lpText)),
		uintptr(unsafe.Pointer(lpCaption)),
		uintptr(flags))

	return int(ret)
}
