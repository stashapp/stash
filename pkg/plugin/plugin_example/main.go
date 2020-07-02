// +build plugin_example

package main

import (
	"context"

	"github.com/stashapp/stash/pkg/plugin/common"
	"github.com/stashapp/stash/pkg/plugin/util"
)

type api struct{}

func (api) Run(input common.PluginInput, output *common.PluginOutput) error {
	client := util.NewClient(input)

	var m struct {
		ReloadScrapers bool `graphql:"reloadScrapers"`
	}

	vars := map[string]interface{}{}
	err := client.Mutate(context.Background(), &m, vars)
	if err != nil {
		return err
	}

	outputStr := "ok"
	*output = common.PluginOutput{
		Output: &common.PluginArgValue{
			Str: &outputStr,
		},
	}

	return nil
}

func main() {
	err := common.ServePlugin(api{})
	if err != nil {
		panic(err)
	}
}
