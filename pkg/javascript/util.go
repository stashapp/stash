package javascript

import (
	"fmt"
	"time"
)

type Util struct{}

func (u *Util) sleepFunc(ms int64) {
	time.Sleep(time.Millisecond * time.Duration(ms))
}

func (u *Util) AddToVM(globalName string, vm *VM) error {
	util := vm.NewObject()
	if err := util.Set("Sleep", u.sleepFunc); err != nil {
		return fmt.Errorf("unable to set sleep func: %w", err)
	}

	if err := vm.Set(globalName, util); err != nil {
		return fmt.Errorf("unable to set util: %w", err)
	}

	return nil
}
