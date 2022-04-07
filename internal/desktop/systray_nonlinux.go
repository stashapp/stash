//go:build windows || darwin
// +build windows darwin

package desktop

import (
	"strings"

	"github.com/kermieisinthehouse/systray"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/logger"
)

// MUST be run on the main goroutine or will have no effect on macOS
func startSystray(shutdownHandler ShutdownHandler, faviconProvider FaviconProvider) {

	// Shows a small notification to inform that Stash will no longer show a terminal window,
	// and instead will be available in the tray. Will only show the first time a pre-desktop integration
	// system is started from a non-terminal method, e.g. double-clicking an icon.
	c := config.GetInstance()
	if c.GetShowOneTimeMovedNotification() {
		SendNotification("Stash has moved!", "Stash now runs in your tray, instead of a terminal window.")
		c.Set(config.ShowOneTimeMovedNotification, false)
		if err := c.Write(); err != nil {
			logger.Errorf("Error while writing configuration file: %s", err.Error())
		}
	}

	// Listen for changes to rerender systray
	// TODO: This is disabled for now. The systray package does not clean up all of its resources when Quit() is called.
	// TODO: This results in this only working once, or changes being ignored. Our fork of systray fixes a crash(!) on macOS here.
	// go func() {
	// 	for {
	// 		<-config.GetInstance().GetConfigUpdatesChannel()
	// 		systray.Quit()
	// 	}
	// }()

	for {
		systray.Run(func() {
			systrayInitialize(shutdownHandler, faviconProvider)
		}, nil)
	}
}

func systrayInitialize(shutdownHandler ShutdownHandler, faviconProvider FaviconProvider) {
	favicon := faviconProvider.GetFavicon()
	systray.SetTemplateIcon(favicon, favicon)
	systray.SetTooltip("🟢 Stash is Running.")

	openStashButton := systray.AddMenuItem("Open Stash", "Open a browser window to Stash")
	var menuItems []string
	systray.AddSeparator()
	c := config.GetInstance()
	if !c.IsNewSystem() {
		menuItems = c.GetMenuItems()
		for _, item := range menuItems {
			titleCaseItem := strings.Title(strings.ToLower(item))
			curr := systray.AddMenuItem(titleCaseItem, "Open to "+titleCaseItem)
			go func(item string) {
				for {
					<-curr.ClickedCh
					if item == "markers" {
						item = "scenes/markers"
					}
					openURLInBrowser(item)
				}
			}(item)
		}
		systray.AddSeparator()
		// TODO - Some ideas for future expansions
		// systray.AddMenuItem("Start a Scan", "Scan all libraries with default settings")
		// systray.AddMenuItem("Start Auto Tagging", "Auto Tag all libraries")
		// systray.AddMenuItem("Check for updates", "Check for a new Stash release")
		// systray.AddSeparator()
	}

	quitStashButton := systray.AddMenuItem("Quit Stash Server", "Quits the Stash server")

	go func() {
		for {
			select {
			case <-openStashButton.ClickedCh:
				openURLInBrowser("")
			case <-quitStashButton.ClickedCh:
				systray.Quit()
				shutdownHandler.Shutdown(0)
			}
		}
	}()
}
