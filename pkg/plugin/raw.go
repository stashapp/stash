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
		ret.Output = &output
	}

	return ret
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
