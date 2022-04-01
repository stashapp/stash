//go:build plugin_example
// +build plugin_example

package main

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"time"

	exampleCommon "github.com/stashapp/stash/pkg/plugin/examples/common"

	"github.com/stashapp/stash/pkg/plugin/common"
	"github.com/stashapp/stash/pkg/plugin/common/log"
	"github.com/stashapp/stash/pkg/plugin/util"
)

// raw plugins may accept the plugin input from stdin, or they can elect
// to ignore it entirely. In this case it optionally reads from the
// command-line parameters.
func main() {
	input := common.PluginInput{}

	if len(os.Args) < 2 {
		inData, _ := io.ReadAll(os.Stdin)
		log.Debugf("Raw input: %s", string(inData))
		decodeErr := json.Unmarshal(inData, &input)

		if decodeErr != nil {
			panic("missing mode argument")
		}
	} else {
		log.Debug("Using command line inputs")
		mode := os.Args[1]
		log.Debugf("Command line inputs: %v", os.Args[1:])
		input.Args = common.ArgsMap{
			"mode": mode,
		}

		// just some hard-coded values
		input.ServerConnection = common.StashServerConnection{
			Scheme: "http",
			Port:   9999,
		}
	}

	output := common.PluginOutput{}
	Run(input, &output)

	out, _ := json.Marshal(output)
	os.Stdout.WriteString(string(out))
}

func Run(input common.PluginInput, output *common.PluginOutput) error {
	modeArg := input.Args.String("mode")
	ctx := context.TODO()
	var err error
	if modeArg == "" || modeArg == "add" {
		client := util.NewClient(input.ServerConnection)
		err = exampleCommon.AddTag(ctx, client)
	} else if modeArg == "remove" {
		client := util.NewClient(input.ServerConnection)
		err = exampleCommon.RemoveTag(ctx, client)
	} else if modeArg == "long" {
		err = doLongTask()
	} else if modeArg == "indef" {
		err = doIndefiniteTask()
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

func doLongTask() error {
	const total = 100
	upTo := 0

	log.Info("Doing long task")
	for upTo < total {
		time.Sleep(time.Second)

		log.Progress(float64(upTo) / float64(total))
		upTo++
	}

	return nil
}

func doIndefiniteTask() error {
	log.Warn("Sleeping indefinitely")
	for {
		time.Sleep(time.Second)
	}
}
