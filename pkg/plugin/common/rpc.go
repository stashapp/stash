package common

import (
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/natefinch/pie"
)

type RPCRunner interface {
	Run(input PluginInput, output *PluginOutput) error
}

func ServePlugin(iface RPCRunner) error {
	p := pie.NewProvider()
	if err := p.RegisterName("RPCRunner", iface); err != nil {
		return err
	}

	p.ServeCodec(jsonrpc.NewServerCodec)
	return nil
}

type PluginClient struct {
	Client *rpc.Client
}

func (p PluginClient) Run(input PluginInput, output *PluginOutput) error {
	return p.Client.Call("RPCRunner.Run", input, output)
}
