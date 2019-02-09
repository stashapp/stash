package ffmpeg

import (
	"fmt"
	"github.com/stashapp/stash/internal/logger"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strconv"
)

var progressRegex = regexp.MustCompile(`time=(\d+):(\d+):(\d+.\d+)`)

type encoder struct {
	Path string
}

func NewEncoder(ffmpegPath string) encoder {
	return encoder{
		Path: ffmpegPath,
	}
}

type ScreenshotOptions struct {
	OutputPath string
	Quality int
	Time float64
	Width int
	Verbosity string
}

type TranscodeOptions struct {
	OutputPath string
}

func (e *encoder) Screenshot(probeResult FFProbeResult, options ScreenshotOptions) {
	if options.Verbosity == "" {
		options.Verbosity = "quiet"
	}
	if options.Quality == 0 {
		options.Quality = 1
	}
	args := []string{
		"-v", options.Verbosity,
		"-ss", fmt.Sprintf("%v", options.Time),
		"-y",
		"-i", probeResult.Path, // TODO: Wrap in quotes?
		"-vframes", "1",
		"-q:v", fmt.Sprintf("%v", options.Quality),
		"-vf", fmt.Sprintf("scale=%v:-1", options.Width),
		"-f", "image2",
		options.OutputPath,
	}
	_, _ = e.run(probeResult, args)
}

func (e *encoder) Transcode(probeResult FFProbeResult, options TranscodeOptions) {
	args := []string{
		"-i", probeResult.Path,
		"-c:v", "libx264",
		"-profile:v", "high",
		"-level", "4.2",
		"-preset", "superfast",
		"-crf", "23",
		"-vf", "scale=iw:-2",
		"-c:a", "aac",
		options.OutputPath,
	}
	_, _ = e.run(probeResult, args)
}

func (e *encoder) run(probeResult FFProbeResult, args []string) (string, error) {
	cmd := exec.Command(e.Path, args...)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		logger.Error("FFMPEG stderr not available: " + err.Error())
	}

	stdout, err := cmd.StdoutPipe()
	if nil != err {
		logger.Error("FFMPEG stdout not available: " + err.Error())
	}

	if err = cmd.Start(); err != nil {
		return "", err
	}

	buf := make([]byte, 80)
	for {
		n, err := stderr.Read(buf)
		if n > 0 {
			data := string(buf[0:n])
			regexResult := progressRegex.FindStringSubmatch(data)
			if len(regexResult) == 4 && probeResult.Duration > 0 {
				h, _ := strconv.ParseFloat(regexResult[1], 64)
				m, _ := strconv.ParseFloat(regexResult[2], 64)
				s, _ := strconv.ParseFloat(regexResult[3], 64)
				hours := h * 3600
				mins := m * 60
				secs := s
				time := hours + mins + secs
				progress := time / probeResult.Duration
				logger.Infof("Progress %.2f", progress)
			}
		}
		if err != nil {
			break
		}
	}

	stdoutData, _ := ioutil.ReadAll(stdout)
	stdoutString := string(stdoutData)

	if err := cmd.Wait(); err != nil {
		return stdoutString, err
	}

	return stdoutString, nil
}