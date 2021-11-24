package desktop

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/getlantern/systray"
	"github.com/pkg/browser"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/utils"
	"golang.org/x/term"
)

func Initialize() {
	if IsDesktop() {
		OpenURLInBrowser(false, "")
		systray.Run(systrayOnReady, nil)
	}

	// we should re-render the systray if the system if finalized, or if an update is available
	// go func() {
	// 	// TODO config channel listen
	// 	systray.Quit()
	// 	go systray.Run(systrayOnReady, nil)
	// }()
}

func systrayOnReady() {
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

		testnotify := systray.AddMenuItem("Test Notification", "Send a test notification")
		go func(testCh chan struct{}) {
			for {
				<-testCh
			}
		}(testnotify.ClickedCh)
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

// OpenURLInBrowser opens a browser to the Stash UI. Path is optional.
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

func IsDesktop() bool {
	// check if running under root
	if os.Getuid() == 0 {
		return false
	}
	// Check if stdin is a terminal
	if term.IsTerminal(int(os.Stdin.Fd())) {
		return false
	}
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
	manager.GetInstance().Shutdown(0)
}
