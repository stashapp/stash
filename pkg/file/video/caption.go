package video

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/asticode/go-astisub"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
	"golang.org/x/text/language"
)

var CaptionExts = []string{"vtt", "srt"} // in a case where vtt and srt files are both provided prioritize vtt file due to native support

// to be used for captions without a language code in the filename
// ISO 639-1 uses 2 or 3 a-z chars for codes so 00 is a safe non valid choise
// https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes
const LangUnknown = "00"

// GetCaptionPath generates the path of a caption
// from a given file path, wanted language and caption sufffix
func GetCaptionPath(path, lang, suffix string) string {
	ext := filepath.Ext(path)
	fn := strings.TrimSuffix(path, ext)
	captionExt := ""
	if len(lang) == 0 || lang == LangUnknown {
		captionExt = suffix
	} else {
		captionExt = lang + "." + suffix
	}
	return fn + "." + captionExt
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
func IsLangInCaptions(lang string, ext string, captions []*models.VideoCaption) bool {
	for _, caption := range captions {
		if lang == caption.LanguageCode && ext == caption.CaptionType {
			return true
		}
	}
	return false
}

func captionIndex(lang string, ext string, captions []*models.VideoCaption) int {
	for i, caption := range captions {
		if lang == caption.LanguageCode && ext == caption.CaptionType {
			return i
		}
	}

	return -1
}

// getCaptionPrefix returns the prefix used to search for video files for the provided caption path
func getCaptionPrefix(captionPath string) string {
	basename := strings.TrimSuffix(captionPath, filepath.Ext(captionPath)) // caption filename without the extension

	// a caption file can be something like scene_filename.srt or scene_filename.en.srt
	// if a language code is present and valid remove it from the basename
	languageExt := filepath.Ext(basename)
	if len(languageExt) > 2 && IsValidLanguage(languageExt[1:]) {
		basename = strings.TrimSuffix(basename, languageExt)
	}

	return basename + "."
}

// GetCaptionsLangFromPath returns the language code from a given captions path
// If no valid language is present LangUknown is returned
func getCaptionsLangFromPath(captionPath string) string {
	langCode := LangUnknown
	basename := strings.TrimSuffix(captionPath, filepath.Ext(captionPath)) // caption filename without the extension
	languageExt := filepath.Ext(basename)
	if len(languageExt) > 2 && IsValidLanguage(languageExt[1:]) {
		langCode = languageExt[1:]
	}
	return langCode
}

type CaptionUpdater interface {
	GetCaptions(ctx context.Context, fileID models.FileID) ([]*models.VideoCaption, error)
	UpdateCaptions(ctx context.Context, fileID models.FileID, captions []*models.VideoCaption) error
}

type GeneratedCaptionPathResolver interface {
	GetGeneratedCaptionPath(filename string) string
}

func GetCaptionFilePath(caption *models.VideoCaption, videoPath string, generatedCaptionPathResolver GeneratedCaptionPathResolver) string {
	if caption.Generated && generatedCaptionPathResolver != nil {
		return generatedCaptionPathResolver.GetGeneratedCaptionPath(caption.Filename)
	}

	return caption.Path(videoPath)
}

func MergeGeneratedCaptions(captions []*models.VideoCaption, generatedCaptions []*models.VideoCaption) []*models.VideoCaption {
	ret := make([]*models.VideoCaption, len(captions))
	copy(ret, captions)

	for _, generated := range generatedCaptions {
		replaced := false
		for i, caption := range ret {
			if caption.LanguageCode == generated.LanguageCode && caption.CaptionType == generated.CaptionType {
				if caption.Generated {
					ret[i] = generated
				}
				replaced = true
				break
			}
		}

		if !replaced {
			ret = append(ret, generated)
		}
	}

	return ret
}

// MatchesCaption returns true if the caption file matches the video file based on the filename
func MatchesCaption(videoPath, captionPath string) bool {
	captionPrefix := getCaptionPrefix(captionPath)
	videoPrefix := strings.TrimSuffix(videoPath, filepath.Ext(videoPath)) + "."
	return captionPrefix == videoPrefix
}

// associates captions to scene/s with the same basename
// returns true if the caption file was matched to a video file and processed, false otherwise
func AssociateCaptions(ctx context.Context, captionPath string, txnMgr txn.Manager, fqb models.FileFinder, w CaptionUpdater) bool {
	captionLang := getCaptionsLangFromPath(captionPath)

	captionPrefix := getCaptionPrefix(captionPath)
	matched := false
	if err := txn.WithTxn(ctx, txnMgr, func(ctx context.Context) error {
		var err error
		files, er := fqb.FindAllByPath(ctx, captionPrefix+"*", true)

		if er != nil {
			return fmt.Errorf("searching for scene %s: %w", captionPrefix, er)
		}

		for _, f := range files {
			// found some files
			// filter out non video files
			switch f.(type) {
			case *models.VideoFile:
				break
			default:
				continue
			}

			fileID := f.Base().ID
			path := f.Base().Path

			logger.Debugf("Matched captions to file %s", path)
			matched = true

			captions, er := w.GetCaptions(ctx, fileID)
			if er != nil {
				return fmt.Errorf("getting captions for file %s: %w", path, er)
			}

			fileExt := filepath.Ext(captionPath)
			ext := fileExt[1:]
			newCaption := &models.VideoCaption{
				LanguageCode: captionLang,
				Filename:     filepath.Base(captionPath),
				CaptionType:  ext,
			}

			existingIndex := captionIndex(captionLang, ext, captions)
			if existingIndex == -1 {
				captions = append(captions, newCaption)
			} else if captions[existingIndex].Generated {
				captions[existingIndex] = newCaption
			} else {
				continue
			}

			er = w.UpdateCaptions(ctx, fileID, captions)
			if er != nil {
				return fmt.Errorf("updating captions for file %s: %w", path, er)
			}

			logger.Debugf("Updated captions for file %s. Added %s", path, captionLang)
		}
		return err
	}); err != nil {
		logger.Error(err.Error())
	}

	return matched
}

// CleanCaptions removes non existent/accessible language codes from captions
func CleanCaptions(ctx context.Context, f *models.VideoFile, txnMgr txn.Manager, w CaptionUpdater, generatedCaptionPathResolver GeneratedCaptionPathResolver) error {
	captions, err := w.GetCaptions(ctx, f.ID)
	if err != nil {
		return fmt.Errorf("getting captions for file %s: %w", f.Path, err)
	}

	if len(captions) == 0 {
		return nil
	}

	filePath := f.Path

	changed := false
	var newCaptions []*models.VideoCaption

	for _, caption := range captions {
		captionPath := GetCaptionFilePath(caption, filePath, generatedCaptionPathResolver)
		_, err := os.Stat(captionPath)
		if errors.Is(err, os.ErrNotExist) {
			logger.Infof("Removing non existent caption %s for %s", caption.Filename, f.Path)
			changed = true
		} else {
			// other errors are ignored for the purposes of cleaning
			newCaptions = append(newCaptions, caption)
		}
	}

	if changed {
		fn := func(ctx context.Context) error {
			return w.UpdateCaptions(ctx, f.ID, newCaptions)
		}

		// possible that we are already in a transaction and txnMgr is nil
		// in that case just call the function directly
		if txnMgr == nil {
			err = fn(ctx)
		} else {
			err = txn.WithTxn(ctx, txnMgr, fn)
		}

		if err != nil {
			return fmt.Errorf("updating captions for file %s: %w", f.Path, err)
		}
	}

	return nil
}
