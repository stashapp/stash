package audio

import (
	"context"
	"errors"
	"fmt"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/models"
)

// Decorator adds audio specific fields to a File.
type Decorator struct {
	FFProbe *ffmpeg.FFProbe
}

func (d *Decorator) Decorate(ctx context.Context, fs models.FS, f models.File) (models.File, error) {
	if d.FFProbe == nil {
		return f, errors.New("ffprobe not configured")
	}

	base := f.Base()
	// TODO - copy to temp file if not an OsFS
	if _, isOs := fs.(*file.OsFS); !isOs {
		return f, fmt.Errorf("audio.constructFile: only OsFS is supported")
	}

	probe := d.FFProbe
	audioFile, err := probe.NewAudioFile(base.Path)
	if err != nil {
		return f, fmt.Errorf("running ffprobe on %q: %w", base.Path, err)
	}

	container, err := ffmpeg.MatchContainer(audioFile.Container, base.Path)
	if err != nil {
		return f, fmt.Errorf("matching container for %q: %w", base.Path, err)
	}

	return &models.AudioFile{
		BaseFile:   base,
		Format:     string(container),
		AudioCodec: audioFile.AudioCodec,
		Duration:   audioFile.FileDuration,
		SampleRate: audioFile.SampleRate,
		BitRate:    audioFile.Bitrate,
	}, nil
}

func (d *Decorator) IsMissingMetadata(ctx context.Context, fs models.FS, f models.File) bool {
	return false
}
