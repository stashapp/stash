package video

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/asticode/go-astisub"
	"github.com/stashapp/stash/pkg/ffmpeg"
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

const EmbeddedCaptionType = "embedded"

// iso639_2BTo1 maps ISO 639-2/B bibliographic codes commonly emitted by ffprobe
// (e.g. from MKV containers) to their ISO 639-1 equivalents.
var iso639_2BTo1 = map[string]string{
	"alb": "sq", "arm": "hy", "baq": "eu", "bur": "my", "chi": "zh",
	"cze": "cs", "dut": "nl", "fre": "fr", "geo": "ka", "ger": "de",
	"gre": "el", "ice": "is", "mac": "mk", "mao": "mi", "may": "ms",
	"per": "fa", "rum": "ro", "slo": "sk", "tib": "bo", "wel": "cy",
}

// NormalizeSubtitleLanguage normalizes an ffprobe subtitle language tag
// (e.g. ISO 639-2 "eng", "jpn", "fre", "und") into an ISO 639-1/base code (e.g. "en", "ja", "fr")
// or LangUnknown ("00") when absent or undetermined.
func NormalizeSubtitleLanguage(raw string) string {
	raw = strings.ToLower(strings.TrimSpace(raw))
	if raw == "" || raw == "und" || raw == LangUnknown {
		return LangUnknown
	}
	if mapped, ok := iso639_2BTo1[raw]; ok {
		return mapped
	}
	base, err := language.ParseBase(raw)
	if err != nil {
		return LangUnknown
	}
	code := base.String()
	if code == "" || code == "und" {
		return LangUnknown
	}
	return code
}

// EmbeddedCaptionsFromStreams converts ffprobe text subtitle streams into
// deduplicated VideoCaption entries suitable for storing in video_captions.
func EmbeddedCaptionsFromStreams(streams []*ffmpeg.FFProbeStream) []*models.VideoCaption {
	captions := make([]*models.VideoCaption, 0, len(streams))
	langCounts := make(map[string]int, len(streams))

	for _, stream := range streams {
		if stream == nil || !ffmpeg.IsTextSubtitleCodec(stream.CodecName) {
			continue
		}
		lang := NormalizeSubtitleLanguage(stream.Tags.Language)
		langCounts[lang]++
		captionType := EmbeddedCaptionType
		if langCounts[lang] > 1 {
			captionType = fmt.Sprintf("%s-%d", EmbeddedCaptionType, langCounts[lang])
		}
		captions = append(captions, models.NewEmbeddedVideoCaption(lang, stream.Index, captionType))
	}

	return captions
}

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
			if !IsLangInCaptions(captionLang, ext, captions) { // only update captions if language code is not present
				newCaption := &models.VideoCaption{
					LanguageCode: captionLang,
					Filename:     filepath.Base(captionPath),
					CaptionType:  ext,
				}
				captions = append(captions, newCaption)
				er = w.UpdateCaptions(ctx, fileID, captions)
				if er != nil {
					return fmt.Errorf("updating captions for file %s: %w", path, er)
				}

				logger.Debugf("Updated captions for file %s. Added %s", path, captionLang)
			}
		}
		return err
	}); err != nil {
		logger.Error(err.Error())
	}

	return matched
}

// CleanCaptions removes non existent/accessible language codes from captions
func CleanCaptions(ctx context.Context, f *models.VideoFile, txnMgr txn.Manager, w CaptionUpdater) error {
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
		if caption.IsEmbedded() {
			newCaptions = append(newCaptions, caption)
			continue
		}
		captionPath := caption.Path(filePath)
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

// SyncEmbeddedCaptions merges the embedded subtitle tracks discovered on f with
// any existing sidecar captions in storage and persists updates when changed.
func SyncEmbeddedCaptions(ctx context.Context, f *models.VideoFile, txnMgr txn.Manager, w CaptionUpdater) error {
	if f == nil || f.EmbeddedCaptions == nil {
		return nil
	}

	existing, err := w.GetCaptions(ctx, f.ID)
	if err != nil {
		return fmt.Errorf("getting captions for file %s: %w", f.Path, err)
	}

	var merged []*models.VideoCaption
	var existingEmbedded []*models.VideoCaption
	for _, c := range existing {
		if c.IsEmbedded() {
			existingEmbedded = append(existingEmbedded, c)
		} else {
			merged = append(merged, c)
		}
	}

	if captionsEqual(existingEmbedded, f.EmbeddedCaptions) {
		return nil
	}

	merged = append(merged, f.EmbeddedCaptions...)

	fn := func(ctx context.Context) error {
		return w.UpdateCaptions(ctx, f.ID, merged)
	}

	if txnMgr == nil {
		err = fn(ctx)
	} else {
		err = txn.WithTxn(ctx, txnMgr, fn)
	}
	if err != nil {
		return fmt.Errorf("updating embedded captions for file %s: %w", f.Path, err)
	}

	return nil
}

func captionsEqual(a, b []*models.VideoCaption) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].LanguageCode != b[i].LanguageCode ||
			a[i].Filename != b[i].Filename ||
			a[i].CaptionType != b[i].CaptionType {
			return false
		}
	}
	return true
}
