package plugin

import (
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin/common"
)

func toOperationValue(value *common.PluginArgValue) *models.OperationValue {
	if value == nil {
		return nil
	}

	ret := &models.OperationValue{
		Str: value.Str,
		B:   value.B,
		F:   value.F,
		I:   value.I,
	}

	for _, v := range value.A {
		ret.A = append(ret.A, toOperationValue(v))
	}

	for _, v := range value.O {
		ret.O = append(ret.O, &models.OperationKeyValue{
			Key:   v.Key,
			Value: toOperationValue(v.Value),
		})
	}

	return ret
}

func toOperationResult(pluginOutput common.PluginOutput) *models.OperationResult {
	ret := &models.OperationResult{
		Error: pluginOutput.Error,
	}

	ret.Result = toOperationValue(pluginOutput.Output)

	return ret
}
