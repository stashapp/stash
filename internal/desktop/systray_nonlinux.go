//go:build (windows || darwin) && cgo

package desktop

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/kermieisinthehouse/systray"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/logger"
)

// MUST be run on the main goroutine or will have no effect on macOS
func startSystray(exit chan int, faviconProvider FaviconProvider) {
	// Shows a small notification to inform that Stash will no longer show a terminal window,
	// and instead will be available in the tray. Will only show the first time a pre-desktop integration
	// system is started from a non-terminal method, e.g. double-clicking an icon.
	c := config.GetInstance()
	if c.GetShowOneTimeMovedNotification() {
		// Use platform-appropriate terminology
		location := "tray"
		if runtime.GOOS == "darwin" {
			location = "menu bar"
		}
		SendNotification("Stash has moved!", "Stash now runs in your "+location+", instead of a terminal window.")
		c.SetBool(config.ShowOneTimeMovedNotification, false)
		if err := c.Write(); err != nil {
			logger.Errorf("Error while writing configuration file: %v", err)
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

	// "intercept" an exit code to quit the systray, allowing the call to systray.Run() below to return.
	go func() {
		exitCode := <-exit
		systray.Quit()
		exit <- exitCode
	}()

	systray.Run(func() {
		systrayInitialize(exit, faviconProvider)
	}, nil)
}

func systrayInitialize(exit chan<- int, faviconProvider FaviconProvider) {
	favicon := faviconProvider.GetFavicon()
	systray.SetTemplateIcon(favicon, favicon)
	c := config.GetInstance()
	systray.SetTooltip(fmt.Sprintf("🟢 Stash is Running on port %d.", c.GetPort()))

	openStashButton := systray.AddMenuItem("Open Stash", "Open a browser window to Stash")
	var menuItems []string
	systray.AddSeparator()
	if !c.IsNewSystem() {
		menuItems = c.GetMenuItems()
		for _, item := range menuItems {
			c := cases.Title(language.Und)
			titleCaseItem := c.String(strings.ToLower(item))
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
				exit <- 0
				return
			}
		}
	}()
}
