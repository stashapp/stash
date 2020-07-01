package plugin

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os/exec"
	"strings"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/plugin/common"
)

func RunOperationPlugin(command []string, input common.PluginInput, out *common.PluginOutput) error {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = config.GetScrapersPath()

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	go func() {
		defer stdin.Close()

		in, _ := json.Marshal(input)
		stdin.Write(in)
	}()

	stderr, err := cmd.StderrPipe()
	if err != nil {
		logger.Error("Scraper stderr not available: " + err.Error())
	}

	stdout, err := cmd.StdoutPipe()
	if nil != err {
		logger.Error("Scraper stdout not available: " + err.Error())
	}

	if err = cmd.Start(); err != nil {
		logger.Error("Error running scraper script: " + err.Error())
		return errors.New("Error running scraper script")
	}

	// TODO - add a timeout here
	decodeErr := json.NewDecoder(stdout).Decode(out)

	stderrData, _ := ioutil.ReadAll(stderr)
	stderrString := string(stderrData)

	err = cmd.Wait()

	if err != nil {
		// error message should be in the stderr stream
		logger.Errorf("error when running command <%s>: %s", strings.Join(cmd.Args, " "), stderrString)
		return errors.New("Error running script")
	}

	if decodeErr != nil {
		logger.Errorf("error decoding performer from data: %s", decodeErr.Error())
		return errors.New("Error decoding performer from script")
	}

	return nil
}
