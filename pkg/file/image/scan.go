package image

import (
	"context"
	"errors"
	"fmt"
	"image"
	"io"
	"os"
	"path/filepath"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/file/video"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	_ "golang.org/x/image/webp"
)

var ErrUnsupportedAVIFInZip = errors.New("AVIF images in zip files is unsupported")

// Decorator adds image specific fields to a File.
type Decorator struct {
	FFProbe *ffmpeg.FFProbe
}

func (d *Decorator) Decorate(ctx context.Context, fs models.FS, f models.File) (models.File, error) {
	base := f.Base()

	// ignore clips in non-OsFS filesystems as ffprobe cannot read them
	// TODO - copy to temp file if not an OsFS
	if _, isOs := fs.(*file.OsFS); !isOs {
		ext := strings.ToLower(filepath.Ext(base.Path))

		// AVIF images inside zip files are not supported
		if ext == ".avif" {
			return nil, fmt.Errorf("%w: %s", ErrUnsupportedAVIFInZip, base.Path)
		}

		// Go cannot decode JXL from a stream, so extract to a temp file and probe by path
		if ext == ".jxl" {
			return d.decorateViaTempFile(fs, f)
		}

		logger.Debugf("assuming ImageFile for non-OsFS file %q", base.Path)
		return decorateFallback(fs, f)
	}

	probe, err := d.FFProbe.NewVideoFile(base.Path)
	if err != nil {
		logger.Warnf("File %q could not be read with ffprobe: %s, assuming ImageFile", base.Path, err)
		return decorateFallback(fs, f)
	}

	// Fallback to catch non-animated avif images that FFProbe detects as video files
	if probe.Bitrate == 0 && probe.VideoCodec == "av1" {
		return &models.ImageFile{
			BaseFile: base,
			Format:   "avif",
			Width:    probe.Width,
			Height:   probe.Height,
		}, nil
	}

	isClip := true
	// This list is derived from ffmpegImageThumbnail in pkg/image/thumbnail. If one gets updated, the other should be as well
	for _, item := range []string{"png", "mjpeg", "webp", "bmp", "jpegxl"} {
		if item == probe.VideoCodec {
			isClip = false
		}
	}
	if isClip {
		videoFileDecorator := video.Decorator{FFProbe: d.FFProbe}
		return videoFileDecorator.Decorate(ctx, fs, f)
	}

	ret := &models.ImageFile{
		BaseFile: base,
		Format:   probe.VideoCodec,
		Width:    probe.Width,
		Height:   probe.Height,
	}

	// FFprobe has a known bug where it returns 0x0 dimensions for some animated WebP files
	// Fall back to image.DecodeConfig in this case.
	// See: https://trac.ffmpeg.org/ticket/4907
	if ret.Width == 0 || ret.Height == 0 {
		logger.Warnf("FFprobe returned invalid dimensions (%dx%d) for %q, trying fallback decoder", ret.Width, ret.Height, base.Path)
		c, format, err := decodeConfig(fs, base.Path)
		if err != nil {
			logger.Warnf("Fallback decoder failed for %q: %s. Proceeding with original FFprobe result", base.Path, err)
		} else {
			ret.Width = c.Width
			ret.Height = c.Height
			// Update format if it differs (fallback decoder may be more accurate)
			if format != "" && format != ret.Format {
				logger.Debugf("Updating format from %q to %q for %q", ret.Format, format, base.Path)
				ret.Format = format
			}
		}
	}

	adjustForOrientation(fs, base.Path, ret)

	return ret, nil
}

func decodeConfig(fs models.FS, path string) (config image.Config, format string, err error) {
	r, err := fs.Open(path)
	if err != nil {
		err = fmt.Errorf("reading image file %q: %w", path, err)
		return
	}
	defer r.Close()

	config, format, err = image.DecodeConfig(r)
	if err != nil {
		err = fmt.Errorf("decoding image file %q: %w", path, err)
		return
	}
	return
}

func decorateFallback(fs models.FS, f models.File) (models.File, error) {
	base := f.Base()
	path := base.Path

	c, format, err := decodeConfig(fs, path)
	if err != nil {
		return f, err
	}

	ret := &models.ImageFile{
		BaseFile: base,
		Format:   format,
		Width:    c.Width,
		Height:   c.Height,
	}

	adjustForOrientation(fs, path, ret)

	return ret, nil
}

// decorateViaTempFile extracts a non-OsFS file (e.g. inside a zip) to a temp file so ffprobe can read it by path, for formats like JXL that ffprobe reads but Go cannot decode from a stream.
func (d *Decorator) decorateViaTempFile(fs models.FS, f models.File) (models.File, error) {
	base := f.Base()

	r, err := fs.Open(base.Path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	// preserve the extension so ffprobe can detect the format
	tmp, err := os.CreateTemp("", "stash-image-*"+filepath.Ext(base.Path))
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmp.Name())

	if _, err := io.Copy(tmp, r); err != nil {
		tmp.Close()
		return nil, err
	}
	tmp.Close()

	probe, err := d.FFProbe.NewVideoFile(tmp.Name())
	if err != nil {
		return nil, err
	}

	ret := &models.ImageFile{
		BaseFile: base,
		Format:   probe.VideoCodec,
		Width:    probe.Width,
		Height:   probe.Height,
	}

	adjustForOrientation(fs, base.Path, ret)

	return ret, nil
}

func (d *Decorator) IsMissingMetadata(ctx context.Context, fs models.FS, f models.File) bool {
	const (
		unsetString = "unset"
		unsetNumber = -1
	)

	imf, isImage := f.(*models.ImageFile)
	vf, isVideo := f.(*models.VideoFile)

	switch {
	case isImage:
		return imf.Format == unsetString || imf.Width == unsetNumber || imf.Height == unsetNumber
	case isVideo:
		videoFileDecorator := video.Decorator{FFProbe: d.FFProbe}
		return videoFileDecorator.IsMissingMetadata(ctx, fs, vf)
	default:
		return true
	}
}
