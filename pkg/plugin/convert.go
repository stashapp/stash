package plugin

import (
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin/common"
)

func toPluginArgs(args []*models.PluginArgInput) []*common.PluginKeyValue {
	var ret []*common.PluginKeyValue
	for _, a := range args {
		ret = append(ret, &common.PluginKeyValue{
			Key:   a.Key,
			Value: toPluginArgValue(a.Value),
		})
	}

	return ret
}

func toPluginArgValue(arg *models.PluginValueInput) *common.PluginArgValue {
	if arg == nil {
		return nil
	}

	ret := &common.PluginArgValue{
		Str: arg.Str,
		I:   arg.I,
		B:   arg.B,
		F:   arg.F,
		O:   toPluginArgs(arg.O),
	}

	for _, v := range arg.A {
		ret.A = append(ret.A, toPluginArgValue(v))
	}

	return ret
}
