package video

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	stashExec "github.com/stashapp/stash/pkg/exec"
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

var supportedSubtitleCodecs = map[string]struct{}{
	"subrip":   {},
	"srt":      {},
	"ass":      {},
	"ssa":      {},
	"webvtt":   {},
	"mov_text": {},
}

// ExtractEmbeddedCaptions extracts embedded subtitle streams from a video file
// to WebVTT sidecar files and registers them via the provided CaptionUpdater.
// Streams with unsupported codecs or that already have a registered VTT caption
// for the same language are silently skipped.
func ExtractEmbeddedCaptions(ctx context.Context, videoPath string, streams []*ffmpeg.FFProbeStream, ffmpegPath string, fileID models.FileID, w CaptionUpdater) error {
	captions, err := w.GetCaptions(ctx, fileID)
	if err != nil {
		return fmt.Errorf("getting captions for %s: %w", videoPath, err)
	}

	changed := false
	for _, s := range streams {
		if _, ok := supportedSubtitleCodecs[s.CodecName]; !ok {
			logger.Debugf("[embedded captions] unsupported codec %q in %s, skipping stream %d", s.CodecName, videoPath, s.Index)
			continue
		}

		lang := LangUnknown
		if IsValidLanguage(s.Tags.Language) {
			lang = s.Tags.Language
		}

		if IsLangInCaptions(lang, "vtt", captions) {
			logger.Debugf("[embedded captions] lang=%s vtt already registered for %s", lang, videoPath)
			continue
		}

		vttPath := GetCaptionPath(videoPath, lang, "vtt")

		if _, statErr := os.Stat(vttPath); statErr != nil {
			if err := extractSubtitleStream(ctx, ffmpegPath, videoPath, s.Index, vttPath); err != nil {
				logger.Warnf("[embedded captions] failed to extract stream %d from %s: %v", s.Index, videoPath, err)
				continue
			}
		}

		captions = append(captions, &models.VideoCaption{
			LanguageCode: lang,
			Filename:     filepath.Base(vttPath),
			CaptionType:  "vtt",
		})
		changed = true
		logger.Infof("[embedded captions] extracted subtitle stream (lang=%s, index=%d) from %s", lang, s.Index, videoPath)
	}

	if changed {
		return w.UpdateCaptions(ctx, fileID, captions)
	}
	return nil
}

func extractSubtitleStream(ctx context.Context, ffmpegPath, videoPath string, streamIndex int, outputPath string) error {
	args := []string{
		"-v", "quiet",
		"-i", videoPath,
		"-map", fmt.Sprintf("0:%d", streamIndex),
		"-c:s", "webvtt",
		outputPath,
	}
	cmd := stashExec.CommandContext(ctx, ffmpegPath, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg exit: %w — output: %s", err, string(out))
	}
	return nil
}
