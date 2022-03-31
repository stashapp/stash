package scene

import (
	"path/filepath"
	"strings"
)

func GetCaptionPath(path string) string {
	ext := filepath.Ext(path)
	fn := strings.TrimSuffix(path, ext)
	return fn + ".en.vtt"
}
