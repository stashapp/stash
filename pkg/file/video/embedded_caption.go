package video

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/txn"
	"golang.org/x/text/language"
)

const embeddedCaptionType = "vtt"

// EmbeddedCaptionEncoder runs ffmpeg commands for embedded caption extraction.
type EmbeddedCaptionEncoder interface {
	Generate(ctx context.Context, args ffmpeg.Args) error
}

type embeddedCaptionCandidate struct {
	streamIndex     int
	captionPath     string
	caption         *models.VideoCaption
	needsExtraction bool
}

// ExtractEmbeddedCaptions extracts embedded subtitle streams to VTT sidecar files and registers them as captions.
func ExtractEmbeddedCaptions(ctx context.Context, f *models.VideoFile, encoder EmbeddedCaptionEncoder, txnMgr models.TxnManager, w CaptionUpdater) error {
	if f == nil || len(f.SubtitleStreams) == 0 || encoder == nil || w == nil {
		return nil
	}

	captions, err := getExistingCaptions(ctx, f.ID, txnMgr, w)
	if err != nil {
		return fmt.Errorf("getting captions for file %s: %w", f.Path, err)
	}

	candidates := embeddedCaptionCandidates(f.Path, f.SubtitleStreams, captions, captionFileExists)
	if len(candidates) == 0 {
		return nil
	}

	var added []*models.VideoCaption
	for _, candidate := range candidates {
		if candidate.needsExtraction {
			if err := extractEmbeddedCaption(ctx, encoder, f.Path, candidate.streamIndex, candidate.captionPath); err != nil {
				logger.Warnf("Error extracting embedded caption stream %d from %s: %v", candidate.streamIndex, f.Path, err)
				continue
			}
		}

		added = append(added, candidate.caption)
	}

	if len(added) == 0 {
		return nil
	}

	return updateEmbeddedCaptions(ctx, f.ID, txnMgr, w, added)
}

func embeddedCaptionCandidates(videoPath string, streams []models.VideoSubtitleStream, captions []*models.VideoCaption, exists func(string) bool) []embeddedCaptionCandidate {
	if exists == nil {
		exists = captionFileExists
	}

	seen := make(map[string]struct{})
	for _, caption := range captions {
		seen[captionKey(caption.LanguageCode, caption.CaptionType)] = struct{}{}
	}

	var ret []embeddedCaptionCandidate
	for _, stream := range streams {
		lang := normalizeCaptionLanguage(stream.LanguageCode)
		key := captionKey(lang, embeddedCaptionType)
		if _, ok := seen[key]; ok {
			continue
		}

		captionPath := GetCaptionPath(videoPath, lang, embeddedCaptionType)
		ret = append(ret, embeddedCaptionCandidate{
			streamIndex: stream.Index,
			captionPath: captionPath,
			caption: &models.VideoCaption{
				LanguageCode: lang,
				Filename:     filepath.Base(captionPath),
				CaptionType:  embeddedCaptionType,
			},
			needsExtraction: !exists(captionPath),
		})
		seen[key] = struct{}{}
	}

	return ret
}

func normalizeCaptionLanguage(lang string) string {
	lang = strings.ToLower(strings.TrimSpace(lang))
	switch lang {
	case "", "und", "unknown":
		return LangUnknown
	}

	base, err := language.ParseBase(lang)
	if err != nil {
		return LangUnknown
	}

	return base.String()
}

func captionKey(lang string, captionType string) string {
	return lang + "." + captionType
}

func captionFileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func extractEmbeddedCaption(ctx context.Context, encoder EmbeddedCaptionEncoder, videoPath string, streamIndex int, captionPath string) error {
	tmpPath := fmt.Sprintf("%s.%d.tmp", captionPath, streamIndex)
	defer os.Remove(tmpPath)

	args := ffmpeg.Args{}.
		LogLevel(ffmpeg.LogLevelError).
		Overwrite().
		Input(videoPath)
	args = append(args,
		"-map", fmt.Sprintf("0:%d", streamIndex),
		"-c:s", "webvtt",
		"-f", "webvtt",
		tmpPath,
	)

	if err := encoder.Generate(ctx, args); err != nil {
		return err
	}

	if captionFileExists(captionPath) {
		return nil
	}

	return os.Rename(tmpPath, captionPath)
}

func getExistingCaptions(ctx context.Context, fileID models.FileID, txnMgr models.TxnManager, w CaptionUpdater) ([]*models.VideoCaption, error) {
	var captions []*models.VideoCaption
	err := withDatabase(ctx, txnMgr, func(ctx context.Context) error {
		var err error
		captions, err = w.GetCaptions(ctx, fileID)
		return err
	})
	return captions, err
}

func updateEmbeddedCaptions(ctx context.Context, fileID models.FileID, txnMgr models.TxnManager, w CaptionUpdater, added []*models.VideoCaption) error {
	return withTxn(ctx, txnMgr, func(ctx context.Context) error {
		captions, err := w.GetCaptions(ctx, fileID)
		if err != nil {
			return err
		}

		changed := false
		for _, caption := range added {
			if IsLangInCaptions(caption.LanguageCode, caption.CaptionType, captions) {
				continue
			}

			captions = append(captions, caption)
			changed = true
		}

		if !changed {
			return nil
		}

		return w.UpdateCaptions(ctx, fileID, captions)
	})
}

func withDatabase(ctx context.Context, txnMgr models.TxnManager, fn txn.TxnFunc) error {
	if txnMgr == nil {
		return fn(ctx)
	}

	return txn.WithDatabase(ctx, txnMgr, fn)
}

func withTxn(ctx context.Context, txnMgr models.TxnManager, fn txn.TxnFunc) error {
	if txnMgr == nil {
		return fn(ctx)
	}

	return txn.WithTxn(ctx, txnMgr, fn)
}
