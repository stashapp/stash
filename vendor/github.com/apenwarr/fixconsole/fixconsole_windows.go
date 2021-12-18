package fixconsole

import (
	"fmt"
	"github.com/apenwarr/w32"
	"golang.org/x/sys/windows"
	"os"
	"syscall"
)

func AttachConsole() error {
	const ATTACH_PARENT_PROCESS = ^uintptr(0)
	proc := syscall.MustLoadDLL("kernel32.dll").MustFindProc("AttachConsole")
	r1, _, err := proc.Call(ATTACH_PARENT_PROCESS)
	if r1 == 0 {
		errno, ok := err.(syscall.Errno)
		if ok && errno == w32.ERROR_INVALID_HANDLE {
			// console handle doesn't exist; not a real
			// error, but the console handle will be
			// invalid.
			return nil
		}
		return err
	} else {
		return nil
	}
}

var oldStdin, oldStdout, oldStderr *os.File

// Windows console output is a mess.
//
// If you compile as "-H windows", then if you launch your program without
// a console, Windows forcibly creates one to use as your stdin/stdout, which
// is silly for a GUI app, so we can't do that.
//
// If you compile as "-H windowsgui", then it doesn't create a console for
// your app... but also doesn't provide a working stdin/stdout/stderr even if
// you *did* launch from the console.  However, you can use AttachConsole()
// to get a handle to your parent process's console, if any, and then
// os.NewFile() to turn that handle into a fd usable as stdout/stderr.
//
// However, then you have the problem that if you redirect stdout or stderr
// from the shell, you end up ignoring the redirection by forcing it to the
// console.
//
// To fix *that*, we have to detect whether there was a pre-existing stdout
// or not. We can check GetStdHandle(), which returns 0 for "should be
// console" and nonzero for "already pointing at a file."
//
// Be careful though!  As soon as you run AttachConsole(), it resets *all*
// the GetStdHandle() handles to point them at the console instead, thus
// throwing away the original file redirects.  So we have to GetStdHandle()
// *before* AttachConsole().
//
// For some reason, powershell redirections provide a valid file handle, but
// writing to that handle doesn't write to the file.  I haven't found a way
// to work around that.  (Windows 10.0.17763.379)
//
// Net result is as follows.
// Before:
//    SHELL            NON-REDIRECTED     REDIRECTED
//    explorer.exe     no console         n/a
//    cmd.exe          broken             works
//    powershell       broken             broken
//    WSL bash         broken             works
// After
//    SHELL            NON-REDIRECTED     REDIRECTED
//    explorer.exe     no console         n/a
//    cmd.exe          works              works
//    powershell       works              broken
//    WSL bash         works              works
//
// We don't seem to make anything worse, at least.
func FixConsoleIfNeeded() error {
	// Retain the original console objects, to prevent Go from automatically
	// closing their file descriptors when they get garbage collected.
	// You never want to close file descriptors 0, 1, and 2.
	oldStdin, oldStdout, oldStderr = os.Stdin, os.Stdout, os.Stderr

	stdin, _ := syscall.GetStdHandle(syscall.STD_INPUT_HANDLE)
	stdout, _ := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
	stderr, _ := syscall.GetStdHandle(syscall.STD_ERROR_HANDLE)

	var invalid syscall.Handle
	con := invalid

	if stdin == invalid || stdout == invalid || stderr == invalid {
		err := AttachConsole()
		if err != nil {
			return fmt.Errorf("attachconsole: %v", err)
		}

		if stdin == invalid {
			stdin, _ = syscall.GetStdHandle(syscall.STD_INPUT_HANDLE)
		}
		if stdout == invalid {
			stdout, _ = syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
			con = stdout
		}
		if stderr == invalid {
			stderr, _ = syscall.GetStdHandle(syscall.STD_ERROR_HANDLE)
			con = stderr
		}
	}

	if con != invalid {
		// Make sure the console is configured to convert
		// \n to \r\n, like Go programs expect.
		h := windows.Handle(con)
		var st uint32
		err := windows.GetConsoleMode(h, &st)
		if err != nil {
			return fmt.Errorf("GetConsoleMode: %v", err)
		}
		err = windows.SetConsoleMode(h, st&^windows.DISABLE_NEWLINE_AUTO_RETURN)
		if err != nil {
			return fmt.Errorf("SetConsoleMode: %v", err)
		}
	}

	if stdin != invalid {
		os.Stdin = os.NewFile(uintptr(stdin), "stdin")
	}
	if stdout != invalid {
		os.Stdout = os.NewFile(uintptr(stdout), "stdout")
	}
	if stderr != invalid {
		os.Stderr = os.NewFile(uintptr(stderr), "stderr")
	}
	return nil
}
