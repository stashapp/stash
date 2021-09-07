package js

import (
	"fmt"
	"time"

	"github.com/robertkrimen/otto"
)

func sleepFunc(call otto.FunctionCall) otto.Value {
	arg := call.Argument(0)
	ms, _ := arg.ToInteger()

	time.Sleep(time.Millisecond * time.Duration(ms))
	return otto.UndefinedValue()
}

func AddUtilAPI(vm *otto.Otto) error {
	util, _ := vm.Object("({})")
	err := util.Set("Sleep", sleepFunc)
	if err != nil {
		return fmt.Errorf("unable to set sleep func: %w", err)
	}

	err = vm.Set("util", util)
	if err != nil {
		return fmt.Errorf("unable to set util: %w", err)
	}

	return nil
}
