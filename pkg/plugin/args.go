package plugin

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

func findArg(args []*PluginArgInput, name string) *PluginArgInput {
	for _, v := range args {
		if v.Key == name {
			return v
		}
	}

	return nil
}

func applyDefaultArgs(args []*PluginArgInput, defaultArgs map[string]string) []*PluginArgInput {
	for k, v := range defaultArgs {
		if arg := findArg(args, k); arg == nil {
			v := v // Copy v, because it's being exported out of the loop
			args = append(args, &PluginArgInput{
				Key: k,
				Value: &PluginValueInput{
					Str: &v,
				},
			})
		}
	}

	return args
}
