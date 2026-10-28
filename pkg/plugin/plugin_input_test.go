package plugin

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stashapp/stash/pkg/javascript"
	"github.com/stashapp/stash/pkg/plugin/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type pluginInputTestConfig struct {
	mu       sync.RWMutex
	settings map[string]map[string]interface{}
}

func (c *pluginInputTestConfig) GetHost() string          { return "localhost" }
func (c *pluginInputTestConfig) GetPort() int             { return 9999 }
func (c *pluginInputTestConfig) GetConfigPathAbs() string { return "" }
func (c *pluginInputTestConfig) HasTLSConfig() bool       { return false }
func (c *pluginInputTestConfig) GetPluginsPath() string   { return "" }
func (c *pluginInputTestConfig) GetDisabledPlugins() []string {
	return nil
}
func (c *pluginInputTestConfig) GetPluginConfiguration(pluginID string) map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.settings[pluginID]
}
func (c *pluginInputTestConfig) GetPythonPath() string { return "" }

func (c *pluginInputTestConfig) setPluginConfiguration(pluginID string, settings map[string]interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.settings == nil {
		c.settings = make(map[string]map[string]interface{})
	}
	c.settings[pluginID] = settings
}

func TestBuildPluginInputIncludesPluginSettings(t *testing.T) {
	settings := map[string]interface{}{
		"enabled": true,
		"count":   3,
		"name":    "example",
	}
	cache := NewCache(&pluginInputTestConfig{settings: map[string]map[string]interface{}{
		"test-plugin": settings,
	}})
	plugin := &Config{id: "test-plugin", path: "test-plugin.yml"}
	operation := &OperationConfig{DefaultArgs: map[string]string{"default": "value"}}

	input := cache.buildPluginInput(plugin, operation, common.StashServerConnection{}, OperationInput{"provided": "value"})

	assert.Equal(t, settings, input.Settings)
	assert.Equal(t, "value", input.Args.String("default"))
	assert.Equal(t, "value", input.Args.String("provided"))
	assert.Equal(t, plugin.getConfigPath(), input.ServerConnection.PluginDir)
}

func TestBuildPluginInputUsesEmptySettingsForUnconfiguredPlugin(t *testing.T) {
	cache := NewCache(&pluginInputTestConfig{})

	input := cache.buildPluginInput(&Config{id: "test-plugin"}, nil, common.StashServerConnection{}, nil)

	require.NotNil(t, input.Settings)
	assert.Empty(t, input.Settings)
}

func TestBuildPluginInputSettingsAreDetachedFromConfiguration(t *testing.T) {
	config := &pluginInputTestConfig{settings: map[string]map[string]interface{}{
		"test-plugin": {"enabled": true},
	}}
	cache := NewCache(config)

	input := cache.buildPluginInput(&Config{id: "test-plugin"}, nil, common.StashServerConnection{}, nil)
	input.Settings["enabled"] = false
	input.Settings["new"] = "value"

	stored := config.GetPluginConfiguration("test-plugin")
	assert.Equal(t, true, stored["enabled"])
	assert.NotContains(t, stored, "new")
}

func TestBuildPluginInputUsesLatestSettingsWithoutChangingExistingInput(t *testing.T) {
	config := &pluginInputTestConfig{}
	config.setPluginConfiguration("test-plugin", map[string]interface{}{"version": 1})
	cache := NewCache(config)
	plugin := &Config{id: "test-plugin"}

	first := cache.buildPluginInput(plugin, nil, common.StashServerConnection{}, nil)
	config.setPluginConfiguration("test-plugin", map[string]interface{}{"version": 2})
	second := cache.buildPluginInput(plugin, nil, common.StashServerConnection{}, nil)

	assert.Equal(t, 1, first.Settings["version"])
	assert.Equal(t, 2, second.Settings["version"])
}

func TestBuildPluginInputOnlyIncludesRequestedPluginSettings(t *testing.T) {
	config := &pluginInputTestConfig{settings: map[string]map[string]interface{}{
		"first-plugin":  {"name": "first"},
		"second-plugin": {"name": "second"},
	}}
	cache := NewCache(config)

	input := cache.buildPluginInput(&Config{id: "first-plugin"}, nil, common.StashServerConnection{}, nil)

	assert.Equal(t, map[string]interface{}{"name": "first"}, input.Settings)
}

// This test should be run with -race. Each input must be a coherent snapshot,
// even while later invocations begin using updated settings.
func TestBuildPluginInputConcurrentSettingsUpdates(t *testing.T) {
	config := &pluginInputTestConfig{}
	config.setPluginConfiguration("test-plugin", map[string]interface{}{
		"version": 0,
		"label":   "version-0",
	})
	cache := NewCache(config)
	plugin := &Config{id: "test-plugin"}

	const iterations = 200
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 1; i <= iterations; i++ {
			config.setPluginConfiguration("test-plugin", map[string]interface{}{
				"version": i,
				"label":   fmt.Sprintf("version-%d", i),
			})
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			input := cache.buildPluginInput(plugin, nil, common.StashServerConnection{}, nil)
			version := input.Settings["version"].(int)
			assert.Equal(t, fmt.Sprintf("version-%d", version), input.Settings["label"])
		}
	}()

	wg.Wait()
}

func TestJSPluginInputIncludesPluginSettings(t *testing.T) {
	task := &jsPluginTask{
		pluginTask: pluginTask{
			plugin: &Config{Name: "Test plugin"},
			input: common.PluginInput{
				Args: common.ArgsMap{},
				Settings: map[string]interface{}{
					"enabled": true,
					"name":    "example",
				},
			},
		},
		vm: javascript.NewVM(),
	}

	require.NoError(t, task.initVM())
	value, err := task.vm.RunString(`JSON.stringify(input.Settings)`)
	require.NoError(t, err)
	assert.JSONEq(t, `{"enabled":true,"name":"example"}`, value.String())
}

func TestJSPluginCannotMutateStoredPluginSettings(t *testing.T) {
	config := &pluginInputTestConfig{settings: map[string]map[string]interface{}{
		"test-plugin": {"enabled": true},
	}}
	cache := NewCache(config)
	task := &jsPluginTask{
		pluginTask: pluginTask{
			plugin: &Config{id: "test-plugin", Name: "Test plugin"},
			input: cache.buildPluginInput(
				&Config{id: "test-plugin"},
				nil,
				common.StashServerConnection{},
				nil,
			),
		},
		vm: javascript.NewVM(),
	}

	require.NoError(t, task.initVM())
	_, err := task.vm.RunString(`input.Settings.enabled = false; input.Settings.added = "value"`)
	require.NoError(t, err)

	stored := config.GetPluginConfiguration("test-plugin")
	assert.Equal(t, true, stored["enabled"])
	assert.NotContains(t, stored, "added")
}
