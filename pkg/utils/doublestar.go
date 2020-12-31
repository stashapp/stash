package utils

import (
	"strings"
)

// EscapeGlobPath escapes doublestar glob special terms from the path string
// so that a glob match pattern can be safely applied
func EscapeGlobPath(path string) string {

	// doublestar glob needs special terms escaped with a backslash ( \ )
	specialEscaped := []string{"*", "\\*", "?", "\\?", "[", "\\[", "]", "\\]", "{", "\\{", "}", "\\}"}

	r := strings.NewReplacer(specialEscaped...)
	escaped := r.Replace(path)
	return escaped
}
