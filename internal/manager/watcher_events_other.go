//go:build !linux
// +build !linux

package manager

import "github.com/syncthing/notify"

func notifyEvents() []notify.Event {
	return []notify.Event{
		notify.Create,
		notify.Rename,
		notify.Write,
	}
}

func notifyShouldScanEvent(ev notify.EventInfo) bool {
	return true
}
