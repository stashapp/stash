package image

import (
	"context"
	"errors"
	"fmt"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/file"
)

type DecoratorClip struct {
	FFProbe ffmpeg.FFProbe
}

func (d *DecoratorClip) Decorate(ctx context.Context, fs file.FS, f file.File) (file.File, error) {
	fmt.Println("test")
	if d.FFProbe == "" {
		return f, errors.New("ffprobe not configured")
	}

	base := f.Base()
	// TODO - copy to temp file if not an OsFS
	if _, isOs := fs.(*file.OsFS); !isOs {
		return f, fmt.Errorf("video.constructFile: only OsFS is supported")
	}

	probe := d.FFProbe
	clipFile, err := probe.NewVideoFile(base.Path)
	if err != nil {
		return f, fmt.Errorf("running ffprobe on %q: %w", base.Path, err)
	}

	container, err := ffmpeg.MatchContainer(clipFile.Container, base.Path)
	if err != nil {
		return f, fmt.Errorf("matching container for %q: %w", base.Path, err)
	}

	return &file.ImageFile{
		BaseFile: base,
		Format:   string(container),
		Width:    clipFile.Width,
		Height:   clipFile.Height,
		Clip:     true,
	}, nil
}

func (d *DecoratorClip) IsMissingMetadata(ctx context.Context, fs file.FS, f file.File) bool {
	const (
		unsetString = "unset"
		unsetNumber = -1
	)

	imf, ok := f.(*file.ImageFile)
	if !ok {
		return true
	}

	return imf.Format == unsetString || imf.Width == unsetNumber || imf.Height == unsetNumber
}
