package javascript

import "fmt"

type console struct {
	Log
}

func (c *console) AddToVM(globalName string, vm *VM) error {
	console := vm.NewObject()
	if err := SetAll(console,
		ObjectValueDef{"log", c.logInfo},
		ObjectValueDef{"error", c.logError},
		ObjectValueDef{"warn", c.logWarn},
		ObjectValueDef{"info", c.logInfo},
		ObjectValueDef{"debug", c.logDebug},
	); err != nil {
		return err
	}

	if err := vm.Set(globalName, console); err != nil {
		return fmt.Errorf("unable to set console: %w", err)
	}

	return nil
}
