package encoder

import (
	"github.com/stashapp/stash/pkg/ffmpeg2"
)

type TranscodeOptions struct {
	Width  int
	Height int
}

func Transcode(encoder ffmpeg2.FFMpeg, input string, output string, options TranscodeOptions) error {
	var videoArgs ffmpeg2.Args
	if options.Width != 0 && options.Height != 0 {
		var videoFilter ffmpeg2.VideoFilter
		videoFilter = videoFilter.ScaleDimensions(options.Width, options.Height)
		videoArgs = videoArgs.VideoFilter(videoFilter)
	}

	videoArgs = append(videoArgs,
		"-pix_fmt", "yuv420p",
		"-profile:v", "high",
		"-level", "4.2",
		"-preset", "superfast",
		"-crf", "23",
	)

	args := ffmpeg2.Transcode(input, ffmpeg2.TranscodeOptions{
		OutputPath: output,
		VideoCodec: ffmpeg2.VideoCodecLibX264,
		VideoArgs:  videoArgs,
		AudioCodec: ffmpeg2.AudioCodecAAC,
	})

	return doGenerate(encoder, input, args)
}

// TranscodeVideo transcodes the video, and removes the audio.
// In some videos where the audio codec is not supported by ffmpeg,
// ffmpeg fails if you try to transcode the audio
func TranscodeVideo(encoder ffmpeg2.FFMpeg, input string, output string, options TranscodeOptions) error {
	var videoArgs ffmpeg2.Args
	if options.Width != 0 && options.Height != 0 {
		var videoFilter ffmpeg2.VideoFilter
		videoFilter = videoFilter.ScaleDimensions(options.Width, options.Height)
		videoArgs = videoArgs.VideoFilter(videoFilter)
	}

	videoArgs = append(videoArgs,
		"-pix_fmt", "yuv420p",
		"-profile:v", "high",
		"-level", "4.2",
		"-preset", "superfast",
		"-crf", "23",
	)

	var audioArgs ffmpeg2.Args
	audioArgs = audioArgs.SkipAudio()

	args := ffmpeg2.Transcode(input, ffmpeg2.TranscodeOptions{
		OutputPath: output,
		VideoCodec: ffmpeg2.VideoCodecLibX264,
		VideoArgs:  videoArgs,
		AudioArgs:  audioArgs,
	})

	return doGenerate(encoder, input, args)
}

// TranscodeAudio will copy the video stream as is, and transcode audio.
func TranscodeAudio(encoder ffmpeg2.FFMpeg, input string, output string, options TranscodeOptions) error {
	var videoArgs ffmpeg2.Args
	if options.Width != 0 && options.Height != 0 {
		var videoFilter ffmpeg2.VideoFilter
		videoFilter = videoFilter.ScaleDimensions(options.Width, options.Height)
		videoArgs = videoArgs.VideoFilter(videoFilter)
	}

	args := ffmpeg2.Transcode(input, ffmpeg2.TranscodeOptions{
		OutputPath: output,
		VideoCodec: ffmpeg2.VideoCodecCopy,
		VideoArgs:  videoArgs,
		AudioCodec: ffmpeg2.AudioCodecAAC,
	})

	return doGenerate(encoder, input, args)
}

// CopyVideo will copy the video stream as is, and drop the audio stream.
func CopyVideo(encoder ffmpeg2.FFMpeg, input string, output string, options TranscodeOptions) error {
	var videoArgs ffmpeg2.Args
	if options.Width != 0 && options.Height != 0 {
		var videoFilter ffmpeg2.VideoFilter
		videoFilter = videoFilter.ScaleDimensions(options.Width, options.Height)
		videoArgs = videoArgs.VideoFilter(videoFilter)
	}

	var audioArgs ffmpeg2.Args
	audioArgs = audioArgs.SkipAudio()

	args := ffmpeg2.Transcode(input, ffmpeg2.TranscodeOptions{
		OutputPath: output,
		VideoCodec: ffmpeg2.VideoCodecCopy,
		VideoArgs:  videoArgs,
		AudioArgs:  audioArgs,
	})

	return doGenerate(encoder, input, args)
}
