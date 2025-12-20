//go:build linux
// +build linux

package manager

import (
	"os"

	"github.com/syncthing/notify"
)

func notifyEvents() []notify.Event {
	return []notify.Event{
		notify.InCloseWrite,
		notify.InCreate,
		notify.InMovedTo,
		notify.Rename,
	}
}

func notifyShouldScanEvent(fs os.FileInfo, ev notify.EventInfo) bool {
	event := ev.Event()

	// directories only, files fire too early
	if event&notify.InCreate != 0 {
		return fs.IsDir()
	}

	return true
}
