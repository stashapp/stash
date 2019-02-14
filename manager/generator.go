package manager

import (
	"fmt"
	"github.com/stashapp/stash/ffmpeg"
	"github.com/stashapp/stash/logger"
	"github.com/stashapp/stash/utils"
	"os/exec"
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
		command := `ffmpeg -nostats -i ` + g.VideoFile.Path + ` -vcodec copy -f rawvideo -y /dev/null 2>&1 | \
                       grep frame | \
                       awk '{split($0,a,"fps")}END{print a[1]}' | \
                       sed 's/.*= *//'`
		commandResult, _ := exec.Command(command).Output()
		numberOfFrames, _ := strconv.Atoi(string(commandResult))
		if numberOfFrames == 0 { // TODO: test
			numberOfFrames = int(framerate * g.VideoFile.Duration)
		}
	}
	g.NumberOfFrames = numberOfFrames
	g.NthFrame = g.NumberOfFrames / g.ChunkCount

	return nil
}
