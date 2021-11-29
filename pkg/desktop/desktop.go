package desktop

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/pkg/browser"
	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/utils"
	"golang.org/x/term"
)

func Initialize() {
	if IsDesktop() {
		OpenURLInBrowser(false, "")
		startSystray()
	}
}

// OpenURLInBrowser opens a browser to the Stash UI. Path can be an empty string for main page.
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

	// Only try to update if downloaded from official sources
	if !config.IsOfficialBuild() {
		return false
	}

	// Avoid updating if installed from package manager
	if runtime.GOOS == "linux" {
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
	database.Close()
	os.Exit(0)
}
