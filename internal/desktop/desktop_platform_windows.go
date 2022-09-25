//go:build windows
// +build windows

package desktop

import (
	"os/exec"
	"syscall"
	"unsafe"

	"github.com/go-toast/toast"
	"github.com/stashapp/stash/pkg/logger"
	"golang.org/x/sys/windows/svc"
)

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	user32   = syscall.NewLazyDLL("user32.dll")
)

func isService() bool {
	result, err := svc.IsWindowsService()
	if err != nil {
		logger.Errorf("Encountered error checking if running as Windows service: %s", err.Error())
		return false
	}
	return result
}

// Detect if windows golang executable file is running via double click or from cmd/shell terminator
// https://stackoverflow.com/questions/8610489/distinguish-if-program-runs-by-clicking-on-the-icon-typing-its-name-in-the-cons?rq=1
// https://github.com/shirou/w32/blob/master/kernel32.go
// https://github.com/kbinani/win/blob/master/kernel32.go#L3268
// win.GetConsoleProcessList(new(uint32), win.DWORD(2))
// from https://gist.github.com/yougg/213250cc04a52e2b853590b06f49d865
func isDoubleClickLaunched() bool {
	lp := kernel32.NewProc("GetConsoleProcessList")
	if lp != nil {
		var pids [2]uint32
		var maxCount uint32 = 2
		ret, _, _ := lp.Call(uintptr(unsafe.Pointer(&pids)), uintptr(maxCount))
		if ret > 1 {
			return false
		}
	}
	return true
}

func hideConsole() {
	const SW_HIDE = 0
	h := getConsoleWindow()
	lp := user32.NewProc("ShowWindow")

	// don't want to check for errors and can't prevent dogsled
	_, _, _ = lp.Call(h, SW_HIDE) //nolint:dogsled
}

func getConsoleWindow() uintptr {
	lp := kernel32.NewProc("GetConsoleWindow")
	ret, _, _ := lp.Call()
	return ret
}

func isServerDockerized() bool {
	return false
}

func sendNotification(notificationTitle string, notificationText string) {
	notification := toast.Notification{
		AppID:   "Stash",
		Title:   notificationTitle,
		Message: notificationText,
		Icon:    getIconPath(),
		Actions: []toast.Action{{
			Type:      "protocol",
			Label:     "Open Stash",
			Arguments: getServerURL(""),
		}},
	}
	err := notification.Push()
	if err != nil {
		logger.Errorf("Error creating Windows notification: %s", err.Error())
	}
}

func revealInFileManager(path string) {
	exec.Command(`explorer`, `\select`, path)
}
