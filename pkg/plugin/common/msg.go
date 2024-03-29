package common

import "net/http"

const (
	HookContextKey = "hookContext"
)

// StashServerConnection represents the connection details needed for a
// plugin instance to connect to its parent stash server.
type StashServerConnection struct {
	// http or https
	Scheme string
	Host   string
	Port   int

	// Cookie for authentication purposes
	SessionCookie *http.Cookie

	// Dir specifies the directory containing the stash server's configuration
	// file.
	Dir string

	// PluginDir specifies the directory containing the plugin configuration
	// file.
	PluginDir string
}

// PluginArgValue represents a single value parameter for plugin operations.
type PluginArgValue interface{}

// ArgsMap is a map of argument key to value.
type ArgsMap map[string]PluginArgValue

// String returns the string field or an empty string if the string field is
// nil
func (m ArgsMap) String(key string) string {
	v, found := m[key]
	var ret string
	if !found {
		return ret
	}
	ret, _ = v.(string)
	return ret
}

// Int returns the int field or 0 if the int field is nil
func (m ArgsMap) Int(key string) int {
	v, found := m[key]
	var ret int
	if !found {
		return ret
	}
	ret, _ = v.(int)
	return ret
}

// Bool returns the boolean field or false if the boolean field is nil
func (m ArgsMap) Bool(key string) bool {
	v, found := m[key]
	var ret bool
	if !found {
		return ret
	}
	ret, _ = v.(bool)
	return ret
}

// Float returns the float field or 0 if the float field is nil
func (m ArgsMap) Float(key string) float64 {
	v, found := m[key]
	var ret float64
	if !found {
		return ret
	}
	ret, _ = v.(float64)
	return ret
}

func (m ArgsMap) ToMap() map[string]interface{} {
	ret := make(map[string]interface{})
	for k, v := range m {
		ret[k] = v
	}
	return ret
}

// PluginInput is the data structure that is sent to plugin instances when they
// are spawned.
type PluginInput struct {
	// Server details to connect to the stash server.
	ServerConnection StashServerConnection `json:"server_connection"`

	// Arguments to the plugin operation.
	Args ArgsMap `json:"args"`
}

// PluginOutput is the data structure that is expected to be output by plugin
// processes when execution has concluded. It is expected that this data will
// be encoded as JSON.
type PluginOutput struct {
	Error  *string     `json:"error"`
	Output interface{} `json:"output"`
}

// SetError is a convenience method that sets the Error field based on the
// provided error.
func (o *PluginOutput) SetError(err error) {
	errStr := err.Error()
	o.Error = &errStr
}

// HookContext is passed as a PluginArgValue and indicates what hook triggered
// this plugin task.
type HookContext struct {
	ID          int         `json:"id,omitempty"`
	Type        string      `json:"type"`
	Input       interface{} `json:"input"`
	InputFields []string    `json:"inputFields,omitempty"`
}
