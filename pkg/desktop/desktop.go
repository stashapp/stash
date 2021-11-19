package desktop

import (
	"io/ioutil"
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
)

func Initialize() {
	if IsDesktop() {
		StartBrowser(false)
		go systray.Run(systrayOnReady, nil)
	}

	// we should re-render the systray if the system if finalized, or if an update is available
	// go func() {
	// 	// TODO config channel listen
	// 	systray.Quit()
	// 	go systray.Run(systrayOnReady, nil)
	// }()
}

func systrayOnReady() {
	systray.SetTitle("Stash")
	systray.SetTooltip("🟢 Stash is Running.")

	openStashButton := systray.AddMenuItem("Open Stash", "Open a browser window to Stash")
	var menuItems []string
	systray.AddSeparator()
	c := config.GetInstance()
	if !c.IsNewSystem() {
		menuItems = c.GetMenuItems()
		for _, item := range menuItems {
			systray.AddMenuItem(strings.Title(strings.ToLower(item)), "Open to "+item)
			// TODO add handlers
		}
		systray.AddSeparator()
		systray.AddMenuItem("Start a Scan", "Scan all libraries with default settings")
		systray.AddMenuItem("Start Auto Tagging", "Auto Tag all libraries")
		systray.AddSeparator()
	}

	quitStashButton := systray.AddMenuItem("Quit Stash Server", "Quits the Stash server")

	go func() {
		for {
			select {
			case <-openStashButton.ClickedCh:
				StartBrowser(true)
			case <-quitStashButton.ClickedCh:
				manager.GetInstance().Shutdown()
				os.Exit(0)
			}
		}
	}()
}

func StartBrowser(force bool) {
	// This can be done before actually starting the server, as modern browsers will
	// automatically reload the page if a local port is closed at page load and then opened.
	c := config.GetInstance()
	if force || (!c.GetNoBrowser() && IsDesktop()) {
		displayHost := c.GetHost()
		if displayHost == "0.0.0.0" {
			displayHost = "localhost"
		}
		displayAddress := displayHost + ":" + strconv.Itoa(c.GetPort())
		if c.HasTLSConfig() {
			displayAddress = "https://" + displayAddress + "/"
		} else {
			displayAddress = "http://" + displayAddress + "/"
		}

		err := browser.OpenURL(displayAddress)
		if err != nil {
			logger.Error("Could not open browser: " + err.Error())
		}
	}
}

func IsDesktop() bool {
	// TODO if desktop integration is disabled

	// check if running under root
	if os.Getuid() == 0 {
		return false
	}
	// check if started by init, e.g. stash is a *nix systemd service / MacOS launchd service
	if os.Getppid() == 1 {
		return false
	}
	if IsServerDockerized() {
		return false
	}
	// TODO: Check if stdin is terminal
	// PERHAPS: Check if windows service

	return true
}

func IsServerDockerized() bool {
	if runtime.GOOS == "linux" {
		_, dockerEnvErr := os.Stat("/.dockerenv")
		cgroups, _ := ioutil.ReadFile("/proc/self/cgroup")
		if os.IsExist(dockerEnvErr) || strings.Contains(string(cgroups), "docker") {
			return true
		}
	}

	return false
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
	if !utils.IsPathInDir("/usr", executablePath) {
		if !utils.IsPathInDir("/opt", executablePath) {
			return !IsServerDockerized()
		}
	}
	return false
}
