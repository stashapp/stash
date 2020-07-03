package plugin

import (
	"bufio"
	"fmt"
	"io"
	"net/rpc/jsonrpc"
	"os/exec"
	"path/filepath"

	"github.com/natefinch/pie"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin/common"
)

func executeRPC(operation *PluginOperationConfig, args []*models.PluginArgInput, serverConnection common.StashServerConnection) (*common.PluginOutput, error) {
	command := operation.Exec
	if len(command) == 0 {
		return nil, fmt.Errorf("empty exec value in operation %s", operation.Name)
	}

	// TODO - this should be the plugin config path, since it may be in a subdir
	_, err := exec.LookPath(command[0])
	if err != nil {
		// change command to use absolute path
		pluginPath := config.GetPluginsPath()
		command[0] = filepath.Join(pluginPath, command[0])
	}

	pluginErrReader, pluginErrWriter := io.Pipe()

	client, err := pie.StartProviderCodec(jsonrpc.NewClientCodec, pluginErrWriter, command[0], command[1:]...)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	go func() {
		// pipe plugin stderr to our logging
		scanner := bufio.NewScanner(pluginErrReader)
		for scanner.Scan() {
			str := scanner.Text()
			if str != "" {
				// TODO - support logging and progress
				logger.Infof("[Plugin] %s", str)
			}
		}

		str := scanner.Text()
		if str != "" {
			// TODO - support logging and progress
			logger.Infof("[Plugin] %s", str)
		}

		pluginErrReader.Close()
	}()

	iface := common.PluginClient{
		Client: client,
	}

	args = applyDefaultArgs(args, operation.DefaultArgs)

	input := buildPluginInput(args, serverConnection)

	output := common.PluginOutput{}
	err = iface.Run(input, &output)
	if err != nil {
		return nil, err
	}

	return &output, nil
}
