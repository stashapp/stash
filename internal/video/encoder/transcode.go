package encoder

import (
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/ffmpeg/transcoder"
)

type TranscodeOptions struct {
	Width  int
	Height int
}

func Transcode(encoder ffmpeg.FFMpeg, input string, output string, options TranscodeOptions) error {
	var videoArgs ffmpeg.Args
	if options.Width != 0 && options.Height != 0 {
		var videoFilter ffmpeg.VideoFilter
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

	args := transcoder.Transcode(input, transcoder.TranscodeOptions{
		OutputPath: output,
		VideoCodec: ffmpeg.VideoCodecLibX264,
		VideoArgs:  videoArgs,
		AudioCodec: ffmpeg.AudioCodecAAC,
	})

	return doGenerate(encoder, input, args)
}

// TranscodeVideo transcodes the video, and removes the audio.
// In some videos where the audio codec is not supported by ffmpeg,
// ffmpeg fails if you try to transcode the audio
func TranscodeVideo(encoder ffmpeg.FFMpeg, input string, output string, options TranscodeOptions) error {
	var videoArgs ffmpeg.Args
	if options.Width != 0 && options.Height != 0 {
		var videoFilter ffmpeg.VideoFilter
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

	var audioArgs ffmpeg.Args
	audioArgs = audioArgs.SkipAudio()

	args := transcoder.Transcode(input, transcoder.TranscodeOptions{
		OutputPath: output,
		VideoCodec: ffmpeg.VideoCodecLibX264,
		VideoArgs:  videoArgs,
		AudioArgs:  audioArgs,
	})

	return doGenerate(encoder, input, args)
}

// TranscodeAudio will copy the video stream as is, and transcode audio.
func TranscodeAudio(encoder ffmpeg.FFMpeg, input string, output string, options TranscodeOptions) error {
	var videoArgs ffmpeg.Args
	if options.Width != 0 && options.Height != 0 {
		var videoFilter ffmpeg.VideoFilter
		videoFilter = videoFilter.ScaleDimensions(options.Width, options.Height)
		videoArgs = videoArgs.VideoFilter(videoFilter)
	}

	args := transcoder.Transcode(input, transcoder.TranscodeOptions{
		OutputPath: output,
		VideoCodec: ffmpeg.VideoCodecCopy,
		VideoArgs:  videoArgs,
		AudioCodec: ffmpeg.AudioCodecAAC,
	})

	return doGenerate(encoder, input, args)
}

// CopyVideo will copy the video stream as is, and drop the audio stream.
func CopyVideo(encoder ffmpeg.FFMpeg, input string, output string, options TranscodeOptions) error {
	var videoArgs ffmpeg.Args
	if options.Width != 0 && options.Height != 0 {
		var videoFilter ffmpeg.VideoFilter
		videoFilter = videoFilter.ScaleDimensions(options.Width, options.Height)
		videoArgs = videoArgs.VideoFilter(videoFilter)
	}

	var audioArgs ffmpeg.Args
	audioArgs = audioArgs.SkipAudio()

	args := transcoder.Transcode(input, transcoder.TranscodeOptions{
		OutputPath: output,
		VideoCodec: ffmpeg.VideoCodecCopy,
		VideoArgs:  videoArgs,
		AudioArgs:  audioArgs,
	})

	return doGenerate(encoder, input, args)
}
