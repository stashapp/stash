// +build !windows

package fixconsole

// On non-windows platforms, we don't need to do anything. The console
// starts off attached already, if it exists.

func AttachConsole() error {
	return nil
}

func FixConsoleIfNeeded() error {
	return nil
}
