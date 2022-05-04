package scene

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/language"

	"github.com/asticode/go-astisub"
	"github.com/stashapp/stash/pkg/models"
)

var CaptionExts = []string{"vtt", "srt"} // in a case where vtt and srt files are both provided prioritize vtt file due to native support

// to be used for captions without a language code in the filename
// ISO 639-1 uses 2 or 3 a-z chars for codes so 00 is a safe non valid choise
// https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes
const LangUknown = "00"

// GetCaptionPath generates the path of a caption
// from a given file path, wanted language and caption sufffix
func GetCaptionPath(scenePath string, caption *models.SceneCaption) string {
	sceneExt := filepath.Ext(scenePath)
	sceneName := strings.TrimSuffix(scenePath, sceneExt)
	return sceneName + "." + caption.LanguageCode + "." + caption.CaptionType
}

// ReadSubs reads a captions file
func ReadSubs(path string) (*astisub.Subtitles, error) {
	return astisub.OpenFile(path)
}

// IsValidLanguage checks whether the given string is a valid
// ISO 639 language code
func IsValidLanguage(lang string) bool {
	_, err := language.ParseBase(lang)
	return err == nil
}

// IsLangInCaptions returns true if lang is present
// in the captions
func IsLangInCaptions(lang string, ext string, captions []*models.SceneCaption) bool {
	for _, caption := range captions {
		if lang == caption.LanguageCode && ext == caption.CaptionType {
			return true
		}
	}
	return false
}

// GenerateCaptionCandidates generates a list of filenames with exts as extensions
// that can associated with the caption
func GenerateCaptionCandidates(captionPath string, exts []string) []string {
	var candidates []string

	basename := strings.TrimSuffix(captionPath, filepath.Ext(captionPath)) // caption filename without the extension

	// a caption file can be something like scene_filename.srt or scene_filename.en.srt
	// if a language code is present and valid remove it from the basename
	languageExt := filepath.Ext(basename)
	if len(languageExt) > 2 && IsValidLanguage(languageExt[1:]) {
		basename = strings.TrimSuffix(basename, languageExt)
	}

	for _, ext := range exts {
		candidates = append(candidates, basename+"."+ext)
	}

	return candidates
}

// GetCaptionsLangFromPath returns the language code from a given captions path
// If no valid language is present LangUknown is returned
func GetCaptionsLangFromPath(captionPath string) string {
	langCode := LangUknown
	basename := strings.TrimSuffix(captionPath, filepath.Ext(captionPath)) // caption filename without the extension
	languageExt := filepath.Ext(basename)
	if len(languageExt) > 2 && IsValidLanguage(languageExt[1:]) {
		langCode = languageExt[1:]
	}
	return langCode
}

// CleanCaptions removes non existent/accessible language codes from captions
func CleanCaptions(scenePath string, captions []*models.SceneCaption) (cleanedCaptions []*models.SceneCaption, changed bool) {
	changed = false
	for _, caption := range captions {
		found := false
		f := GetCaptionPath(scenePath, caption)
		if _, er := os.Stat(f); er == nil {
			cleanedCaptions = append(cleanedCaptions, caption)
			found = true
		}
		if !found {
			changed = true
		}
	}
	return
}
