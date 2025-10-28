package plugin

type OperationInput map[string]interface{}

type PluginArgInput struct {
	Key   string            `json:"key"`
	Value *PluginValueInput `json:"value"`
}

type PluginValueInput struct {
	Str *string             `json:"str"`
	I   *int                `json:"i"`
	B   *bool               `json:"b"`
	F   *float64            `json:"f"`
	O   []*PluginArgInput   `json:"o"`
	A   []*PluginValueInput `json:"a"`
}

func applyDefaultArgs(args OperationInput, defaultArgs map[string]string) {
	for k, v := range defaultArgs {
		_, found := args[k]
		if !found {
			args[k] = v
		}
	}
}
