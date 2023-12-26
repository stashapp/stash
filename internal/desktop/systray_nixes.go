//go:build (!windows && !darwin) || !cgo

package desktop

func startSystray(exit chan int, favicon FaviconProvider) {
	// The systray is not available on Linux because the required libraries (libappindicator3 and gtk+3.0)
	// are not able to be statically compiled. Technically, the systray works perfectly fine when dynamically
	// linked, but we cannot distribute it for compatibility reasons.
	// Additionally, the systray package requires CGo so the dependency cannot be used if building with
	// CGo disabled.
}
