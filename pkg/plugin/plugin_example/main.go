// +build plugin_example

package main

import (
	"context"

	"github.com/shurcooL/graphql"

	"github.com/stashapp/stash/pkg/plugin/common"
)

func main() {
	client := graphql.NewClient("http://localhost:9999/graphql", nil)

	var m struct {
		ReloadScrapers bool `graphql:"reloadScrapers"`
	}

	vars := map[string]interface{}{}
	err := client.Mutate(context.Background(), &m, vars)
	if err != nil {
		common.Error(err)
	}

	o := common.PluginOutput{
		Output: "ok",
	}
	o.Dispatch()
}
