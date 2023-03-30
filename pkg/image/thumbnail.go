package image

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"sync"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/ffmpeg/transcoder"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/models"
)

const ffmpegImageQuality = 5

var vipsPath string
var once sync.Once

var (
	ErrUnsupportedImageFormat = errors.New("unsupported image format")

	// ErrNotSupportedForThumbnail is returned if the image format is not supported for thumbnail generation
	ErrNotSupportedForThumbnail = errors.New("unsupported image format for thumbnail")
)

type ThumbnailGenerator interface {
	GenerateThumbnail(ctx context.Context, i *models.Image, f file.File) error
	GeneratePreview(ctx context.Context, i *models.Image, f file.File) error
}

type ThumbnailEncoder struct {
	ffmpeg     *ffmpeg.FFMpeg
	ffprobe    ffmpeg.FFProbe
	inputArgs  []string
	outputArgs []string
	preset     string
	vips       *vipsEncoder
}

func GetVipsPath() string {
	once.Do(func() {
		vipsPath, _ = exec.LookPath("vips")
	})
	return vipsPath
}

func NewThumbnailEncoder(ffmpegEncoder *ffmpeg.FFMpeg, ffProbe ffmpeg.FFProbe, inputArgs []string, outputArgs []string, preset string) ThumbnailEncoder {
	ret := ThumbnailEncoder{
		ffmpeg:     ffmpegEncoder,
		ffprobe:    ffProbe,
		inputArgs:  inputArgs,
		outputArgs: outputArgs,
		preset:     preset,
	}

	vipsPath := GetVipsPath()
	if vipsPath != "" {
		vipsEncoder := vipsEncoder(vipsPath)
		ret.vips = &vipsEncoder
	}

	return ret
}

// GetThumbnail returns the thumbnail image of the provided image resized to
// the provided max size. It resizes based on the largest X/Y direction.
// It returns nil and an error if an error occurs reading, decoding or encoding
// the image, or if the image is not suitable for thumbnails.
func (e *ThumbnailEncoder) GetThumbnail(f file.File, maxSize int) ([]byte, error) {
	reader, err := f.Open(&file.OsFS{})
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(reader); err != nil {
		return nil, err
	}

	data := buf.Bytes()

	if imageFile, ok := f.(*file.ImageFile); ok {
		format := imageFile.Format
		animated := imageFile.Format == formatGif

		// #2266 - if image is webp, then determine if it is animated
		if format == formatWebP {
			animated = isWebPAnimated(data)
		}

		// #2266 - don't generate a thumbnail for animated images
		if animated {
			return nil, fmt.Errorf("%w: %s", ErrNotSupportedForThumbnail, format)
		}
	}

	// Videofiles can only be thumbnailed with ffmpeg
	if _, ok := f.(*file.VideoFile); ok {
		return e.ffmpegImageThumbnail(buf, maxSize)
	}

	// vips has issues loading files from stdin on Windows
	if e.vips != nil && runtime.GOOS != "windows" {
		return e.vips.ImageThumbnail(buf, maxSize)
	} else {
		return e.ffmpegImageThumbnail(buf, maxSize)
	}
}

// GetPreview returns the preview clip of the provided image clip resized to
// the provided max size. It resizes based on the largest X/Y direction.
// It returns nil and an error if an error occurs reading, decoding or encoding
// the image, or if the image is not suitable for thumbnails.
// It is hardcoded to 30 seconds maximum right now
func (e *ThumbnailEncoder) GetPreview(f file.File, maxSize int) ([]byte, error) {
	reader, err := f.Open(&file.OsFS{})
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(reader); err != nil {
		return nil, err
	}

	fileData, err := e.ffprobe.NewVideoFile(f.Base().Path)
	if err != nil {
		return nil, err
	}
	if fileData.Width <= maxSize {
		maxSize = fileData.Width
	}
	clipDuration := fileData.VideoStreamDuration
	if clipDuration > 30.0 {
		clipDuration = 30.0
	}
	return e.getClipPreview(buf, maxSize, clipDuration, fileData.FrameRate)
}

func (e *ThumbnailEncoder) ffmpegImageThumbnail(image *bytes.Buffer, maxSize int) ([]byte, error) {
	args := transcoder.ImageThumbnail("-", transcoder.ImageThumbnailOptions{
		OutputFormat:  ffmpeg.ImageFormatJpeg,
		OutputPath:    "-",
		MaxDimensions: maxSize,
		Quality:       ffmpegImageQuality,
	})

	return e.ffmpeg.GenerateOutput(context.TODO(), args, image)
}

func (e *ThumbnailEncoder) getClipPreview(image *bytes.Buffer, maxSize int, clipDuration float64, frameRate float64) ([]byte, error) {
	var thumbFilter ffmpeg.VideoFilter
	thumbFilter = thumbFilter.ScaleMaxSize(maxSize)

	var thumbArgs ffmpeg.Args
	thumbArgs = thumbArgs.VideoFilter(thumbFilter)

	thumbArgs = append(thumbArgs,
		"-pix_fmt", "yuv420p",
		"-preset", e.preset,
		"-crf", "25",
		"-threads", "4",
		"-strict", "-2",
		"-f", "webm",
	)

	if frameRate <= 0.01 {
		thumbArgs = append(thumbArgs, "-vsync", "2")
	}

	thumbOptions := transcoder.TranscodeOptions{
		OutputPath: "-",
		StartTime:  0,
		Duration:   clipDuration,

		XError:   true,
		SlowSeek: false,

		VideoCodec: ffmpeg.VideoCodecVP9,
		VideoArgs:  thumbArgs,

		ExtraInputArgs:  e.inputArgs,
		ExtraOutputArgs: e.outputArgs,
	}

	args := transcoder.Transcode("-", thumbOptions)
	return e.ffmpeg.GenerateOutput(context.TODO(), args, image)
}
