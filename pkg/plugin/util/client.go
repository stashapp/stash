package util

import (
	"strconv"

	"github.com/shurcooL/graphql"

	"github.com/stashapp/stash/pkg/plugin/common"
)

func NewClient(provider common.StashServerProvider) *graphql.Client {
	// TODO - handle https
	portStr := strconv.Itoa(provider.GetPort())
	return graphql.NewClient("http://localhost:"+portStr+"/graphql", nil)
}
