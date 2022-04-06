//go:build linux || freebsd
// +build linux freebsd

package desktop

func startSystray(shutdownHandler ShutdownHandler, favicon FaviconProvider) {
	// The systray is not available on linux because the required libraries (libappindicator3 and gtk+3.0)
	// are not able to be statically compiled. Technically, the systray works perfectly fine when dynamically
	// linked, but we cannot distribute it for compatibility reasons.
}
