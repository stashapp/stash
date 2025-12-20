//go:build !linux
// +build !linux

package manager

import (
	"os"

	"github.com/syncthing/notify"
)

func notifyEvents() []notify.Event {
	return []notify.Event{
		notify.Create,
		notify.Rename,
		notify.Write,
	}
}

func notifyShouldScanEvent(fs os.FileInfo, ev notify.EventInfo) bool {
	return true
}
