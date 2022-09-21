package plugin

import (
	"github.com/stashapp/stash/pkg/plugin/common"
)

func toPluginArgs(args []*PluginArgInput) common.ArgsMap {
	ret := make(common.ArgsMap)
	for _, a := range args {
		ret[a.Key] = toPluginArgValue(a.Value)
	}

	return ret
}

func toPluginArgValue(arg *PluginValueInput) common.PluginArgValue {
	if arg == nil {
		return nil
	}

	switch {
	case arg.Str != nil:
		return common.PluginArgValue(*arg.Str)
	case arg.I != nil:
		return common.PluginArgValue(*arg.I)
	case arg.B != nil:
		return common.PluginArgValue(*arg.B)
	case arg.F != nil:
		return common.PluginArgValue(*arg.F)
	case arg.O != nil:
		return common.PluginArgValue(toPluginArgs(arg.O))
	case arg.A != nil:
		var ret []common.PluginArgValue
		for _, v := range arg.A {
			ret = append(ret, toPluginArgValue(v))
		}
		return common.PluginArgValue(ret)
	}

	return nil
}
