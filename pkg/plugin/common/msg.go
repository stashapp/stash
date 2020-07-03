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

func (v PluginArgValue) String() string {
	var ret string
	if v.Str != nil {
		ret = *v.Str
	}

	return ret
}

func (v PluginArgValue) Int() int {
	var ret int
	if v.I != nil {
		ret = *v.I
	}

	return ret
}

func (v PluginArgValue) Bool() bool {
	var ret bool
	if v.B != nil {
		ret = *v.B
	}

	return ret
}

func (v PluginArgValue) Float() float64 {
	var ret float64
	if v.F != nil {
		ret = *v.F
	}

	return ret
}

type PluginInput struct {
	ServerPort int               `json:"server_port"`
	Args       []*PluginKeyValue `json:"args"`
}

func (i PluginInput) GetPort() int {
	return i.ServerPort
}

func GetValue(keyValues []*PluginKeyValue, name string) *PluginArgValue {
	for _, v := range keyValues {
		if name == v.Key {
			return v.Value
		}
	}

	return nil
}

type PluginOutput struct {
	Error  *string `json:"error"`
	Output *string `json:"output"`
}
