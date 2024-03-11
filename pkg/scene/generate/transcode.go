package generate

import (
	"context"

	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/ffmpeg/transcoder"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
)

type TranscodeOptions struct {
	Width  int
	Height int
}

func (g Generator) Transcode(ctx context.Context, input string, hash string, options TranscodeOptions) error {
	lockCtx := g.LockManager.ReadLock(ctx, input)
	defer lockCtx.Cancel()

	return g.makeTranscode(lockCtx, hash, g.transcode(input, options))
}

// TranscodeVideo transcodes the video, and removes the audio.
// In some videos where the audio codec is not supported by ffmpeg,
// ffmpeg fails if you try to transcode the audio
func (g Generator) TranscodeVideo(ctx context.Context, input string, hash string, options TranscodeOptions) error {
	lockCtx := g.LockManager.ReadLock(ctx, input)
	defer lockCtx.Cancel()

	return g.makeTranscode(lockCtx, hash, g.transcodeVideo(input, options))
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

func (g Generator) transcode(input string, options TranscodeOptions) generateFn {
	return func(lockCtx *fsutil.LockContext, tmpFn string) error {
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
			OutputPath: tmpFn,
			VideoCodec: ffmpeg.VideoCodecLibX264,
			VideoArgs:  videoArgs,
			AudioCodec: ffmpeg.AudioCodecAAC,

			ExtraInputArgs:  g.FFMpegConfig.GetTranscodeInputArgs(),
			ExtraOutputArgs: g.FFMpegConfig.GetTranscodeOutputArgs(),
		})

		return g.generate(lockCtx, args)
	}
}

func (g Generator) transcodeVideo(input string, options TranscodeOptions) generateFn {
	return func(lockCtx *fsutil.LockContext, tmpFn string) error {
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
			OutputPath: tmpFn,
			VideoCodec: ffmpeg.VideoCodecLibX264,
			VideoArgs:  videoArgs,
			AudioArgs:  audioArgs,

			ExtraInputArgs:  g.FFMpegConfig.GetTranscodeInputArgs(),
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
