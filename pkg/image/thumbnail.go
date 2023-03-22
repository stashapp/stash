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
	GenerateThumbnail(ctx context.Context, i *models.Image, f *file.ImageFile) error
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
func (e *ThumbnailEncoder) GetThumbnail(f *file.ImageFile, maxSize int) ([]byte, error) {
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

	format := f.Format
	animated := f.Format == formatGif

	// #2266 - if image is webp, then determine if it is animated
	if format == formatWebP {
		animated = isWebPAnimated(data)
	}

	// #2266 - don't generate a thumbnail for animated images
	if animated {
		return nil, fmt.Errorf("%w: %s", ErrNotSupportedForThumbnail, format)
	}

	if f.Clip {
		fileData, err := e.ffprobe.NewVideoFile(f.Path)
		if err != nil {
			return nil, err
		}
		if f.Width <= maxSize {
			maxSize = f.Width
		}
		clipDuration := fileData.VideoStreamDuration
		if clipDuration > 30.0 {
			clipDuration = 30.0
		}
		return e.getClipThumbnail(buf, maxSize, clipDuration, fileData.FrameRate)
	}

	// vips has issues loading files from stdin on Windows
	if e.vips != nil && runtime.GOOS != "windows" {
		return e.vips.ImageThumbnail(buf, maxSize)
	} else {
		return e.ffmpegImageThumbnail(buf, format, maxSize)
	}
}

func (e *ThumbnailEncoder) ffmpegImageThumbnail(image *bytes.Buffer, format string, maxSize int) ([]byte, error) {
	var ffmpegFormat ffmpeg.ImageFormat

	// These Options are used in isImage and isClip in internal/manager/manager_tasks. If one gets updated, the other should be as well
	switch format {
	case "mjpeg":
		ffmpegFormat = ffmpeg.ImageFormatJpeg
        case "jpeg":
                ffmpegFormat = ffmpeg.ImageFormatJpeg
	case "png":
		ffmpegFormat = ffmpeg.ImageFormatPng
	case "webp":
		ffmpegFormat = ffmpeg.ImageFormatWebp
	default:
		return nil, ErrUnsupportedImageFormat
	}

	args := transcoder.ImageThumbnail("-", transcoder.ImageThumbnailOptions{
		InputFormat:   ffmpegFormat,
		OutputPath:    "-",
		MaxDimensions: maxSize,
		Quality:       ffmpegImageQuality,
	})

	return e.ffmpeg.GenerateOutput(context.TODO(), args, image)
}

func (e *ThumbnailEncoder) getClipThumbnail(image *bytes.Buffer, maxSize int, clipDuration float64, frameRate float64) ([]byte, error) {
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
