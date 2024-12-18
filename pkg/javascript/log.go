package javascript

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"

	"github.com/dop251/goja"
	"github.com/stashapp/stash/pkg/logger"
)

// Log provides log wrappers for usable from the JS VM.
type Log struct {
	// Logger is the LoggerImpl to forward log messages to.
	Logger logger.LoggerImpl
	// Prefix is the prefix to prepend to log messages.
	Prefix string
	// ProgressChan is a channel that receives float64s indicating the current progress of an operation.
	ProgressChan chan float64
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
	l.Logger.Trace(l.Prefix, l.argToString(call))
	return nil
}

func (l *Log) logDebug(call goja.FunctionCall) goja.Value {
	l.Logger.Debug(l.Prefix, l.argToString(call))
	return nil
}

func (l *Log) logInfo(call goja.FunctionCall) goja.Value {
	l.Logger.Info(l.Prefix, l.argToString(call))
	return nil
}

func (l *Log) logWarn(call goja.FunctionCall) goja.Value {
	l.Logger.Warn(l.Prefix, l.argToString(call))
	return nil
}

func (l *Log) logError(call goja.FunctionCall) goja.Value {
	l.Logger.Error(l.Prefix, l.argToString(call))
	return nil
}

// Progress logs the current progress value. The progress value should be
// between 0 and 1.0 inclusively, with 1 representing that the task is
// complete. Values outside of this range will be clamp to be within it.
func (l *Log) logProgress(value float64) {
	value = math.Min(math.Max(0, value), 1)
	l.ProgressChan <- value
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
