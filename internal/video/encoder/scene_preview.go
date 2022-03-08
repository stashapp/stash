package encoder

import (
	"github.com/stashapp/stash/pkg/ffmpeg2"
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

func ScenePreviewVideoChunk(encoder ffmpeg2.FFMpeg, fn string, options ScenePreviewChunkOptions, fallback bool) error {
	var videoFilter ffmpeg2.VideoFilter
	videoFilter = videoFilter.ScaleWidth(scenePreviewWidth)

	var videoArgs ffmpeg2.Args
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

	trimOptions := ffmpeg2.TranscodeOptions{
		OutputPath: options.OutputPath,
		StartTime:  options.StartTime,
		Duration:   options.Duration,

		XError:   !fallback,
		SlowSeek: fallback,

		VideoCodec: ffmpeg2.VideoCodecLibX264,
		VideoArgs:  videoArgs,
	}

	if options.Audio {
		var audioArgs ffmpeg2.Args
		audioArgs = audioArgs.AudioBitrate(scenePreviewAudioBitrate)

		trimOptions.AudioCodec = ffmpeg2.AudioCodecAAC
		trimOptions.AudioArgs = audioArgs
	}

	args := ffmpeg2.Transcode(fn, trimOptions)

	return doGenerate(encoder, fn, args)
}

func ScenePreviewVideoToImage(encoder ffmpeg2.FFMpeg, fn string, outputPath string) error {
	var videoFilter ffmpeg2.VideoFilter
	videoFilter = videoFilter.ScaleWidth(scenePreviewWidth)
	videoFilter = videoFilter.Fps(scenePreviewImageFPS)

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

	encodeOptions := ffmpeg2.TranscodeOptions{
		OutputPath: outputPath,

		VideoCodec: ffmpeg2.VideoCodecLibWebP,
		VideoArgs:  videoArgs,
	}

	args := ffmpeg2.Transcode(fn, encodeOptions)

	return doGenerate(encoder, fn, args)
}

func ScenePreviewVideoChunkCombine(encoder ffmpeg2.FFMpeg, concatFilePath string, outputPath string) error {
	spliceOptions := ffmpeg2.SpliceOptions{
		OutputPath: outputPath,
	}

	args := ffmpeg2.Splice(concatFilePath, spliceOptions)

	return doGenerate(encoder, concatFilePath, args)
}
