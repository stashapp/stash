package js

import (
	"time"

	"github.com/robertkrimen/otto"
)

func sleepFunc(call otto.FunctionCall) otto.Value {
	arg := call.Argument(0)
	ms, _ := arg.ToInteger()

	time.Sleep(time.Millisecond * time.Duration(ms))
	return otto.UndefinedValue()
}

func AddUtilAPI(vm *otto.Otto) {
	util, _ := vm.Object("({})")
	util.Set("Sleep", sleepFunc)

	vm.Set("util", util)
}
