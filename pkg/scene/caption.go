package scene

import (
	"path/filepath"
	"strings"

	"github.com/asticode/go-astisub"
)

// GetCaptionPath generates the path of a caption
// from a given file path wanted language and caption sufffix
func GetCaptionPath(path, lang, suffix string) string {
	ext := filepath.Ext(path)
	fn := strings.TrimSuffix(path, ext)
	captionExt := ""
	if len(lang) == 0 {
		captionExt = suffix
	} else {
		captionExt = "." + suffix
	}
	return fn + "." + lang + captionExt
}

// ReadSubs reads a captions file
func ReadSubs(path string) (*astisub.Subtitles, error) {
	return astisub.OpenFile(path)
}
