package plugin

import (
	"context"
	"errors"
	"path/filepath"
	"sync"

	"github.com/robertkrimen/otto"
	"github.com/stashapp/stash/pkg/logger"
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

type testObject struct{}

func (o *testObject) Foo(ctx context.Context) {
	logger.Info(ctx.Value("foo"))
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

	vm.Set("api", t.api)
	vm.Set("ctx", ctx)
	vm.Set("test", &testObject{})
	_, err = vm.Run(script)

	return err
}

func (t *jsPluginTask) Wait() {
}

func (t *jsPluginTask) Stop() error {
	return nil
}
