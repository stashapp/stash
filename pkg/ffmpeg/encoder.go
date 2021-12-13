package ffmpeg

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/stashapp/stash/pkg/desktop"
	"github.com/stashapp/stash/pkg/logger"
)

type Encoder string

var (
	runningEncoders      = make(map[string][]*os.Process)
	runningEncodersMutex = sync.RWMutex{}
)

func registerRunningEncoder(path string, process *os.Process) {
	runningEncodersMutex.Lock()
	processes := runningEncoders[path]

	runningEncoders[path] = append(processes, process)
	runningEncodersMutex.Unlock()
}

func deregisterRunningEncoder(path string, process *os.Process) {
	runningEncodersMutex.Lock()
	defer runningEncodersMutex.Unlock()
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
	runningEncodersMutex.RLock()
	processes := runningEncoders[path]
	runningEncodersMutex.RUnlock()

	for _, process := range processes {
		// assume it worked, don't check for error
		logger.Infof("Killing encoder process for file: %s", path)
		if err := process.Kill(); err != nil {
			logger.Warnf("failed to kill process %v: %v", process.Pid, err)
		}

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

// FFmpeg runner with progress output, used for transcodes
func (e *Encoder) runTranscode(probeResult VideoFile, args []string) (string, error) {
	cmd := exec.Command(string(*e), args...)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		logger.Error("FFMPEG stderr not available: " + err.Error())
	}

	stdout, err := cmd.StdoutPipe()
	if nil != err {
		logger.Error("FFMPEG stdout not available: " + err.Error())
	}

	desktop.HideExecShell(cmd)
	if err = cmd.Start(); err != nil {
		return "", err
	}

	buf := make([]byte, 80)
	lastProgress := 0.0
	var errBuilder strings.Builder
	for {
		n, err := stderr.Read(buf)
		if n > 0 {
			data := string(buf[0:n])
			time := GetTimeFromRegex(data)
			if time > 0 && probeResult.Duration > 0 {
				progress := time / probeResult.Duration

				if progress > lastProgress+0.01 {
					logger.Infof("Progress %.2f", progress)
					lastProgress = progress
				}
			}

			errBuilder.WriteString(data)
		}
		if err != nil {
			break
		}
	}

	stdoutData, _ := io.ReadAll(stdout)
	stdoutString := string(stdoutData)

	registerRunningEncoder(probeResult.Path, cmd.Process)
	err = waitAndDeregister(probeResult.Path, cmd)

	if err != nil {
		// error message should be in the stderr stream
		logger.Errorf("ffmpeg error when running command <%s>: %s", strings.Join(cmd.Args, " "), errBuilder.String())
		return stdoutString, err
	}

	return stdoutString, nil
}

func (e *Encoder) run(sourcePath string, args []string, stdin io.Reader) (string, error) {
	cmd := exec.Command(string(*e), args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = stdin

	desktop.HideExecShell(cmd)
	if err := cmd.Start(); err != nil {
		return "", err
	}

	var err error
	if sourcePath != "" {
		registerRunningEncoder(sourcePath, cmd.Process)
		err = waitAndDeregister(sourcePath, cmd)
	} else {
		err = cmd.Wait()
	}

	if err != nil {
		// error message should be in the stderr stream
		logger.Errorf("ffmpeg error when running command <%s>: %s", strings.Join(cmd.Args, " "), stderr.String())
		return stdout.String(), err
	}

	return stdout.String(), nil
}
