package desktop

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/pkg/browser"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"golang.org/x/term"
)

type ShutdownHandler interface {
	Shutdown(code int)
}

type FaviconProvider interface {
	GetFavicon() []byte
	GetFaviconPng() []byte
}

// Start starts the desktop icon process. It blocks until the process exits.
// MUST be run on the main goroutine or will have no effect on macOS
func Start(shutdownHandler ShutdownHandler, faviconProvider FaviconProvider) {
	if IsDesktop() {
		c := config.GetInstance()
		if !c.GetNoBrowser() {
			openURLInBrowser("")
		}
		writeStashIcon(faviconProvider)
		startSystray(shutdownHandler, faviconProvider)
	}
}

// openURLInBrowser opens a browser to the Stash UI. Path can be an empty string for main page.
func openURLInBrowser(path string) {
	// This can be done before actually starting the server, as modern browsers will
	// automatically reload the page if a local port is closed at page load and then opened.
	serverAddress := getServerURL(path)

	err := browser.OpenURL(serverAddress)
	if err != nil {
		logger.Error("Could not open browser: " + err.Error())
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
	// Check if running under root
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

// writeStashIcon writes the current stash logo to config/icon.png
func writeStashIcon(faviconProvider FaviconProvider) {
	c := config.GetInstance()
	if !c.IsNewSystem() {
		iconPath := path.Join(c.GetConfigPath(), "icon.png")
		err := ioutil.WriteFile(iconPath, faviconProvider.GetFaviconPng(), 0644)
		if err != nil {
			logger.Errorf("Couldn't write icon file: %s", err.Error())
		}
	}
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
		if fsutil.IsPathInDir("/usr", executablePath) || fsutil.IsPathInDir("/opt", executablePath) {
			return false
		}

		if isServerDockerized() {
			return false
		}
	}

	return true
}

func getIconPath() string {
	return path.Join(config.GetInstance().GetConfigPath(), "icon.png")
}

func RevealInFileManager(path string) {
	exists, err := fsutil.FileExists(path)
	if err != nil {
		logger.Errorf("Error checking file: %s", err)
		return
	}
	if exists && IsDesktop() {
		revealInFileManager(path)
	}
}

func getServerURL(path string) string {
	c := config.GetInstance()
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

	return serverAddress
}
