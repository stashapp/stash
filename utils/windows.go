package utils

import (
	"runtime"
	"strings"
)

// Sometimes the \ isn't recognized as valid on windows
func FixWindowsPath(str string) string {
	if runtime.GOOS == "windows" {
		return strings.Replace(str, "\\", "/", -1)
	} else {
		return str
	}
}