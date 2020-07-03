package plugin

import (
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin/common"
)

func findArg(args []*models.PluginArgInput, name string) *models.PluginArgInput {
	for _, v := range args {
		if v.Key == name {
			return v
		}
	}

	return nil
}

func applyDefaultArgs(args []*models.PluginArgInput, defaultArgs map[string]string) []*models.PluginArgInput {
	for k, v := range defaultArgs {
		if arg := findArg(args, k); arg == nil {
			args = append(args, &models.PluginArgInput{
				Key: k,
				Value: &models.PluginValueInput{
					Str: &v,
				},
			})
		}
	}

	return args
}

func buildPluginInput(args []*models.PluginArgInput, serverConnection common.StashServerConnection) common.PluginInput {
	return common.PluginInput{
		ServerConnection: serverConnection,
		Args:             toPluginArgs(args),
	}
}
