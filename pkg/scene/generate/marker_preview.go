package generate

import (
	"context"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/ffmpeg/transcoder"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
)

const (
	markerPreviewWidth        = 640
	markerPreviewDuration     = 20
	markerPreviewAudioBitrate = "64k"

	markerImageDuration = 5
	markerWebpFPS       = 12

	markerScreenshotQuality = 2
)

func (g Generator) MarkerPreviewVideo(ctx context.Context, input string, hash string, seconds int, includeAudio bool) error {
	output := g.MarkerPaths.GetVideoPreviewPath(hash, seconds)
	if !g.Overwrite {
		if exists, _ := fsutil.FileExists(output); exists {
			return nil
		}
	}

	if err := g.generateFile(ctx, g.MarkerPaths, mp4Pattern, output, g.markerPreviewVideo(input, sceneMarkerOptions{
		Seconds: seconds,
		Audio:   includeAudio,
	})); err != nil {
		return err
	}

	logger.Debug("created marker video: ", output)

	return nil
}

type sceneMarkerOptions struct {
	Seconds int
	Audio   bool
}

func (g Generator) markerPreviewVideo(input string, options sceneMarkerOptions) generateFn {
	return func(ctx context.Context, tmpFn string) error {
		var videoFilter ffmpeg.VideoFilter
		videoFilter = videoFilter.ScaleWidth(markerPreviewWidth)

		var videoArgs ffmpeg.Args
		videoArgs = videoArgs.VideoFilter(videoFilter)

		videoArgs = append(videoArgs,
			"-pix_fmt", "yuv420p",
			"-profile:v", "high",
			"-level", "4.2",
			"-preset", "veryslow",
			"-crf", "24",
			"-movflags", "+faststart",
			"-threads", "4",
			"-sws_flags", "lanczos",
			"-strict", "-2",
		)

		trimOptions := transcoder.TranscodeOptions{
			Duration:   markerPreviewDuration,
			StartTime:  float64(options.Seconds),
			OutputPath: tmpFn,
			VideoCodec: ffmpeg.VideoCodecLibX264,
			VideoArgs:  videoArgs,
		}

		if options.Audio {
			var audioArgs ffmpeg.Args
			audioArgs = audioArgs.AudioBitrate(markerPreviewAudioBitrate)

			trimOptions.AudioCodec = ffmpeg.AudioCodecAAC
			trimOptions.AudioArgs = audioArgs
		}

		args := transcoder.Transcode(input, trimOptions)

		return ffmpeg.Generate(ctx, g.Encoder, args)
	}
}

func (g Generator) SceneMarkerWebp(ctx context.Context, input string, hash string, seconds int) error {
	output := g.MarkerPaths.GetWebpPreviewPath(hash, seconds)
	if !g.Overwrite {
		if exists, _ := fsutil.FileExists(output); exists {
			return nil
		}
	}

	if err := g.generateFile(ctx, g.MarkerPaths, webpPattern, output, g.sceneMarkerWebp(input, sceneMarkerOptions{
		Seconds: seconds,
	})); err != nil {
		return err
	}

	logger.Debug("created marker image: ", output)

	return nil
}

func (g Generator) sceneMarkerWebp(input string, options sceneMarkerOptions) generateFn {
	return func(ctx context.Context, tmpFn string) error {
		var videoFilter ffmpeg.VideoFilter
		videoFilter = videoFilter.ScaleWidth(markerPreviewWidth)
		videoFilter = videoFilter.Fps(markerWebpFPS)

		var videoArgs ffmpeg.Args
		videoArgs = videoArgs.VideoFilter(videoFilter)
		videoArgs = append(videoArgs,
			"-lossless", "1",
			"-q:v", "70",
			"-compression_level", "6",
			"-preset", "default",
			"-loop", "0",
			"-threads", "4",
		)

		trimOptions := transcoder.TranscodeOptions{
			Duration:   markerImageDuration,
			StartTime:  float64(options.Seconds),
			OutputPath: tmpFn,
			VideoCodec: ffmpeg.VideoCodecLibWebP,
			VideoArgs:  videoArgs,
		}

		args := transcoder.Transcode(input, trimOptions)

		return ffmpeg.Generate(ctx, g.Encoder, args)
	}
}

func (g Generator) SceneMarkerScreenshot(ctx context.Context, input string, hash string, seconds int, width int) error {
	output := g.MarkerPaths.GetScreenshotPath(hash, seconds)
	if !g.Overwrite {
		if exists, _ := fsutil.FileExists(output); exists {
			return nil
		}
	}

	if err := g.generateFile(ctx, g.MarkerPaths, jpgPattern, output, g.sceneMarkerScreenshot(input, SceneMarkerScreenshotOptions{
		Seconds: seconds,
		Width:   width,
	})); err != nil {
		return err
	}

	logger.Debug("created marker screenshot: ", output)

	return nil
}

type SceneMarkerScreenshotOptions struct {
	Seconds int
	Width   int
}

func (g Generator) sceneMarkerScreenshot(input string, options SceneMarkerScreenshotOptions) generateFn {
	return func(ctx context.Context, tmpFn string) error {
		ssOptions := transcoder.ScreenshotOptions{
			OutputPath: tmpFn,
			OutputType: transcoder.ScreenshotOutputTypeImage2,
			Quality:    markerScreenshotQuality,
			Width:      options.Width,
		}

		args := transcoder.ScreenshotTime(input, float64(options.Seconds), ssOptions)

		return ffmpeg.Generate(ctx, g.Encoder, args)
	}
}
