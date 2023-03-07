package ffmpeg

import (
	"bytes"
	"context"
	"math"
	"regexp"
	"strconv"
)

// FrameInfo contains the number of frames and the frame rate for a video file.
type FrameInfo struct {
	FrameRate      float64
	NumberOfFrames int
}

// CalculateFrameRate calculates the frame rate and number of frames of the video file.
// Used where the frame rate or NbFrames is missing or invalid in the ffprobe output.
func (f *FFMpeg) CalculateFrameRate(ctx context.Context, v *VideoFile) (*FrameInfo, error) {
	var args Args
	args = append(args, "-nostats")
	args = args.Input(v.Path).
		VideoCodec(VideoCodecCopy).
		Format(FormatRawVideo).
		Overwrite().
		NullOutput()

	command := f.Command(ctx, args)
	var stdErrBuffer bytes.Buffer
	command.Stderr = &stdErrBuffer // Frames go to stderr rather than stdout
	err := command.Run()
	if err == nil {
		var ret FrameInfo
		stdErrString := stdErrBuffer.String()
		ret.NumberOfFrames = getFrameFromRegex(stdErrString)

		time := getTimeFromRegex(stdErrString)
		ret.FrameRate = math.Round((float64(ret.NumberOfFrames)/time)*100) / 100

		return &ret, nil
	}

	return nil, err
}

var timeRegex = regexp.MustCompile(`time=\s*(\d+):(\d+):(\d+.\d+)`)
var frameRegex = regexp.MustCompile(`frame=\s*([0-9]+)`)

func getTimeFromRegex(str string) float64 {
	regexResult := timeRegex.FindStringSubmatch(str)

	// Bail early if we don't have the results we expect
	if len(regexResult) != 4 {
		return 0
	}

	h, _ := strconv.ParseFloat(regexResult[1], 64)
	m, _ := strconv.ParseFloat(regexResult[2], 64)
	s, _ := strconv.ParseFloat(regexResult[3], 64)
	hours := h * 3600
	minutes := m * 60
	seconds := s
	return hours + minutes + seconds
}

func getFrameFromRegex(str string) int {
	regexResult := frameRegex.FindStringSubmatch(str)

	// Bail early if we don't have the results we expect
	if len(regexResult) < 2 {
		return 0
	}

	result, _ := strconv.Atoi(regexResult[1])
	return result
}
