package common

import (
	"net/rpc/jsonrpc"

	"github.com/natefinch/pie"
)

// RPCRunner is the interface that RPC plugins are expected to fulfil.
type RPCRunner interface {
	// Perform the operation, using the provided input and populating the
	// output object.
	Run(input PluginInput, output *PluginOutput) error

	// Stop any running operations, if possible. No input is sent and any
	// output is ignored.
	Stop(input struct{}, output *bool) error
}

// ServePlugin is used by plugin instances to serve the plugin via RPC, using
// the provided RPCRunner interface.
func ServePlugin(iface RPCRunner) error {
	p := pie.NewProvider()
	if err := p.RegisterName("RPCRunner", iface); err != nil {
		return err
	}

	p.ServeCodec(jsonrpc.NewServerCodec)
	return nil
}
