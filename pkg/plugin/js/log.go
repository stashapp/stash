package js

import (
	"encoding/json"
	"math"

	"github.com/robertkrimen/otto"
	"github.com/stashapp/stash/pkg/logger"
)

func argToString(call otto.FunctionCall) string {
	arg := call.Argument(0)
	if arg.IsObject() {
		o, _ := arg.Export()
		data, _ := json.Marshal(o)
		return string(data)
	}

	return arg.String()
}

func logTrace(call otto.FunctionCall) otto.Value {
	logger.Trace(argToString(call))
	return otto.UndefinedValue()
}

func logDebug(call otto.FunctionCall) otto.Value {
	logger.Debug(argToString(call))
	return otto.UndefinedValue()
}

func logInfo(call otto.FunctionCall) otto.Value {
	logger.Info(argToString(call))
	return otto.UndefinedValue()
}

func logWarn(call otto.FunctionCall) otto.Value {
	logger.Warn(argToString(call))
	return otto.UndefinedValue()
}

func logError(call otto.FunctionCall) otto.Value {
	logger.Error(argToString(call))
	return otto.UndefinedValue()
}

// Progress logs the current progress value. The progress value should be
// between 0 and 1.0 inclusively, with 1 representing that the task is
// complete. Values outside of this range will be clamp to be within it.
func logProgressFunc(c chan float64) func(call otto.FunctionCall) otto.Value {
	return func(call otto.FunctionCall) otto.Value {
		arg := call.Argument(0)
		if !arg.IsNumber() {
			return otto.UndefinedValue()
		}

		progress, _ := arg.ToFloat()
		progress = math.Min(math.Max(0, progress), 1)
		c <- progress

		return otto.UndefinedValue()
	}
}

func AddLogAPI(vm *otto.Otto, progress chan float64) {
	log, _ := vm.Object("({})")
	log.Set("Trace", logTrace)
	log.Set("Debug", logDebug)
	log.Set("Info", logInfo)
	log.Set("Warn", logWarn)
	log.Set("Error", logError)
	log.Set("Progress", logProgressFunc(progress))

	vm.Set("log", log)
}
