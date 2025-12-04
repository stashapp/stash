package api

import (
	"context"

	"github.com/stashapp/stash/internal/manager/config"
)

func (r *configDisableDropdownCreateResolver) Gallery(ctx context.Context, obj *config.ConfigDisableDropdownCreate) (bool, error) {
	return obj.Gallery, nil
}

func (r *configResultResolver) Plugins(ctx context.Context, obj *ConfigResult, include []string) (map[string]map[string]interface{}, error) {
	if len(include) == 0 {
		ret := config.GetInstance().GetAllPluginConfiguration()
		return ret, nil
	}

	ret := make(map[string]map[string]interface{})

	for _, plugin := range include {
		c := config.GetInstance().GetPluginConfiguration(plugin)
		if len(c) > 0 {
			ret[plugin] = c
		}
	}

	return ret, nil
}
