package utils

import (
	"runtime"
	"strings"
)

// FixWindowsPath replaces \ with / in the given path because sometimes the \ isn't recognized as valid on windows
func FixWindowsPath(str string) string {
	if runtime.GOOS == "windows" {
		return strings.ReplaceAll(str, "\\", "/")
	}
	return str
}
