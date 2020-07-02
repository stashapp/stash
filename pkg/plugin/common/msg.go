package common

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

type PluginOutput struct {
	Error  *string `json:"error"`
	Output string  `json:"output"`
}
