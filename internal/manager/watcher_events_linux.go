//go:build linux
// +build linux

package manager

import "github.com/syncthing/notify"

func notifyEvents() []notify.Event {
	return []notify.Event{
		notify.InCloseWrite,
		notify.InMovedTo,
		notify.Rename,
	}
}
