package plugin

import (
	"fmt"
	"io"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os/exec"
	"path/filepath"

	"github.com/natefinch/pie"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin/common"
)

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

type RPCPluginTask struct {
	PluginTask

	client *rpc.Client
	done   chan *rpc.Call
}

func newRPCPluginTask(operation *PluginOperationConfig, args []*models.PluginArgInput, serverConnection common.StashServerConnection) *RPCPluginTask {
	return &RPCPluginTask{
		PluginTask: newPluginTask(operation, args, serverConnection),
	}
}

func (t *RPCPluginTask) Start() error {
	command := t.Operation.Exec
	if len(command) == 0 {
		return fmt.Errorf("empty exec value in operation %s", t.Operation.Name)
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

	go handleStderr(pluginErrReader)

	iface := rpcPluginClient{
		Client: t.client,
	}

	args := applyDefaultArgs(t.Args, t.Operation.DefaultArgs)

	input := buildPluginInput(args, t.ServerConnection)

	t.done = make(chan *rpc.Call, 1)
	t.WaitGroup.Add(1)
	iface.RunAsync(input, &t.result, t.done)
	go t.waitToFinish()
	return nil
}

func (t *RPCPluginTask) waitToFinish() {
	defer t.client.Close()
	defer t.WaitGroup.Done()
	<-t.done
}

func (t *RPCPluginTask) Wait() {
	t.WaitGroup.Wait()
}

func (t *RPCPluginTask) Stop() error {
	iface := rpcPluginClient{
		Client: t.client,
	}

	return iface.Stop()
}

func (t *RPCPluginTask) GetProgress() *int {
	return nil
}
