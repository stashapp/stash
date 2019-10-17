package ffmpeg

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/logger"
)

type Encoder struct {
	Path string
}

var runningEncoders map[string][]*os.Process = make(map[string][]*os.Process)

func NewEncoder(ffmpegPath string) Encoder {
	return Encoder{
		Path: ffmpegPath,
	}
}

func registerRunningEncoder(path string, process *os.Process) {
	processes := runningEncoders[path]

	runningEncoders[path] = append(processes, process)
}

func deregisterRunningEncoder(path string, process *os.Process) {
	processes := runningEncoders[path]

	for i, v := range processes {
		if v == process {
			runningEncoders[path] = append(processes[:i], processes[i+1:]...)
			return
		}
	}
}

func waitAndDeregister(path string, cmd *exec.Cmd) error {
	err := cmd.Wait()
	deregisterRunningEncoder(path, cmd.Process)

	return err
}

func KillRunningEncoders(path string) {
	processes := runningEncoders[path]

	for _, process := range processes {
		// assume it worked, don't check for error
		fmt.Printf("Killing encoder process for file: %s", path)
		process.Kill()

		// wait for the process to die before returning
		// don't wait more than a few seconds
		done := make(chan error)
		go func() { 
			_, err := process.Wait() 
			done <- err
		}()

		select {
		case <-done:
			return
		case <-time.After(5 * time.Second):
			return
		}
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
			time := GetTimeFromRegex(data)
			if time > 0 && probeResult.Duration > 0 {
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

	registerRunningEncoder(probeResult.Path, cmd.Process)
	err = waitAndDeregister(probeResult.Path, cmd)

	if err != nil {
		logger.Errorf("ffmpeg error when running command <%s>: %s", strings.Join(cmd.Args, " "), stdoutString)
		return stdoutString, err
	}

	return stdoutString, nil
}

func (e *Encoder) stream(probeResult VideoFile, args []string) (io.ReadCloser, *os.Process, error) {
	cmd := exec.Command(e.Path, args...)

	stdout, err := cmd.StdoutPipe()
	if nil != err {
		logger.Error("FFMPEG stdout not available: " + err.Error())
	}

	if err = cmd.Start(); err != nil {
		return nil, nil, err
	}

	registerRunningEncoder(probeResult.Path, cmd.Process)
	go waitAndDeregister(probeResult.Path, cmd)

	return stdout, cmd.Process, nil
}
