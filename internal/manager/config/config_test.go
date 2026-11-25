package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_SetPassword(t *testing.T) {
	i := InitializeEmpty()

	pw := strings.Repeat("a", maxPasswordBytes)
	assert.NoError(t, i.SetPassword(pw))
	assert.NotEmpty(t, i.GetPasswordHash())

	// error must not overwrite the existing hash
	before := i.GetPasswordHash()
	assert.Error(t, i.SetPassword(pw+"a"))
	assert.Equal(t, before, i.GetPasswordHash())

	// 36 two-byte runes = 72 bytes ok; 37 = 74 bytes rejected
	assert.NoError(t, i.SetPassword(strings.Repeat("é", 36)))
	assert.Error(t, i.SetPassword(strings.Repeat("é", 37)))

	assert.NoError(t, i.SetPassword(""))
	assert.Empty(t, i.GetPasswordHash())
}

func TestConfig_GetAllPluginConfiguration(t *testing.T) {
	i := InitializeEmpty()

	assert.Equal(t, i.GetAllPluginConfiguration(), map[string]map[string]interface{}{})

	i.SetPluginConfiguration("plugin1", map[string]interface{}{"key1": "value1"})

	assert.Equal(t, map[string]map[string]interface{}{
		"plugin1": {"key1": "value1"},
	}, i.GetAllPluginConfiguration())

	i.SetPluginConfiguration("plugin2", map[string]interface{}{"key2": "value2"})

	assert.Equal(t, map[string]map[string]interface{}{
		"plugin1": {"key1": "value1"},
		"plugin2": {"key2": "value2"},
	}, i.GetAllPluginConfiguration())

	// ensure SetPluginConfiguration overwrites existing configuration
	i.SetPluginConfiguration("plugin2", map[string]interface{}{"key3": "value3"})

	assert.Equal(t, map[string]map[string]interface{}{
		"plugin1": {"key1": "value1"},
		"plugin2": {"key3": "value3"},
	}, i.GetAllPluginConfiguration())
}
