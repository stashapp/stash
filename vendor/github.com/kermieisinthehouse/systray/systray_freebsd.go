//go:build freebsd
// +build freebsd

package systray

func registerSystray() {
}

func nativeLoop() {
}

func quit() {
}

// SetIcon sets the systray icon.
// iconBytes should be the content of .ico for windows and .ico/.jpg/.png
// for other platforms.
func SetIcon(iconBytes []byte) {
}

// SetTitle sets the systray title, only available on Mac and Linux.
func SetTitle(title string) {
}

// SetTooltip sets the systray tooltip to display on mouse hover of the tray icon,
// only available on Mac and Windows.
func SetTooltip(tooltip string) {
}

func addOrUpdateMenuItem(item *MenuItem) {
}

func addSeparator(id uint32) {
}

func hideMenuItem(item *MenuItem) {
}

func showMenuItem(item *MenuItem) {
}

//export systray_ready
func systray_ready() {
}

//export systray_on_exit
func systray_on_exit() {
}

//export systray_menu_item_selected
func systray_menu_item_selected(cID int) {
}
