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
	switch ev.Event() {
	case notify.InCreate:
		// directories only, files fire too early
		return ev.Event()&unix.IN_ISDIR != 0
	}

	return true
}
