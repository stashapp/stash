package generate

import (
	"context"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/ffmpeg/transcoder"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
)

const transcodeCRF = 23

type TranscodeOptions struct {
	// MaxSize is the maximum resolution of the smaller dimension.
	// 0 means no scaling.
	MaxSize int
}

func (g Generator) Transcode(ctx context.Context, path string, width, height int, hash string, options TranscodeOptions) error {
	lockCtx := g.LockManager.ReadLock(ctx, path)
	defer lockCtx.Cancel()

	return g.makeTranscode(lockCtx, hash, g.transcode(path, width, height, options))
}

// TranscodeVideo transcodes the video, and removes the audio.
// In some videos where the audio codec is not supported by ffmpeg,
// ffmpeg fails if you try to transcode the audio
func (g Generator) TranscodeVideo(ctx context.Context, path string, width, height int, hash string, options TranscodeOptions) error {
	lockCtx := g.LockManager.ReadLock(ctx, path)
	defer lockCtx.Cancel()

	return g.makeTranscode(lockCtx, hash, g.transcodeVideo(path, width, height, options))
}

// TranscodeAudio will copy the video stream as is, and transcode audio.
func (g Generator) TranscodeAudio(ctx context.Context, input string, hash string) error {
	lockCtx := g.LockManager.ReadLock(ctx, input)
	defer lockCtx.Cancel()

	return g.makeTranscode(lockCtx, hash, g.transcodeAudio(input))
}

// TranscodeCopyVideo will copy the video stream as is, and drop the audio stream.
func (g Generator) TranscodeCopyVideo(ctx context.Context, input string, hash string) error {
	lockCtx := g.LockManager.ReadLock(ctx, input)
	defer lockCtx.Cancel()

	return g.makeTranscode(lockCtx, hash, g.transcodeCopyVideo(input))
}

func (g Generator) makeTranscode(lockCtx *fsutil.LockContext, hash string, generateFn generateFn) error {
	output := g.ScenePaths.GetTranscodePath(hash)
	if !g.Overwrite {
		if exists, _ := fsutil.FileExists(output); exists {
			return nil
		}
	}

	if err := g.generateFile(lockCtx, g.ScenePaths, mp4Pattern, output, generateFn); err != nil {
		return err
	}

	logger.Debug("created transcode: ", output)

	return nil
}

func (g Generator) transcode(path string, width, height int, options TranscodeOptions) generateFn {
	return func(lockCtx *fsutil.LockContext, tmpFn string) error {
		targetHeight := height
		var videoFilter ffmpeg.VideoFilter
		if options.MaxSize != 0 && options.MaxSize < min(width, height) {
			if width >= height {
				targetHeight = options.MaxSize
			} else {
				targetHeight = ffmpeg.ScaledHeight(width, height, options.MaxSize)
			}
		}

		codec, fullhw := g.DetermineCodecAndHW(lockCtx, path, width, height, targetHeight)

		if options.MaxSize != 0 {
			videoFilter = g.Encoder.HWMaxResFilter(codec, width, height, targetHeight, fullhw)
		}

		videoArgs := codec.ExtraArgsHQ("superfast", transcodeCRF)
		videoArgs = videoArgs.VideoFilter(videoFilter)

		extraInputArgs := g.Encoder.HWDeviceInit(g.FFMpegConfig.GetTranscodeInputArgs(), codec, fullhw)

		args := transcoder.Transcode(path, transcoder.TranscodeOptions{
			OutputPath: tmpFn,
			VideoCodec: codec,
			VideoArgs:  videoArgs,
			AudioCodec: ffmpeg.AudioCodecAAC,

			ExtraInputArgs:  extraInputArgs,
			ExtraOutputArgs: g.FFMpegConfig.GetTranscodeOutputArgs(),
		})

		return g.generate(lockCtx, args)
	}
}

func (g Generator) transcodeVideo(path string, width, height int, options TranscodeOptions) generateFn {
	return func(lockCtx *fsutil.LockContext, tmpFn string) error {
		targetHeight := height
		var videoFilter ffmpeg.VideoFilter
		if options.MaxSize != 0 && options.MaxSize < min(width, height) {
			if width >= height {
				targetHeight = options.MaxSize
			} else {
				targetHeight = ffmpeg.ScaledHeight(width, height, options.MaxSize)
			}
		}

		codec, fullhw := g.DetermineCodecAndHW(lockCtx, path, width, height, targetHeight)

		if options.MaxSize != 0 {
			videoFilter = g.Encoder.HWMaxResFilter(codec, width, height, targetHeight, fullhw)
		}

		videoArgs := codec.ExtraArgsHQ("superfast", transcodeCRF)
		videoArgs = videoArgs.VideoFilter(videoFilter)

		var audioArgs ffmpeg.Args
		audioArgs = audioArgs.SkipAudio()

		extraInputArgs := g.Encoder.HWDeviceInit(g.FFMpegConfig.GetTranscodeInputArgs(), codec, fullhw)

		args := transcoder.Transcode(path, transcoder.TranscodeOptions{
			OutputPath: tmpFn,
			VideoCodec: codec,
			VideoArgs:  videoArgs,
			AudioArgs:  audioArgs,

			ExtraInputArgs:  extraInputArgs,
			ExtraOutputArgs: g.FFMpegConfig.GetTranscodeOutputArgs(),
		})

		return g.generate(lockCtx, args)
	}
}

func (g Generator) transcodeAudio(input string) generateFn {
	return func(lockCtx *fsutil.LockContext, tmpFn string) error {
		args := transcoder.Transcode(input, transcoder.TranscodeOptions{
			OutputPath: tmpFn,
			VideoCodec: ffmpeg.VideoCodecCopy,
			AudioCodec: ffmpeg.AudioCodecAAC,
		})

		return g.generate(lockCtx, args)
	}
}

func (g Generator) transcodeCopyVideo(input string) generateFn {
	return func(lockCtx *fsutil.LockContext, tmpFn string) error {

		var audioArgs ffmpeg.Args
		audioArgs = audioArgs.SkipAudio()

		args := transcoder.Transcode(input, transcoder.TranscodeOptions{
			OutputPath: tmpFn,
			VideoCodec: ffmpeg.VideoCodecCopy,
			AudioArgs:  audioArgs,
		})

		return g.generate(lockCtx, args)
	}
}
