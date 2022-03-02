package manager

import (
	"bytes"
	"fmt"
	"math"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/stashapp/stash/pkg/desktop"
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/utils"
)

type GeneratorInfo struct {
	ChunkCount     int
	FrameRate      float64
	NumberOfFrames int

	// NthFrame used for sprite generation
	NthFrame int

	ChunkDuration float64
	ExcludeStart  string
	ExcludeEnd    string

	VideoFile ffmpeg.VideoFile

	Audio bool // used for preview generation
}

func newGeneratorInfo(videoFile ffmpeg.VideoFile) (*GeneratorInfo, error) {
	exists, err := utils.FileExists(videoFile.Path)
	if !exists {
		logger.Errorf("video file not found")
		return nil, err
	}

	generator := &GeneratorInfo{VideoFile: videoFile}
	return generator, nil
}

func (g *GeneratorInfo) calculateFrameRate(videoStream *ffmpeg.FFProbeStream) error {
	var framerate float64
	if g.VideoFile.FrameRate == 0 {
		framerate, _ = strconv.ParseFloat(videoStream.RFrameRate, 64)
	} else {
		framerate = g.VideoFile.FrameRate
	}

	numberOfFrames, _ := strconv.Atoi(videoStream.NbFrames)

	if numberOfFrames == 0 && utils.IsValidFloat64(framerate) && g.VideoFile.Duration > 0 { // TODO: test
		numberOfFrames = int(framerate * g.VideoFile.Duration)
	}

	// If we are missing the frame count or frame rate then seek through the file and extract the info with regex
	if numberOfFrames == 0 || !utils.IsValidFloat64(framerate) {
		args := []string{
			"-nostats",
			"-i", g.VideoFile.Path,
			"-vcodec", "copy",
			"-f", "rawvideo",
			"-y",
		}
		if runtime.GOOS == "windows" {
			args = append(args, "nul") // https://stackoverflow.com/questions/313111/is-there-a-dev-null-on-windows
		} else {
			args = append(args, "/dev/null")
		}

		command := exec.Command(string(instance.FFMPEG), args...)
		desktop.HideExecShell(command)
		var stdErrBuffer bytes.Buffer
		command.Stderr = &stdErrBuffer // Frames go to stderr rather than stdout
		if err := command.Run(); err == nil {
			stdErrString := stdErrBuffer.String()
			if numberOfFrames == 0 {
				numberOfFrames = ffmpeg.GetFrameFromRegex(stdErrString)
			}
			if !utils.IsValidFloat64(framerate) {
				time := ffmpeg.GetTimeFromRegex(stdErrString)
				framerate = math.Round((float64(numberOfFrames)/time)*100) / 100
			}
		}
	}

	// Something seriously wrong with this file
	if numberOfFrames == 0 || !utils.IsValidFloat64(framerate) {
		logger.Errorf(
			"number of frames or framerate is 0.  nb_frames <%s> framerate <%f> duration <%f>",
			videoStream.NbFrames,
			framerate,
			g.VideoFile.Duration,
		)
	}

	g.FrameRate = framerate
	g.NumberOfFrames = numberOfFrames

	return nil
}

func (g *GeneratorInfo) configure() error {
	videoStream := g.VideoFile.VideoStream
	if videoStream == nil {
		return fmt.Errorf("missing video stream")
	}

	if err := g.calculateFrameRate(videoStream); err != nil {
		return err
	}

	// #2250 - ensure ChunkCount is valid
	if g.ChunkCount < 1 {
		logger.Warnf("[generator] Segment count (%d) must be > 0. Using 1 instead.", g.ChunkCount)
		g.ChunkCount = 1
	}

	g.NthFrame = g.NumberOfFrames / g.ChunkCount

	return nil
}

func (g GeneratorInfo) getExcludeValue(v string) float64 {
	if strings.HasSuffix(v, "%") && len(v) > 1 {
		// proportion of video duration
		v = v[0 : len(v)-1]
		prop, _ := strconv.ParseFloat(v, 64)
		return prop / 100.0 * g.VideoFile.Duration
	}

	prop, _ := strconv.ParseFloat(v, 64)
	return prop
}

// getStepSizeAndOffset calculates the step size for preview generation and
// the starting offset.
//
// Step size is calculated based on the duration of the video file, minus the
// excluded duration. The offset is based on the ExcludeStart. If the total
// excluded duration exceeds the duration of the video, then offset is 0, and
// the video duration is used to calculate the step size.
func (g GeneratorInfo) getStepSizeAndOffset() (stepSize float64, offset float64) {
	duration := g.VideoFile.Duration
	excludeStart := g.getExcludeValue(g.ExcludeStart)
	excludeEnd := g.getExcludeValue(g.ExcludeEnd)

	if duration > excludeStart+excludeEnd {
		duration = duration - excludeStart - excludeEnd
		offset = excludeStart
	}

	stepSize = duration / float64(g.ChunkCount)
	return
}
