package plugin

import (
	"errors"
	"fmt"
	"io"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/natefinch/pie"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/plugin/common"
)

type rpcTaskBuilder struct{}

func (*rpcTaskBuilder) build(task pluginTask) Task {
	return &rpcPluginTask{
		pluginTask: task,
	}
}

type rpcPluginClient struct {
	Client *rpc.Client
}

func (p rpcPluginClient) Run(input common.PluginInput, output *common.PluginOutput) error {
	return p.Client.Call("RPCRunner.Run", input, output)
}

func (p rpcPluginClient) RunAsync(input common.PluginInput, output *common.PluginOutput, done chan *rpc.Call) *rpc.Call {
	return p.Client.Go("RPCRunner.Run", input, output, done)
}

func (p rpcPluginClient) Stop() error {
	var resp interface{}
	return p.Client.Call("RPCRunner.Stop", nil, &resp)
}

type rpcPluginTask struct {
	pluginTask

	started   bool
	client    *rpc.Client
	waitGroup sync.WaitGroup
	done      chan *rpc.Call
}

func (t *rpcPluginTask) Start() error {
	if t.started {
		return errors.New("task already started")
	}

	command := t.plugin.getExecCommand(t.operation)
	if len(command) == 0 {
		return fmt.Errorf("empty exec value in operation %s", t.operation.Name)
	}

	// TODO - this should be the plugin config path, since it may be in a subdir
	_, err := exec.LookPath(command[0])
	if err != nil {
		// change command to use absolute path
		pluginPath := config.GetPluginsPath()
		command[0] = filepath.Join(pluginPath, command[0])
	}

	pluginErrReader, pluginErrWriter := io.Pipe()

	t.client, err = pie.StartProviderCodec(jsonrpc.NewClientCodec, pluginErrWriter, command[0], command[1:]...)
	if err != nil {
		return err
	}

	go t.handleStderr(pluginErrReader)

	iface := rpcPluginClient{
		Client: t.client,
	}

	args := applyDefaultArgs(t.args, t.operation.DefaultArgs)

	input := buildPluginInput(args, t.serverConnection)

	t.done = make(chan *rpc.Call, 1)
	result := common.PluginOutput{}
	t.waitGroup.Add(1)
	iface.RunAsync(input, &result, t.done)
	go t.waitToFinish(&result)

	t.started = true
	return nil
}

func (t *rpcPluginTask) waitToFinish(result *common.PluginOutput) {
	defer t.client.Close()
	defer t.waitGroup.Done()
	<-t.done

	t.result = result
}

func (t *rpcPluginTask) Wait() {
	t.waitGroup.Wait()
}

func (t *rpcPluginTask) Stop() error {
	iface := rpcPluginClient{
		Client: t.client,
	}

	return iface.Stop()
}
