package desktop

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/getlantern/systray"
	"github.com/pkg/browser"
	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/utils"
)

func Initialize() {
	if IsDesktop() {
		OpenURLInBrowser(false, "")
		c := config.GetInstance()
		go systray.Run(systrayInitialize, nil)

		// Shows a small notification to inform that Stash will no longer show a terminal window,
		// and instead will be available in the tray. Will only show the first time a pre-desktop integration
		// system is started from a non-terminal method, e.g. double-clicking an icon.
		if c.GetShowOneTimeMovedNotification() {
			SendNotification("Stash has moved!", "Stash now runs in your tray, instead of a terminal window.")
			c.Set(config.ShowOneTimeMovedNotification, false)
			if err := c.Write(); err != nil {
				logger.Errorf("Error while writing configuration file: %s", err.Error())
			}
		}
	}
}

func systrayInitialize() {
	systray.SetTemplateIcon(favicon, favicon)
	systray.SetTitle("Stash")
	systray.SetTooltip("🟢 Stash is Running.")

	openStashButton := systray.AddMenuItem("Open Stash", "Open a browser window to Stash")
	var menuItems []string
	systray.AddSeparator()
	c := config.GetInstance()
	if !c.IsNewSystem() {
		menuItems = c.GetMenuItems()
		for _, item := range menuItems {
			curr := systray.AddMenuItem(strings.Title(strings.ToLower(item)), "Open to "+item)
			go func(item string) {
				for {
					<-curr.ClickedCh
					if item == "markers" {
						item = "scenes/markers"
					}
					OpenURLInBrowser(true, item)
				}
			}(item)
		}
		systray.AddSeparator()
		// systray.AddMenuItem("Start a Scan", "Scan all libraries with default settings")
		// systray.AddMenuItem("Start Auto Tagging", "Auto Tag all libraries")
		// systray.AddSeparator()
	}

	quitStashButton := systray.AddMenuItem("Quit Stash Server", "Quits the Stash server")

	go func() {
		for {
			select {
			case <-openStashButton.ClickedCh:
				OpenURLInBrowser(true, "")
			case <-quitStashButton.ClickedCh:
				Shutdown()
			}
		}
	}()
}

// OpenURLInBrowser opens a browser to the Stash UI. Path can be an empty string.
func OpenURLInBrowser(force bool, path string) {
	// This can be done before actually starting the server, as modern browsers will
	// automatically reload the page if a local port is closed at page load and then opened.
	c := config.GetInstance()
	if force || (!c.GetNoBrowser() && IsDesktop()) {
		serverAddress := c.GetHost()
		if serverAddress == "0.0.0.0" {
			serverAddress = "localhost"
		}
		serverAddress = serverAddress + ":" + strconv.Itoa(c.GetPort())

		proto := ""
		if c.HasTLSConfig() {
			proto = "https://"
		} else {
			proto = "http://"
		}
		serverAddress = proto + serverAddress + "/"

		if path != "" {
			serverAddress += strings.TrimPrefix(path, "/")
		}

		err := browser.OpenURL(serverAddress)
		if err != nil {
			logger.Error("Could not open browser: " + err.Error())
		}
	}
}

func SendNotification(title string, text string) {
	if IsDesktop() {
		c := config.GetInstance()
		if c.GetNotificationsEnabled() {
			sendNotification(title, text)
		}
	}
}

func IsDesktop() bool {
	// check if running under root
	if os.Getuid() == 0 {
		return false
	}
	// // Check if stdin is a terminal
	// if term.IsTerminal(int(os.Stdin.Fd())) {
	// 	return false
	// }
	if isService() {
		return false
	}
	if IsServerDockerized() {
		return false
	}

	return true
}

func IsServerDockerized() bool {
	return isServerDockerized()
}

// IsAllowedAutoUpdate tries to determine if the stash binary was installed from a
// package manager or if touching the executable is otherwise a bad idea
func IsAllowedAutoUpdate() bool {

	executablePath, err := os.Executable()
	if err != nil {
		logger.Errorf("Cannot get executable path: %s", err)
		return false
	}
	executablePath, err = filepath.EvalSymlinks(executablePath)
	if err != nil {
		logger.Errorf("Cannot get executable path: %s", err)
		return false
	}

	if runtime.GOOS == "linux" {
		if utils.IsPathInDir("/usr", executablePath) || utils.IsPathInDir("/opt", executablePath) {
			return false
		}
		if isServerDockerized() {
			return false
		}
	}

	return true
}

func Shutdown() {
	systray.Quit()
	database.Close()
	os.Exit(0)
}
