package scene

import (
	"path/filepath"
	"strings"

	"github.com/asticode/go-astisub"
)

// GetCaptionPath generates the path of a subtitle
// from a given file path wanted language and subtitle sufffix
func GetCaptionPath(path, lang, suffix string) string {
	ext := filepath.Ext(path)
	fn := strings.TrimSuffix(path, ext)
	return fn + "." + lang + "." + suffix
}

// ReadSubs reads a subtitles file
func ReadSubs(path string) (*astisub.Subtitles, error) {
	return astisub.OpenFile(path)
}