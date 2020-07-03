package common

import (
	"net/rpc/jsonrpc"

	"github.com/natefinch/pie"
)

type RPCRunner interface {
	Run(input PluginInput, output *PluginOutput) error
	Stop(input struct{}, output *bool) error
}

func ServePlugin(iface RPCRunner) error {
	p := pie.NewProvider()
	if err := p.RegisterName("RPCRunner", iface); err != nil {
		return err
	}

	p.ServeCodec(jsonrpc.NewServerCodec)
	return nil
}
