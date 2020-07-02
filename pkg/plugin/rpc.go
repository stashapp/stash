package plugin

import (
	"fmt"
	"net/rpc/jsonrpc"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/natefinch/pie"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin/common"
)

func executeRPC(operation *PluginOperationConfig, args []*models.OperationArgInput) common.PluginOutput {
	command := operation.Exec
	if len(command) == 0 {
		return makeErrorOutput(fmt.Errorf("empty exec value in operation %s", operation.Name))
	}

	// TODO - this should be the plugin config path, since it may be in a subdir
	_, err := exec.LookPath(command[0])
	if err != nil {
		// change command to use absolute path
		pluginPath := config.GetPluginsPath()
		command[0] = filepath.Join(pluginPath, command[0])
	}

	client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, os.Stderr, command[0], command[1:]...)
	if err != nil {
		return makeErrorOutput(err)
	}
	defer client.Close()

	iface := common.PluginClient{
		Client: client,
	}

	input := common.PluginInput{
		ServerPort: config.GetPort(),
	}

	output := common.PluginOutput{}
	err = iface.Run(input, &output)
	if err != nil {
		return makeErrorOutput(err)
	}

	return output
}
