package js

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/robertkrimen/otto"
	"github.com/stashapp/stash/pkg/logger"
)

const pluginPrefix = "[Plugin] "

func argToString(call otto.FunctionCall) string {
	arg := call.Argument(0)
	if arg.IsObject() {
		o, _ := arg.Export()
		data, err := json.Marshal(o)
		if err != nil {
			logger.Warnf("Couldn't json encode object")
		}
		return string(data)
	}

	return arg.String()
}

func logTrace(call otto.FunctionCall) otto.Value {
	logger.Trace(pluginPrefix + argToString(call))
	return otto.UndefinedValue()
}

func logDebug(call otto.FunctionCall) otto.Value {
	logger.Debug(pluginPrefix + argToString(call))
	return otto.UndefinedValue()
}

func logInfo(call otto.FunctionCall) otto.Value {
	logger.Info(pluginPrefix + argToString(call))
	return otto.UndefinedValue()
}

func logWarn(call otto.FunctionCall) otto.Value {
	logger.Warn(pluginPrefix + argToString(call))
	return otto.UndefinedValue()
}

func logError(call otto.FunctionCall) otto.Value {
	logger.Error(pluginPrefix + argToString(call))
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

func AddLogAPI(vm *otto.Otto, progress chan float64) error {
	log, _ := vm.Object("({})")
	if err := log.Set("Trace", logTrace); err != nil {
		return fmt.Errorf("error setting Trace: %w", err)
	}
	if err := log.Set("Debug", logDebug); err != nil {
		return fmt.Errorf("error setting Debug: %w", err)
	}
	if err := log.Set("Info", logInfo); err != nil {
		return fmt.Errorf("error setting Info: %w", err)
	}
	if err := log.Set("Warn", logWarn); err != nil {
		return fmt.Errorf("error setting Warn: %w", err)
	}
	if err := log.Set("Error", logError); err != nil {
		return fmt.Errorf("error setting Error: %w", err)
	}
	if err := log.Set("Progress", logProgressFunc(progress)); err != nil {
		return fmt.Errorf("error setting Progress: %v", err)
	}
	if err := vm.Set("log", log); err != nil {
		return fmt.Errorf("unable to set log: %w", err)
	}

	return nil
}
