package javascript

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"

	"github.com/dop251/goja"
	"github.com/stashapp/stash/pkg/logger"
)

const pluginPrefix = "[Plugin] "

type Log struct {
	Progress chan float64
}

func (l *Log) argToString(call goja.FunctionCall) string {
	arg := call.Argument(0)
	var o map[string]interface{}
	if arg.ExportType() == reflect.TypeOf(o) {
		ii := arg.Export()
		o = ii.(map[string]interface{})
		data, err := json.Marshal(o)
		if err != nil {
			logger.Warnf("Couldn't json encode object")
		}
		return string(data)
	}

	return arg.String()
}

func (l *Log) logTrace(call goja.FunctionCall) goja.Value {
	logger.Trace(pluginPrefix + l.argToString(call))
	return nil
}

func (l *Log) logDebug(call goja.FunctionCall) goja.Value {
	logger.Debug(pluginPrefix + l.argToString(call))
	return nil
}

func (l *Log) logInfo(call goja.FunctionCall) goja.Value {
	logger.Info(pluginPrefix + l.argToString(call))
	return nil
}

func (l *Log) logWarn(call goja.FunctionCall) goja.Value {
	logger.Warn(pluginPrefix + l.argToString(call))
	return nil
}

func (l *Log) logError(call goja.FunctionCall) goja.Value {
	logger.Error(pluginPrefix + l.argToString(call))
	return nil
}

// Progress logs the current progress value. The progress value should be
// between 0 and 1.0 inclusively, with 1 representing that the task is
// complete. Values outside of this range will be clamp to be within it.
func (l *Log) logProgress(value float64) {
	value = math.Min(math.Max(0, value), 1)
	l.Progress <- value
}

func (l *Log) AddToVM(globalName string, vm *VM) error {
	log := vm.NewObject()
	if err := SetAll(log,
		ObjectValueDef{"Trace", l.logTrace},
		ObjectValueDef{"Debug", l.logDebug},
		ObjectValueDef{"Info", l.logInfo},
		ObjectValueDef{"Warn", l.logWarn},
		ObjectValueDef{"Error", l.logError},
		ObjectValueDef{"Progress", l.logProgress},
	); err != nil {
		return err
	}

	if err := vm.Set(globalName, log); err != nil {
		return fmt.Errorf("unable to set log: %w", err)
	}

	return nil
}
