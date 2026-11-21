package transcoder

import (
	"fmt"

	"github.com/stashapp/stash/pkg/ffmpeg"
)

// SubtitleExtract builds ffmpeg arguments to extract the embedded subtitle
// stream at streamIndex from input and convert it to WebVTT at outputPath.
func SubtitleExtract(input string, streamIndex int, outputPath string) ffmpeg.Args {
	var args ffmpeg.Args
	args = args.LogLevel(ffmpeg.LogLevelError)
	args = args.Overwrite()
	args = args.Input(input)
	args = args.Map(fmt.Sprintf("0:%d", streamIndex))
	args = args.SkipVideo()
	args = args.SkipAudio()
	args = args.Format(ffmpeg.Format("webvtt"))
	args = args.Output(outputPath)
	return args
}
