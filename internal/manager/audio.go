package manager

import (
	"fmt"
	"net/url"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/models"
)

type AudioStreamEndpoint struct {
	URL      string  `json:"url"`
	MimeType *string `json:"mime_type"`
	Label    *string `json:"label"`
}

const directAudioStreamLabel = "Direct stream"

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

	mimeType := ffmpeg.MimeMp3Audio
	label := directAudioStreamLabel

	return []*AudioStreamEndpoint{
		{
			URL:      directStreamURL.String(),
			MimeType: &mimeType,
			Label:    &label,
		},
	}, nil
}
