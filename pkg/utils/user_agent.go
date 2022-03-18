package utils

import "runtime"

// valid UA from https://user-agents.net
const Safari = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Safari/605.1.15/iY0wnXbs-59"
const FirefoxWindows = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:88.0) Gecko/20100101 Firefox/88.0"
const FirefoxLinux = "Mozilla/5.0 (X11; Linux x86_64; rv:88.0) Gecko/20100101 Firefox/88.0"
const FirefoxLinuxArm = "Mozilla/5.0 (X11; Linux armv7l; rv:86.0) Gecko/20100101 Firefox/86.0"
const FirefoxLinuxArm64 = "Mozilla/5.0 (X11; Linux aarch64; rv:86.0) Gecko/20100101 Firefox/86.0"

// getUserAgent returns a valid User Agent string that matches the running os/arch
func getUserAgent() string {
	arch := runtime.GOARCH
	os := runtime.GOOS

	switch os {
	case "darwin":
		return Safari
	case "windows":
		return FirefoxWindows
	case "linux":
		switch arch {
		case "arm":
			return FirefoxLinuxArm
		case "arm64":
			return FirefoxLinuxArm64
		case "amd64":
			return FirefoxLinux
		default:
			return FirefoxLinux
		}
	default:
		return FirefoxLinux
	}
}
