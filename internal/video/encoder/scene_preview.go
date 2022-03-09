package encoder

import (
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/ffmpeg/transcoder"
)

const (
	scenePreviewWidth        = 640
	scenePreviewAudioBitrate = "128k"

	scenePreviewImageFPS = 12
)

type ScenePreviewChunkOptions struct {
	StartTime  float64
	Duration   float64
	OutputPath string
	Audio      bool
	Preset     string
}

func ScenePreviewVideoChunk(encoder ffmpeg.FFMpeg, fn string, options ScenePreviewChunkOptions, fallback bool) error {
	var videoFilter ffmpeg.VideoFilter
	videoFilter = videoFilter.ScaleWidth(scenePreviewWidth)

	var videoArgs ffmpeg.Args
	videoArgs = videoArgs.VideoFilter(videoFilter)

	videoArgs = append(videoArgs,
		"-pix_fmt", "yuv420p",
		"-profile:v", "high",
		"-level", "4.2",
		"-preset", options.Preset,
		"-crf", "21",
		"-threads", "4",
		"-strict", "-2",
	)

	trimOptions := transcoder.TranscodeOptions{
		OutputPath: options.OutputPath,
		StartTime:  options.StartTime,
		Duration:   options.Duration,

		XError:   !fallback,
		SlowSeek: fallback,

		VideoCodec: ffmpeg.VideoCodecLibX264,
		VideoArgs:  videoArgs,
	}

	if options.Audio {
		var audioArgs ffmpeg.Args
		audioArgs = audioArgs.AudioBitrate(scenePreviewAudioBitrate)

		trimOptions.AudioCodec = ffmpeg.AudioCodecAAC
		trimOptions.AudioArgs = audioArgs
	}

	args := transcoder.Transcode(fn, trimOptions)

	return doGenerate(encoder, fn, args)
}

func ScenePreviewVideoToImage(encoder ffmpeg.FFMpeg, fn string, outputPath string) error {
	var videoFilter ffmpeg.VideoFilter
	videoFilter = videoFilter.ScaleWidth(scenePreviewWidth)
	videoFilter = videoFilter.Fps(scenePreviewImageFPS)

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

	encodeOptions := transcoder.TranscodeOptions{
		OutputPath: outputPath,

		VideoCodec: ffmpeg.VideoCodecLibWebP,
		VideoArgs:  videoArgs,
	}

	args := transcoder.Transcode(fn, encodeOptions)

	return doGenerate(encoder, fn, args)
}

func ScenePreviewVideoChunkCombine(encoder ffmpeg.FFMpeg, concatFilePath string, outputPath string) error {
	spliceOptions := transcoder.SpliceOptions{
		OutputPath: outputPath,
	}

	args := transcoder.Splice(concatFilePath, spliceOptions)

	return doGenerate(encoder, concatFilePath, args)
}
