package video

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/models"
	"golang.org/x/text/language"
)

const embeddedCaptionOutputType = "srt"

var supportedEmbeddedCaptionCodecs = map[string]struct{}{
	"ass":      {},
	"mov_text": {},
	"ssa":      {},
	"subrip":   {},
	"webvtt":   {},
}

type embeddedCaptionProbe interface {
	NewVideoFile(videoPath string) (*ffmpeg.VideoFile, error)
}

type embeddedCaptionEncoder interface {
	Generate(ctx context.Context, args ffmpeg.Args) error
}

type embeddedCaptionTrack struct {
	streamIndex int
	language    string
	captionType string
}

// EmbeddedCaptionExtractor extracts supported embedded text subtitle streams
// into sidecar caption files that can be handled by the existing caption flow.
type EmbeddedCaptionExtractor struct {
	FFProbe embeddedCaptionProbe
	FFMpeg  embeddedCaptionEncoder
}

func (e *EmbeddedCaptionExtractor) Extract(ctx context.Context, videoFile *models.VideoFile) ([]string, error) {
	if e == nil || e.FFProbe == nil || e.FFMpeg == nil {
		return nil, errors.New("embedded caption extractor is not configured")
	}
	if videoFile == nil {
		return nil, errors.New("video file is nil")
	}
	if videoFile.ZipFileID != nil {
		return nil, nil
	}

	probed, err := e.FFProbe.NewVideoFile(videoFile.Path)
	if err != nil {
		return nil, fmt.Errorf("probing video captions: %w", err)
	}

	var ret []string
	for _, track := range embeddedCaptionTracks(probed.JSON.Streams) {
		captionPath := GetCaptionPath(videoFile.Path, track.language, track.captionType)
		exists, err := captionFileExists(captionPath)
		if err != nil {
			return ret, err
		}
		if !exists {
			if err := e.extractTrack(ctx, videoFile.Path, track, captionPath); err != nil {
				_ = os.Remove(captionPath)
				return ret, err
			}
		}

		ret = append(ret, captionPath)
	}

	return ret, nil
}

func (e *EmbeddedCaptionExtractor) extractTrack(ctx context.Context, videoPath string, track embeddedCaptionTrack, outputPath string) error {
	args := ffmpeg.Args{
		"-v", "error",
		"-nostdin",
		"-i", videoPath,
		"-map", fmt.Sprintf("0:%d", track.streamIndex),
		"-vn",
		"-an",
		"-c:s", track.captionType,
		outputPath,
	}

	if err := e.FFMpeg.Generate(ctx, args); err != nil {
		return fmt.Errorf("extracting embedded caption stream %d: %w", track.streamIndex, err)
	}

	return nil
}

func embeddedCaptionTracks(streams []ffmpeg.FFProbeStream) []embeddedCaptionTrack {
	seen := make(map[string]struct{})
	var ret []embeddedCaptionTrack

	for _, stream := range streams {
		if stream.CodecType != "subtitle" || !isSupportedEmbeddedCaptionCodec(stream.CodecName) {
			continue
		}

		lang := normalizeEmbeddedCaptionLanguage(stream.Tags.Language)
		captionType := embeddedCaptionOutputType
		key := lang + "." + captionType
		if _, found := seen[key]; found {
			continue
		}
		seen[key] = struct{}{}

		ret = append(ret, embeddedCaptionTrack{
			streamIndex: stream.Index,
			language:    lang,
			captionType: captionType,
		})
	}

	return ret
}

func isSupportedEmbeddedCaptionCodec(codec string) bool {
	_, ok := supportedEmbeddedCaptionCodecs[strings.ToLower(codec)]
	return ok
}

func normalizeEmbeddedCaptionLanguage(lang string) string {
	lang = strings.TrimSpace(strings.ToLower(lang))
	if lang == "" || lang == "und" {
		return LangUnknown
	}

	if tag, err := language.Parse(lang); err == nil {
		if base, confidence := tag.Base(); confidence != language.No && base.String() != "und" {
			return base.String()
		}
	}

	if base, err := language.ParseBase(lang); err == nil && base.String() != "und" {
		return base.String()
	}

	return LangUnknown
}

func captionFileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return !info.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, fmt.Errorf("checking caption file %q: %w", path, err)
}
