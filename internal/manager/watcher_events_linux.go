//go:build linux
// +build linux

package manager

import (
	"github.com/syncthing/notify"
	"golang.org/x/sys/unix"
)

func notifyEvents() []notify.Event {
	return []notify.Event{
		notify.InCloseWrite,
		notify.InCreate,
		notify.InMovedTo,
		notify.Rename,
	}
}

func notifyShouldScanEvent(ev notify.EventInfo) bool {
	event := ev.Event()

	// directories only, files fire too early
	if event&notify.InCreate != 0 {
		return event&unix.IN_ISDIR != 0
	}

	return true
}
