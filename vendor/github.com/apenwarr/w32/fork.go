// Copyright 2010-2012 The W32 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package w32

// #include <stdlib.h>
//import (
//	"C"
//)

// Based on C code found here https://gist.github.com/juntalis/4366916
// Original code license:
/*
 * fork.c
 * Experimental fork() on Windows.  Requires NT 6 subsystem or
 * newer.
 *
 * Copyright (c) 2012 William Pitcock <nenolod@dereferenced.org>
 *
 * Permission to use, copy, modify, and/or distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * This software is provided 'as is' and without any warranty, express or
 * implied.  In no event shall the authors be liable for any damages arising
 * from the use of this software.
 */

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	ntdll = syscall.NewLazyDLL("ntdll.dll")

	procRtlCloneUserProcess = ntdll.NewProc("RtlCloneUserProcess")
	procAllocConsole        = modkernel32.NewProc("AllocConsole")
	procOpenProcess         = modkernel32.NewProc("OpenProcess")
	procOpenThread          = modkernel32.NewProc("OpenThread")
	procResumeThread        = modkernel32.NewProc("ResumeThread")
)

func OpenProcess(desiredAccess int, inheritHandle bool, processId uintptr) (h HANDLE, e error) {
	inherit := uintptr(0)
	if inheritHandle {
		inherit = 1
	}

	ret, _, lastErr := procOpenProcess.Call(
		uintptr(desiredAccess),
		inherit,
		uintptr(processId),
	)

	if ret == 0 {
		e = lastErr
	}

	h = HANDLE(ret)
	return
}

func OpenThread(desiredAccess int, inheritHandle bool, threadId uintptr) (h HANDLE, e error) {
	inherit := uintptr(0)
	if inheritHandle {
		inherit = 1
	}

	ret, _, lastErr := procOpenThread.Call(
		uintptr(desiredAccess),
		inherit,
		uintptr(threadId),
	)

	if ret == 0 {
		e = lastErr
	}

	h = HANDLE(ret)
	return
}

// DWORD WINAPI ResumeThread(
//   _In_ HANDLE hThread
// );
func ResumeThread(ht HANDLE) (e error) {

	ret, _, lastErr := procResumeThread.Call(
		uintptr(ht),
	)
	if ret == ^uintptr(0) { // -1
		e = lastErr
	}
	return
}

// BOOL WINAPI AllocConsole(void);
func AllocConsole() (e error) {
	ret, _, lastErr := procAllocConsole.Call()
	if ret != ERROR_SUCCESS {
		e = lastErr
	}
	return
}

// NTSYSAPI
// NTSTATUS
// NTAPI RtlCloneUserProcess  (
//   _In_ ULONG  ProcessFlags,
//   _In_opt_ PSECURITY_DESCRIPTOR  ProcessSecurityDescriptor,
//   _In_opt_ PSECURITY_DESCRIPTOR  ThreadSecurityDescriptor,
//   _In_opt_ HANDLE  DebugPort,
//   _Out_ PRTL_USER_PROCESS_INFORMATION  ProcessInformation
//  )

func RtlCloneUserProcess(
	ProcessFlags uint32,
	ProcessSecurityDescriptor, ThreadSecurityDescriptor *SECURITY_DESCRIPTOR, // in advapi32_typedef.go
	DebugPort HANDLE,
	ProcessInformation *RTL_USER_PROCESS_INFORMATION,
) (status uintptr) {

	status, _, _ = procRtlCloneUserProcess.Call(
		uintptr(ProcessFlags),
		uintptr(unsafe.Pointer(ProcessSecurityDescriptor)),
		uintptr(unsafe.Pointer(ThreadSecurityDescriptor)),
		uintptr(DebugPort),
		uintptr(unsafe.Pointer(ProcessInformation)),
	)

	return
}

// Fork creates a clone of the current process using the undocumented
// RtlCloneUserProcess call in ntdll, similar to unix fork(). The
// return value in the parent is the child PID. In the child it is 0.
func Fork() (pid uintptr, e error) {

	pi := &RTL_USER_PROCESS_INFORMATION{}

	ret := RtlCloneUserProcess(
		RTL_CLONE_PROCESS_FLAGS_CREATE_SUSPENDED|RTL_CLONE_PROCESS_FLAGS_INHERIT_HANDLES,
		nil,
		nil,
		HANDLE(0),
		pi,
	)

	switch ret {
	case RTL_CLONE_PARENT:
		pid = pi.ClientId.UniqueProcess
		ht, err := OpenThread(THREAD_ALL_ACCESS, false, pi.ClientId.UniqueThread)
		if err != nil {
			e = fmt.Errorf("OpenThread: %s", err)
		}
		err = ResumeThread(ht)
		if err != nil {
			e = fmt.Errorf("ResumeThread: %s", err)
		}
		CloseHandle(ht)
	case RTL_CLONE_CHILD:
		pid = 0
		err := AllocConsole()
		if err != nil {
			e = fmt.Errorf("AllocConsole: %s", err)
		}
	default:
		e = fmt.Errorf("0x%x", ret)
	}
	return
}
