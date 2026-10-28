package common

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPluginInputJSONIncludesSettings(t *testing.T) {
	input := PluginInput{
		Args: ArgsMap{"mode": "test"},
		Settings: map[string]interface{}{
			"enabled": true,
			"count":   3,
			"name":    "example",
		},
	}

	encoded, err := json.Marshal(input)
	require.NoError(t, err)

	var decoded map[string]interface{}
	require.NoError(t, json.Unmarshal(encoded, &decoded))
	assert.Equal(t, map[string]interface{}{
		"enabled": true,
		"count":   float64(3),
		"name":    "example",
	}, decoded["settings"])
}

func TestPluginInputJSONUsesEmptyObjectForEmptySettings(t *testing.T) {
	input := PluginInput{
		Args:     ArgsMap{},
		Settings: map[string]interface{}{},
	}

	encoded, err := json.Marshal(input)
	require.NoError(t, err)

	var decoded map[string]interface{}
	require.NoError(t, json.Unmarshal(encoded, &decoded))
	require.NotNil(t, decoded["settings"])
	assert.Empty(t, decoded["settings"])
}
