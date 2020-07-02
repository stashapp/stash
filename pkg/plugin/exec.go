package plugin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin/common"
)

func writeInput(cmd *exec.Cmd, operation *PluginOperationConfig, args []*models.OperationArgInput) error {
	if operation.Interface == "" || operation.Interface == common.InterfaceJsonV1 {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return err
		}

		input := common.PluginInput{
			ServerPort: config.GetPort(),
		}

		go func() {
			defer stdin.Close()

			in, _ := json.Marshal(input)
			stdin.Write(in)
		}()
	}

	return nil
}

func readRawOutput(stdout io.ReadCloser) common.PluginOutput {
	output, err := ioutil.ReadAll(stdout)

	var errStr *string
	if err != nil {
		str := err.Error()
		errStr = &str
	}

	return common.PluginOutput{
		Output: string(output),
		Error:  errStr,
	}
}

func makeErrorOutput(err error) common.PluginOutput {
	str := err.Error()
	return common.PluginOutput{
		Error: &str,
	}
}

func readOutput(stdout io.ReadCloser, operation *PluginOperationConfig) common.PluginOutput {
	if operation.Interface == "" || operation.Interface == common.InterfaceJsonV1 {
		out := common.PluginOutput{}

		outStr, err := ioutil.ReadAll(stdout)
		if err != nil {
			return makeErrorOutput(err)
		}

		strReader := ioutil.NopCloser(bytes.NewReader(outStr))
		decodeErr := json.NewDecoder(strReader).Decode(&out)
		if decodeErr != nil {
			strReader = ioutil.NopCloser(bytes.NewReader(outStr))

			out = readRawOutput(strReader)
			if out.Error == nil {
				str := fmt.Sprintf("error decoding PluginOutput from stdout: %s", decodeErr.Error())
				out.Error = &str
			}
		}

		return out
	}

	return readRawOutput(stdout)
}

func executeOperation(operation *PluginOperationConfig, args []*models.OperationArgInput) common.PluginOutput {
	command := operation.Exec
	if len(command) == 0 {
		return makeErrorOutput(fmt.Errorf("empty exec value in operation %s", operation.Name))
	}

	// TODO - this should be the plugin config path, since it may be in a subdir
	_, err := exec.LookPath(command[0])
	if err != nil {
		// change command to use absolute path
		pluginPath := config.GetPluginsPath()
		command[0] = filepath.Join(pluginPath, command[0])
	}

	cmd := exec.Command(command[0], command[1:]...)

	// TODO - allow cwd to be changed in the operation config
	// TODO - this should be the plugin config path, since it may be in a subdir
	cmd.Dir = config.GetPluginsPath()

	err = writeInput(cmd, operation, args)
	if err != nil {
		return makeErrorOutput(err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		logger.Error("Plugin stderr not available: " + err.Error())
	}

	stdout, err := cmd.StdoutPipe()
	if nil != err {
		logger.Error("Plugin stdout not available: " + err.Error())
	}

	if err = cmd.Start(); err != nil {
		logger.Error("Error running plugin operation: " + err.Error())
		return makeErrorOutput(err)
	}

	// TODO - add a timeout here
	out := readOutput(stdout, operation)

	stderrData, _ := ioutil.ReadAll(stderr)
	stderrString := string(stderrData)

	err = cmd.Wait()

	if err != nil {
		// error message should be in the stderr stream
		logger.Errorf("error when running command <%s>: %s", strings.Join(cmd.Args, " "), stderrString)
		return makeErrorOutput(err)
	}

	return out
}
