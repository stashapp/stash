// Copyright 2010-2012 The W32 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package w32

import (
	"syscall"
	"unsafe"
)

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")

	procCreateProcessW      = kernel32.NewProc("CreateProcessW")
	procTerminateProcess    = kernel32.NewProc("TerminateProcess")
	procGetExitCodeProcess  = kernel32.NewProc("GetExitCodeProcess")
	procWaitForSingleObject = kernel32.NewProc("WaitForSingleObject")
)

// WINBASEAPI WINBOOL WINAPI
// CreateProcessW (
// LPCWSTR lpApplicationName,
// LPWSTR lpCommandLine,
// LPSECURITY_ATTRIBUTES lpProcessAttributes,
// LPSECURITY_ATTRIBUTES lpThreadAttributes
// WINBOOL bInheritHandles
// DWORD dwCreationFlags
// LPVOID lpEnvironment
// LPCWSTR lpCurrentDirectory
// LPSTARTUPINFOW lpStartupInfo
// LPPROCESS_INFORMATION lpProcessInformation
//);
func CreateProcessW(
	lpApplicationName, lpCommandLine string,
	lpProcessAttributes, lpThreadAttributes *SECURITY_ATTRIBUTES,
	bInheritHandles BOOL,
	dwCreationFlags uint32,
	lpEnvironment unsafe.Pointer,
	lpCurrentDirectory string,
	lpStartupInfo *STARTUPINFOW,
	lpProcessInformation *PROCESS_INFORMATION,
) (e error) {

	var lpAN, lpCL, lpCD *uint16
	if len(lpApplicationName) > 0 {
		lpAN, e = syscall.UTF16PtrFromString(lpApplicationName)
		if e != nil {
			return
		}
	}
	if len(lpCommandLine) > 0 {
		lpCL, e = syscall.UTF16PtrFromString(lpCommandLine)
		if e != nil {
			return
		}
	}
	if len(lpCurrentDirectory) > 0 {
		lpCD, e = syscall.UTF16PtrFromString(lpCurrentDirectory)
		if e != nil {
			return
		}
	}

	ret, _, lastErr := procCreateProcessW.Call(
		uintptr(unsafe.Pointer(lpAN)),
		uintptr(unsafe.Pointer(lpCL)),
		uintptr(unsafe.Pointer(lpProcessAttributes)),
		uintptr(unsafe.Pointer(lpProcessInformation)),
		uintptr(bInheritHandles),
		uintptr(dwCreationFlags),
		uintptr(lpEnvironment),
		uintptr(unsafe.Pointer(lpCD)),
		uintptr(unsafe.Pointer(lpStartupInfo)),
		uintptr(unsafe.Pointer(lpProcessInformation)),
	)

	if ret == 0 {
		e = lastErr
	}

	return
}

func CreateProcessQuick(cmd string) (pi PROCESS_INFORMATION, e error) {
	si := &STARTUPINFOW{}
	e = CreateProcessW(
		"",
		cmd,
		nil,
		nil,
		0,
		0,
		unsafe.Pointer(nil),
		"",
		si,
		&pi,
	)
	return
}

func TerminateProcess(hProcess HANDLE, exitCode uint32) (e error) {
	ret, _, lastErr := procTerminateProcess.Call(
		uintptr(hProcess),
		uintptr(exitCode),
	)

	if ret == 0 {
		e = lastErr
	}

	return
}

func GetExitCodeProcess(hProcess HANDLE) (code uintptr, e error) {
	ret, _, lastErr := procGetExitCodeProcess.Call(
		uintptr(hProcess),
		uintptr(unsafe.Pointer(&code)),
	)

	if ret == 0 {
		e = lastErr
	}

	return
}

// DWORD WINAPI WaitForSingleObject(
//   _In_ HANDLE hHandle,
//   _In_ DWORD  dwMilliseconds
// );

func WaitForSingleObject(hHandle HANDLE, msecs uint32) (ok bool, e error) {

	ret, _, lastErr := procWaitForSingleObject.Call(
		uintptr(hHandle),
		uintptr(msecs),
	)

	if ret == WAIT_OBJECT_0 {
		ok = true
		return
	}

	// don't set e for timeouts, or it will be ERROR_SUCCESS which is
	// confusing
	if ret != WAIT_TIMEOUT {
		e = lastErr
	}
	return

}
