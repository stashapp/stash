package common

import "net/http"

// StashServerConnection represents the connection details needed for a
// plugin instance to connect to its parent stash server.
type StashServerConnection struct {
	// http or https
	Scheme string

	Port int

	// Cookie for authentication purposes
	SessionCookie *http.Cookie
}

// PluginKeyValue represents a key/value pair for sending arguments to
// plugin operations.
type PluginKeyValue struct {
	Key   string          `json:"key"`
	Value *PluginArgValue `json:"value"`
}

// PluginArgValue represents a single value parameter for plugin operations.PluginArgValue
// Only one of its fields should be non-nil.
type PluginArgValue struct {
	Str *string           `json:"str"`
	I   *int              `json:"i"`
	B   *bool             `json:"b"`
	F   *float64          `json:"f"`
	O   []*PluginKeyValue `json:"o"`
	A   []*PluginArgValue `json:"a"`
}

// String returns the string field or an empty string if the string field is
// nil
func (v PluginArgValue) String() string {
	var ret string
	if v.Str != nil {
		ret = *v.Str
	}

	return ret
}

// Int returns the int field or 0 if the int field is nil
func (v PluginArgValue) Int() int {
	var ret int
	if v.I != nil {
		ret = *v.I
	}

	return ret
}

// Bool returns the boolean field or false if the boolean field is nil
func (v PluginArgValue) Bool() bool {
	var ret bool
	if v.B != nil {
		ret = *v.B
	}

	return ret
}

// Float returns the float field or 0 if the float field is nil
func (v PluginArgValue) Float() float64 {
	var ret float64
	if v.F != nil {
		ret = *v.F
	}

	return ret
}

// PluginInput is the data structure that is sent to plugin instances when they
// are spawned.
type PluginInput struct {
	// Server details to connect to the stash server.
	ServerConnection StashServerConnection `json:"server_connection"`

	// Arguments to the plugin operation.
	Args []*PluginKeyValue `json:"args"`
}

// GetValue gets the PluginArgValue whose key matches the provided name.
// Returns nil if not found. This will be encoded as a JSON string.
func GetValue(keyValues []*PluginKeyValue, name string) *PluginArgValue {
	for _, v := range keyValues {
		if name == v.Key {
			return v.Value
		}
	}

	return nil
}

// PluginOutput is the data structure that is expected to be output by plugin
// processes when execution has concluded. It is expected that this data will
// be encoded as JSON.
type PluginOutput struct {
	Error  *string `json:"error"`
	Output *string `json:"output"`
}

// SetError is a convenience method that sets the Error field based on the
// provided error.
func (o *PluginOutput) SetError(err error) {
	errStr := err.Error()
	o.Error = &errStr
}
