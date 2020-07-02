package common

type StashServerProvider interface {
	GetPort() int
}

type PluginKeyValue struct {
	Key   string          `json:"key"`
	Value *PluginArgValue `json:"value"`
}

type PluginArgValue struct {
	Str *string           `json:"str"`
	I   *int              `json:"i"`
	B   *bool             `json:"b"`
	F   *float64          `json:"f"`
	O   []*PluginKeyValue `json:"o"`
	A   []*PluginArgValue `json:"a"`
}

type PluginInput struct {
	ServerPort int              `json:"server_port"`
	Args       []PluginKeyValue `json:"args"`
}

func (i PluginInput) GetPort() int {
	return i.ServerPort
}

type PluginOutput struct {
	Error  *string         `json:"error"`
	Output *PluginArgValue `json:"output"`
}
