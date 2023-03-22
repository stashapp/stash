package image

import (
	"context"
	"fmt"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/file"
	_ "golang.org/x/image/webp"
)

// Decorator adds image specific fields to a File.
type Decorator struct {
	FFProbe ffmpeg.FFProbe
}

func (d *Decorator) Decorate(ctx context.Context, fs file.FS, f file.File) (file.File, error) {
	base := f.Base()
	r, err := fs.Open(base.Path)
	if err != nil {
		return f, fmt.Errorf("reading image file %q: %w", base.Path, err)
	}
	defer r.Close()

	probe, err := d.FFProbe.NewVideoFile(base.Path)
	if err != nil {
		return f, fmt.Errorf("reading image with ffprobe %q: %w", base.Path, err)
	}
	isClip := true
	// This list is derived from ffmpegImageThumbnail in pkg/image/thumbnail. If one gets updated, the other should be as well
	for _, item := range []string{"png", "mjpeg", "webp"} {
		if item == probe.VideoCodec {
			isClip = false
		}
	}

	return &file.ImageFile{
		BaseFile: base,
		Format:   probe.VideoCodec,
		Width:    probe.Width,
		Height:   probe.Height,
		Clip:     isClip,
	}, nil
}

func (d *Decorator) IsMissingMetadata(ctx context.Context, fs file.FS, f file.File) bool {
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
