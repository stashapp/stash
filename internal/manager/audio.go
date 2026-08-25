package manager

import (
	"fmt"
	"mime"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/models"
)

type AudioStreamEndpoint struct {
	URL      string  `json:"url"`
	MimeType *string `json:"mime_type"`
	Label    *string `json:"label"`
}

const directAudioStreamLabel = "Direct stream"

// audioContainerMimeTypes maps the container reported by ffprobe (see
// ffmpeg.MatchContainer) to the mime type used to serve the direct stream.
var audioContainerMimeTypes = map[string]string{
	"mp3":      ffmpeg.MimeMp3Audio,
	"mp4":      ffmpeg.MimeMp4Audio,
	"m4v":      ffmpeg.MimeMp4Audio,
	"mov":      ffmpeg.MimeMp4Audio,
	"flac":     ffmpeg.MimeFlacAudio,
	"ogg":      ffmpeg.MimeOggAudio,
	"wav":      ffmpeg.MimeWavAudio,
	"aac":      ffmpeg.MimeAacAudio,
	"webm":     ffmpeg.MimeWebmAudio,
	"matroska": ffmpeg.MimeMkvAudio,
}

// audioMimeType returns the mime type for the audio file's container, falling
// back to the extension and finally to mpeg.
func audioMimeType(f *models.AudioFile) string {
	if mimeType, found := audioContainerMimeTypes[strings.ToLower(f.Format)]; found {
		return mimeType
	}

	// fall back to the extension - covers containers ffprobe reports under a
	// name we don't recognise
	if mimeType := mime.TypeByExtension(strings.ToLower(filepath.Ext(f.Base().Path))); mimeType != "" {
		return mimeType
	}

	return ffmpeg.MimeMp3Audio
}

// GetAudioStreamPaths returns the stream endpoints for the audio. Audio files are
// only ever served directly, so this is always the single direct stream endpoint.
func GetAudioStreamPaths(audio *models.Audio, directStreamURL *url.URL) ([]*AudioStreamEndpoint, error) {
	if audio == nil {
		return nil, fmt.Errorf("nil audio")
	}

	pf := audio.Files.Primary()
	if pf == nil {
		return nil, nil
	}

	mimeType := audioMimeType(pf)
	label := directAudioStreamLabel

	return []*AudioStreamEndpoint{
		{
			URL:      directStreamURL.String(),
			MimeType: &mimeType,
			Label:    &label,
		},
	}, nil
}
