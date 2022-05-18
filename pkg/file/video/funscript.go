package video

import (
	"path/filepath"
	"strings"
)

// GetFunscriptPath returns the path of a file
// with the extension changed to .funscript
func GetFunscriptPath(path string) string {
	ext := filepath.Ext(path)
	fn := strings.TrimSuffix(path, ext)
	return fn + ".funscript"
}
