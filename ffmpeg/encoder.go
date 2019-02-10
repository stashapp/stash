package ffmpeg

import (
	"github.com/stashapp/stash/logger"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strconv"
)

var progressRegex = regexp.MustCompile(`time=(\d+):(\d+):(\d+.\d+)`)

type Encoder struct {
	Path string
}

func NewEncoder(ffmpegPath string) Encoder {
	return Encoder{
		Path: ffmpegPath,
	}
}

func (e *Encoder) run(probeResult VideoFile, args []string) (string, error) {
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