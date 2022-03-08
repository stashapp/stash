package encoder

import (
	"github.com/stashapp/stash/pkg/ffmpeg2"
)

const (
	markerPreviewWidth        = 640
	markerPreviewDuration     = 20
	markerPreviewAudioBitrate = "64k"

	markerImageDuration = 5

	markerScreenshotQuality = 2
)

type SceneMarkerOptions struct {
	Seconds    int
	OutputPath string
	Audio      bool
}

func SceneMarkerVideo(encoder ffmpeg2.FFMpeg, fn string, options SceneMarkerOptions) error {
	var videoFilter ffmpeg2.VideoFilter
	videoFilter = videoFilter.ScaleWidth(markerPreviewWidth)

	var videoArgs ffmpeg2.Args
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

	trimOptions := ffmpeg2.TranscodeOptions{
		Duration:   markerPreviewDuration,
		StartTime:  float64(options.Seconds),
		OutputPath: options.OutputPath,
		VideoCodec: ffmpeg2.VideoCodecLibX264,
		VideoArgs:  videoArgs,
	}

	if options.Audio {
		var audioArgs ffmpeg2.Args
		audioArgs = audioArgs.AudioBitrate(markerPreviewAudioBitrate)

		trimOptions.AudioCodec = ffmpeg2.AudioCodecAAC
		trimOptions.AudioArgs = audioArgs
	}

	args := ffmpeg2.Transcode(fn, trimOptions)

	return doGenerate(encoder, fn, args)
}

func SceneMarkerImage(encoder ffmpeg2.FFMpeg, fn string, options SceneMarkerOptions) error {
	var videoFilter ffmpeg2.VideoFilter
	videoFilter = videoFilter.ScaleWidth(markerPreviewWidth)
	videoFilter = videoFilter.Fps(12)

	var videoArgs ffmpeg2.Args
	videoArgs = videoArgs.VideoFilter(videoFilter)
	videoArgs = append(videoArgs,
		"-lossless", "1",
		"-q:v", "70",
		"-compression_level", "6",
		"-preset", "default",
		"-loop", "0",
		"-threads", "4",
	)

	trimOptions := ffmpeg2.TranscodeOptions{
		Duration:   markerImageDuration,
		StartTime:  float64(options.Seconds),
		OutputPath: options.OutputPath,
		VideoCodec: ffmpeg2.VideoCodecLibWebP,
		VideoArgs:  videoArgs,
	}

	args := ffmpeg2.Transcode(fn, trimOptions)

	return doGenerate(encoder, fn, args)
}

type SceneMarkerScreenshotOptions struct {
	Seconds    int
	OutputPath string
	Width      int
}

func SceneMarkerScreenshot(encoder ffmpeg2.FFMpeg, fn string, options SceneMarkerScreenshotOptions) error {
	ssOptions := ffmpeg2.ScreenshotOptions{
		OutputPath: options.OutputPath,
		OutputType: ffmpeg2.ScreenshotOutputTypeImage2,
		Quality:    markerScreenshotQuality,
		Width:      options.Width,
	}

	args := ffmpeg2.ScreenshotTime(fn, float64(options.Seconds), ssOptions)

	return doGenerate(encoder, fn, args)
}
