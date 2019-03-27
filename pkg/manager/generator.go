package manager

import (
	"bytes"
	"fmt"
	"github.com/stashapp/stash/pkg/ffmpeg"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/utils"
	"os/exec"
	"regexp"
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
	g.FrameRate = framerate

	numberOfFrames, _ := strconv.Atoi(videoStream.NbFrames)
	if numberOfFrames == 0 {
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
			re := regexp.MustCompile(`frame[=] ([0-9]+)`)
			frames := re.FindStringSubmatch(stdErrBuffer.String())
			if frames != nil && len(frames) > 1 {
				numberOfFrames, _ = strconv.Atoi(frames[1])
			}
		}
	}
	if numberOfFrames == 0 { // TODO: test
		numberOfFrames = int(framerate * g.VideoFile.Duration)
	}
	if numberOfFrames == 0 {
		logger.Errorf(
			"number of frames is 0.  nb_frames <%s> framerate <%s> duration <%s>",
			videoStream.NbFrames,
			framerate,
			g.VideoFile.Duration,
			)
	}

	g.NumberOfFrames = numberOfFrames
	g.NthFrame = g.NumberOfFrames / g.ChunkCount

	return nil
}
