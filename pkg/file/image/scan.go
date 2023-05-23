package image

import (
	"context"
	"fmt"
	"image"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/file/video"
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
		fmt.Printf("Warning: File %q could not be read with ffprobe: %s, assuming ImageFile", base.Path, err)
		c, format, err := image.DecodeConfig(r)
		if err != nil {
			return f, fmt.Errorf("decoding image file %q: %w", base.Path, err)
		}
		return &file.ImageFile{
			BaseFile: base,
			Format:   format,
			Width:    c.Width,
			Height:   c.Height,
		}, nil
	}

	isClip := true
	// This list is derived from ffmpegImageThumbnail in pkg/image/thumbnail. If one gets updated, the other should be as well
	for _, item := range []string{"png", "mjpeg", "webp"} {
		if item == probe.VideoCodec {
			isClip = false
		}
	}
	if isClip {
		videoFileDecorator := video.Decorator{FFProbe: d.FFProbe}
		return videoFileDecorator.Decorate(ctx, fs, f)
	}

	return &file.ImageFile{
		BaseFile: base,
		Format:   probe.VideoCodec,
		Width:    probe.Width,
		Height:   probe.Height,
	}, nil
}

func (d *Decorator) IsMissingMetadata(ctx context.Context, fs file.FS, f file.File) bool {
	const (
		unsetString = "unset"
		unsetNumber = -1
	)

	imf, isImage := f.(*file.ImageFile)
	vf, isVideo := f.(*file.VideoFile)

	switch {
	case isImage:
		return imf.Format == unsetString || imf.Width == unsetNumber || imf.Height == unsetNumber
	case isVideo:
		interactive := false
		if _, err := fs.Lstat(video.GetFunscriptPath(vf.Base().Path)); err == nil {
			interactive = true
		}

		return vf.VideoCodec == unsetString || vf.AudioCodec == unsetString ||
			vf.Format == unsetString || vf.Width == unsetNumber ||
			vf.Height == unsetNumber || vf.FrameRate == unsetNumber ||
			vf.Duration == unsetNumber ||
			vf.BitRate == unsetNumber || interactive != vf.Interactive
	default:
		return true
	}
}
