package plugin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	stashExec "github.com/stashapp/stash/pkg/exec"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/plugin/common"
	"github.com/stashapp/stash/pkg/python"
)

type rawTaskBuilder struct{}

func (*rawTaskBuilder) build(task pluginTask) Task {
	return &rawPluginTask{
		pluginTask: task,
	}
}

type rawPluginTask struct {
	pluginTask

	started   bool
	waitGroup sync.WaitGroup
	cmd       *exec.Cmd
	done      chan bool
}

func (t *rawPluginTask) Start() error {
	if t.started {
		return errors.New("task already started")
	}

	command := t.plugin.getExecCommand(t.operation)
	if len(command) == 0 {
		return fmt.Errorf("empty exec value")
	}

	var cmd *exec.Cmd
	if python.IsPythonCommand(command[0]) {
		pythonPath := t.serverConfig.GetPythonPath()
		p, err := python.Resolve(pythonPath)

		if err != nil {
			logger.Warnf("%s", err)
		} else {
			cmd = p.Command(context.TODO(), command[1:])

			envVariable, _ := filepath.Abs(filepath.Dir(filepath.Dir(t.plugin.path)))
			python.AppendPythonPath(cmd, envVariable)
		}
	}

	if cmd == nil {
		// if could not find python, just use the command args as-is
		cmd = stashExec.Command(command[0], command[1:]...)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("error getting plugin process stdin: %v", err)
	}

	go func() {
		defer stdin.Close()

		inBytes, err := json.Marshal(t.input)
		if err != nil {
			logger.Warnf("error marshalling raw command input")
		}

		// Defensive: ensure no invalid JSON escape sequences reach the plugin
		// Go's json.Marshal should already produce valid JSON, but this adds a safety check
		inStr := string(inBytes)
		fixed := []rune{}
		for i, ch := range inStr {
			if ch == '\\' && i+1 < len(inStr) {
				nextCh := rune(inStr[i+1])
				// Valid JSON escape chars: " \ / b f n r t u
				validEscape := false
				for _, validCh := range `"\/bfnrtu` {
					if nextCh == validCh {
						validEscape = true
						break
					}
				}
				if validEscape {
					fixed = append(fixed, ch)
				} else {
					// Escape the backslash
					fixed = append(fixed, '\\', ch)
				}
			} else {
				fixed = append(fixed, ch)
			}
		}
		inBytes = []byte(string(fixed))

		if k, err := stdin.Write(inBytes); err != nil {
			logger.Warnf("error writing input to plugins stdin (wrote %v bytes out of %v): %v", k, len(string(inBytes)), err)
		}
	}()

	stderr, err := cmd.StderrPipe()
	if err != nil {
		logger.Error("plugin stderr not available: " + err.Error())
	}

	stdout, err := cmd.StdoutPipe()
	if nil != err {
		logger.Error("plugin stdout not available: " + err.Error())
	}

	t.waitGroup.Add(1)
	t.done = make(chan bool, 1)
	if err = cmd.Start(); err != nil {
		return fmt.Errorf("error running plugin: %v", err)
	}

	go t.handlePluginStderr(t.plugin.Name, stderr)
	t.cmd = cmd

	logger.Debugf("Plugin %s started: %s", t.plugin.Name, strings.Join(cmd.Args, " "))

	// send the stdout to the plugin output
	go func() {
		defer t.waitGroup.Done()
		defer close(t.done)
		stdoutData, _ := io.ReadAll(stdout)
		stdoutString := string(stdoutData)

		output := t.getOutput(stdoutString)

		err := cmd.Wait()
		if err != nil && output.Error == nil {
			errStr := err.Error()
			output.Error = &errStr
		}
		logger.Debugf("Plugin %s finished", t.plugin.Name)

		t.result = &output
	}()

	t.started = true
	return nil
}

func (t *rawPluginTask) getOutput(output string) common.PluginOutput {
	// try to parse the output as a PluginOutput json. If it fails just
	// get the raw output
	ret := common.PluginOutput{}
	decodeErr := json.Unmarshal([]byte(output), &ret)

	if decodeErr != nil {
		// Attempt to fix common invalid backslash escapes in JSON output
		fixed := fixInvalidJSONBackslashEscapes(output)
		tryRet := common.PluginOutput{}
		if err := json.Unmarshal([]byte(fixed), &tryRet); err == nil {
			return tryRet
		}

		ret.Output = &output
	}

	return ret
}

// fixInvalidJSONBackslashEscapes doubles backslashes that are not part of a
// valid JSON escape sequence so that the JSON decoder can parse outputs
// which contain unescaped backslashes (e.g., Windows paths).
func fixInvalidJSONBackslashEscapes(raw string) string {
	var b strings.Builder
	i := 0
	for i < len(raw) {
		if raw[i] == '\\' {
			if i+1 >= len(raw) {
				b.WriteString("\\\\")
				i++
				continue
			}
			next := raw[i+1]
			switch next {
			case '"', '\\', '/', 'b', 'f', 'n', 'r', 't', 'u':
				b.WriteByte('\\')
				b.WriteByte(next)
			default:
				b.WriteString("\\\\")
				b.WriteByte(next)
			}
			i += 2
		} else {
			b.WriteByte(raw[i])
			i++
		}
	}
	return b.String()
}

func (t *rawPluginTask) Wait() {
	t.waitGroup.Wait()
}

func (t *rawPluginTask) Stop() error {
	if t.cmd == nil {
		return nil
	}

	return t.cmd.Process.Kill()
}
