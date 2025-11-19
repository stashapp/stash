package image

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/ffmpeg/transcoder"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/models"
	"golang.org/x/image/draw"
)

const ffmpegImageQuality = 5

var vipsPath string
var once sync.Once

var (
	ErrUnsupportedImageFormat = errors.New("unsupported image format")

	// ErrNotSupportedForThumbnail is returned if the image format is not supported for thumbnail generation
	ErrNotSupportedForThumbnail = errors.New("unsupported image format for thumbnail")
)

type ThumbnailEncoder struct {
	FFMpeg             *ffmpeg.FFMpeg
	FFProbe            *ffmpeg.FFProbe
	ClipPreviewOptions ClipPreviewOptions
	vips               *vipsEncoder
}

type ClipPreviewOptions struct {
	InputArgs  []string
	OutputArgs []string
	Preset     string
}

func GetVipsPath() string {
	once.Do(func() {
		vipsPath, _ = exec.LookPath("vips")
	})
	return vipsPath
}

func NewThumbnailEncoder(ffmpegEncoder *ffmpeg.FFMpeg, ffProbe *ffmpeg.FFProbe, clipPreviewOptions ClipPreviewOptions) ThumbnailEncoder {
	ret := ThumbnailEncoder{
		FFMpeg:             ffmpegEncoder,
		FFProbe:            ffProbe,
		ClipPreviewOptions: clipPreviewOptions,
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
func (e *ThumbnailEncoder) GetThumbnail(f models.File, maxSize int) ([]byte, error) {
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

	format := ""
	if imageFile, ok := f.(*models.ImageFile); ok {
		format = imageFile.Format
		animated := imageFile.Format == formatGif

		// #2266 - if image is webp, then determine if it is animated
		if format == formatWebP {
			animated = isWebPAnimated(data)
		}

		// #2266 - don't generate a thumbnail for animated images
		if animated {
			return nil, fmt.Errorf("%w: %s", ErrNotSupportedForThumbnail, format)
		}

		// AVIF cannot be read from stdin, must use file path
		// However, files in zip archives must be read from the buffer since the path is virtual
		if format == "avif" && f.Base().ZipFileID == nil {
			return e.ffmpegImageThumbnailFromPath(f.Base().Path, maxSize)
		}
	}

	// Videofiles can only be thumbnailed with ffmpeg
	if _, ok := f.(*models.VideoFile); ok {
		return e.ffmpegImageThumbnail(buf, maxSize, format)
	}

	// AVIF cannot be piped to FFmpeg, use vips or alternative for zip files
	// For AVIF in zip, we must avoid FFmpeg entirely
	if format == "avif" && f.Base().ZipFileID != nil {
		if e.vips != nil {
			return e.vips.ImageThumbnail(buf, maxSize)
		}
		// If no vips, use Go image processing as fallback
		return e.goImageThumbnail(buf, maxSize)
	}

	// vips has issues loading files from stdin on Windows
	if e.vips != nil && runtime.GOOS != "windows" {
		return e.vips.ImageThumbnail(buf, maxSize)
	} else {
		return e.ffmpegImageThumbnail(buf, maxSize, format)
	}
}

// GetPreview returns the preview clip of the provided image clip resized to
// the provided max size. It resizes based on the largest X/Y direction.
// It is hardcoded to 30 seconds maximum right now
func (e *ThumbnailEncoder) GetPreview(inPath string, outPath string, maxSize int) error {
	fileData, err := e.FFProbe.NewVideoFile(inPath)
	if err != nil {
		return err
	}
	if fileData.Width <= maxSize {
		maxSize = fileData.Width
	}
	clipDuration := fileData.VideoStreamDuration
	if clipDuration > 30.0 {
		clipDuration = 30.0
	}
	return e.getClipPreview(inPath, outPath, maxSize, clipDuration, fileData.FrameRate)
}

func (e *ThumbnailEncoder) ffmpegImageThumbnail(image *bytes.Buffer, maxSize int, format string) ([]byte, error) {
	options := transcoder.ImageThumbnailOptions{
		OutputFormat:  ffmpeg.ImageFormatJpeg,
		OutputPath:    "-",
		MaxDimensions: maxSize,
		Quality:       ffmpegImageQuality,
	}

	args := transcoder.ImageThumbnail("-", options)

	return e.FFMpeg.GenerateOutput(context.TODO(), args, image)
}

// ffmpegImageThumbnailFromPath generates a thumbnail from a file path (used for AVIF which can't be piped)
func (e *ThumbnailEncoder) ffmpegImageThumbnailFromPath(inputPath string, maxSize int) ([]byte, error) {
	options := transcoder.ImageThumbnailOptions{
		OutputFormat:  ffmpeg.ImageFormatJpeg,
		OutputPath:    "-",
		MaxDimensions: maxSize,
		Quality:       ffmpegImageQuality,
	}

	args := transcoder.ImageThumbnail(inputPath, options)

	return e.FFMpeg.GenerateOutput(context.TODO(), args, nil)
}

func (e *ThumbnailEncoder) getClipPreview(inPath string, outPath string, maxSize int, clipDuration float64, frameRate float64) error {
	var thumbFilter ffmpeg.VideoFilter
	thumbFilter = thumbFilter.ScaleMaxSize(maxSize)

	var thumbArgs ffmpeg.Args
	thumbArgs = thumbArgs.VideoFilter(thumbFilter)

	o := e.ClipPreviewOptions

	thumbArgs = append(thumbArgs,
		"-pix_fmt", "yuv420p",
		"-preset", o.Preset,
		"-crf", "25",
		"-threads", "4",
		"-strict", "-2",
		"-f", "webm",
	)

	if frameRate <= 0.01 {
		thumbArgs = append(thumbArgs, "-vsync", "2")
	}

	thumbOptions := transcoder.TranscodeOptions{
		OutputPath: outPath,
		StartTime:  0,
		Duration:   clipDuration,

		XError:   true,
		SlowSeek: false,

		VideoCodec: ffmpeg.VideoCodecVP9,
		VideoArgs:  thumbArgs,

		ExtraInputArgs:  o.InputArgs,
		ExtraOutputArgs: o.OutputArgs,
	}

	if err := fsutil.EnsureDirAll(filepath.Dir(outPath)); err != nil {
		return err
	}
	args := transcoder.Transcode(inPath, thumbOptions)
	return e.FFMpeg.Generate(context.TODO(), args)
}

// goImageThumbnail generates a thumbnail using Go's native image processing
// This is used as a fallback for AVIF images in zip files when vips is not available
func (e *ThumbnailEncoder) goImageThumbnail(buf *bytes.Buffer, maxSize int) ([]byte, error) {
	// Decode the image
	img, _, err := image.Decode(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return nil, fmt.Errorf("decoding image: %w", err)
	}

	// Calculate thumbnail dimensions
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Determine scaling factor
	var scale float64
	if width > height {
		scale = float64(maxSize) / float64(width)
	} else {
		scale = float64(maxSize) / float64(height)
	}

	// Don't upscale
	if scale > 1.0 {
		scale = 1.0
	}

	newWidth := int(float64(width) * scale)
	newHeight := int(float64(height) * scale)

	// Create thumbnail image
	thumbnail := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Use high-quality scaling
	draw.CatmullRom.Scale(thumbnail, thumbnail.Bounds(), img, bounds, draw.Over, nil)

	// Encode as JPEG
	out := new(bytes.Buffer)
	opts := &jpeg.Options{Quality: 95}
	if err := jpeg.Encode(out, thumbnail, opts); err != nil {
		return nil, fmt.Errorf("encoding thumbnail: %w", err)
	}

	return out.Bytes(), nil
}
