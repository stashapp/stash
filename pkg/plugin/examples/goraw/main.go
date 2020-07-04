// +build plugin_example

package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	exampleCommon "github.com/stashapp/stash/pkg/plugin/examples/common"

	"github.com/stashapp/stash/pkg/plugin/common"
	"github.com/stashapp/stash/pkg/plugin/common/log"
	"github.com/stashapp/stash/pkg/plugin/util"
)

type api struct {
	stopping bool
}

func (a *api) Run(input common.PluginInput, output *common.PluginOutput) error {
	modeArg := input.Args.String("mode")

	var err error
	if modeArg == "" || modeArg == "add" {
		client := util.NewClient(input.ServerConnection)
		err = exampleCommon.AddTag(client)
	} else if modeArg == "remove" {
		client := util.NewClient(input.ServerConnection)
		err = exampleCommon.RemoveTag(client)
	} else if modeArg == "long" {
		err = a.doLongTask()
	} else if modeArg == "indef" {
		err = a.doIndefiniteTask()
	}

	if err != nil {
		errStr := err.Error()
		*output = common.PluginOutput{
			Error: &errStr,
		}
		return nil
	}

	outputStr := "ok"
	*output = common.PluginOutput{
		Output: &outputStr,
	}

	return nil
}

func (a *api) doLongTask() error {
	const total = 100
	upTo := 0

	log.Info("Doing long task")
	for upTo < total {
		time.Sleep(time.Second)
		if a.stopping {
			return nil
		}

		log.Progress(float64(upTo) / float64(total))
		upTo++
	}

	return nil
}

func (a *api) doIndefiniteTask() error {
	log.Warn("Sleeping indefinitely")
	for {
		time.Sleep(time.Second)
		if a.stopping {
			return nil
		}
	}
}

func main() {
	input := common.PluginInput{}

	if len(os.Args) < 2 {
		log.Debug("Unmarshalling plugin input")
		inData, _ := ioutil.ReadAll(os.Stdin)
		log.Debugf("Raw input: %s", string(inData))
		decodeErr := json.Unmarshal(inData, &input)

		if decodeErr != nil {
			panic("missing mode argument")
		}
	} else {
		log.Debug("Using command line inputs")
		mode := os.Args[1]
		input.Args = common.ArgsMap{
			"mode": mode,
		}

		// just some hard-coded values
		input.ServerConnection = common.StashServerConnection{
			Scheme: "http",
			Port:   9999,
		}
	}

	a := api{}
	output := common.PluginOutput{}
	a.Run(input, &output)

	out, _ := json.Marshal(output)
	os.Stdout.WriteString(string(out))
}
