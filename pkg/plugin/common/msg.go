package common

import (
	"encoding/json"
	"os"
)

type PluginInput struct {
	ServerPort int `json:"server_port"`
}

type PluginOutput struct {
	Error  string `json:"error"`
	Output string `json:"output"`
	Log    string `json:"log"`
}

func (o PluginOutput) Dispatch() {
	ret, _ := json.Marshal(o)
	os.Stdout.Write(ret)
}

func Error(err error) {
	o := PluginOutput{
		Error: err.Error(),
	}

	o.Dispatch()
	os.Exit(1)
}
