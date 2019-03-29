package manager

import (
	"bytes"
	"fmt"
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/utils"
	"math"
	"os/exec"
	"runtime"
	"strconv"
)

type GeneratorInfo struct {
	ChunkCount     int
	FrameRate      float64
	NumberOfFrames int
	NthFrame       int

	VideoFile ffmpeg.VideoFile
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

func (g *GeneratorInfo) configure() error {
	videoStream := g.VideoFile.VideoStream
	if videoStream == nil {
		return fmt.Errorf("missing video stream")
	}

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

		command := exec.Command(instance.FFMPEGPath, args...)
		var stdErrBuffer bytes.Buffer
		command.Stderr = &stdErrBuffer // Frames go to stderr rather than stdout
		if err := command.Run(); err == nil {
			stdErrString := stdErrBuffer.String()
			if numberOfFrames == 0 {
				numberOfFrames = ffmpeg.GetFrameFromRegex(stdErrString)
			}
			if utils.IsValidFloat64(framerate) {
				time := ffmpeg.GetTimeFromRegex(stdErrString)
				framerate = math.Round((float64(numberOfFrames)/time)*100) / 100
			}
		}
	}

	// Something seriously wrong with this file
	if numberOfFrames == 0 || !utils.IsValidFloat64(framerate) {
		logger.Errorf(
			"number of frames or framerate is 0.  nb_frames <%s> framerate <%s> duration <%s>",
			videoStream.NbFrames,
			framerate,
			g.VideoFile.Duration,
		)
	}

	g.FrameRate = framerate
	g.NumberOfFrames = numberOfFrames
	g.NthFrame = g.NumberOfFrames / g.ChunkCount

	return nil
}
