//go:build plugin_example
// +build plugin_example

package main

import (
	"context"
	"time"

	exampleCommon "github.com/stashapp/stash/pkg/plugin/examples/common"

	"github.com/stashapp/stash/pkg/plugin/common"
	"github.com/stashapp/stash/pkg/plugin/common/log"
	"github.com/stashapp/stash/pkg/plugin/util"
)

func main() {
	// serves the plugin, providing an object that satisfies the
	// common.RPCRunner interface
	err := common.ServePlugin(&api{})
	if err != nil {
		panic(err)
	}
}

type api struct {
	stopping bool
}

func (a *api) Stop(input struct{}, output *bool) error {
	log.Info("Stopping...")
	a.stopping = true
	*output = true
	return nil
}

// Run is the main work function of the plugin. It interprets the input and
// acts accordingly.
func (a *api) Run(input common.PluginInput, output *common.PluginOutput) error {
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

	return nil
}
