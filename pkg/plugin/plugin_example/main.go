// +build plugin_example

package main

import (
	"context"

	"github.com/stashapp/stash/pkg/plugin/common"
	"github.com/stashapp/stash/pkg/plugin/util"
)

func main() {
	input, err := common.ReadPluginInput()
	if err != nil {
		common.Error(err)
	}

	client := util.NewClient(input)

	var m struct {
		ReloadScrapers bool `graphql:"reloadScrapers"`
	}

	vars := map[string]interface{}{}
	err = client.Mutate(context.Background(), &m, vars)
	if err != nil {
		common.Error(err)
	}

	o := common.PluginOutput{
		Output: "ok",
	}
	o.Dispatch()
}
