package plugin

import (
	"context"
	"errors"
	"path/filepath"
	"sync"

	"github.com/robertkrimen/otto"
	"github.com/stashapp/stash/pkg/plugin/common"
)

type jsTaskBuilder struct{}

func (*jsTaskBuilder) build(task pluginTask) Task {
	return &jsPluginTask{
		pluginTask: task,
	}
}

type jsPluginTask struct {
	pluginTask

	started   bool
	waitGroup sync.WaitGroup
}

func (t *jsPluginTask) makeOutput(o otto.Value) {
	t.result = &common.PluginOutput{}

	asObj := o.Object()
	if asObj == nil {
		return
	}

	t.result.Output, _ = asObj.Get("Output")
	err, _ := asObj.Get("Error")
	if !err.IsUndefined() {
		errStr := err.String()
		t.result.Error = &errStr
	}
}

func (t *jsPluginTask) Start() error {
	if t.started {
		return errors.New("task already started")
	}

	t.started = true

	if len(t.plugin.Exec) == 0 {
		return errors.New("no script specified in exec")
	}

	scriptFile := t.plugin.Exec[0]

	// otto doesn't like interfaces much. Need to provide context with
	// a value, otherwise it resolves to an int.
	ctx := context.WithValue(context.Background(), "foo", "bar")

	vm := otto.New()
	pluginPath := t.plugin.getConfigPath()
	script, err := vm.Compile(filepath.Join(pluginPath, scriptFile), nil)
	if err != nil {
		return err
	}

	input := t.buildPluginInput()

	vm.Set("input", input)
	vm.Set("api", &jsAPI{r: t.api, ctx: ctx})
	// TODO - vm.Set("log")
	output, err := vm.Run(script)

	t.makeOutput(output)

	return err
}

func (t *jsPluginTask) Wait() {
}

func (t *jsPluginTask) Stop() error {
	return nil
}
