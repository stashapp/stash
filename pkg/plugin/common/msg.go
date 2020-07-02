package common

import (
	"encoding/json"
	"os"
)

const InterfaceJsonV1 = "json.v1"

type StashServerProvider interface {
	GetPort() int
}

type PluginInput struct {
	ServerPort int `json:"server_port"`
}

func (i PluginInput) GetPort() int {
	return i.ServerPort
}

func ReadPluginInput() (*PluginInput, error) {
	out := PluginInput{}
	decodeErr := json.NewDecoder(os.Stdin).Decode(&out)

	return &out, decodeErr
}

type PluginOutput struct {
	Error  *string `json:"error"`
	Output string  `json:"output"`
}

func (o PluginOutput) Dispatch() {
	ret, _ := json.Marshal(o)
	os.Stdout.Write(ret)
}

func Error(err error) {
	str := err.Error()
	o := PluginOutput{
		Error: &str,
	}

	o.Dispatch()
	os.Exit(1)
}
