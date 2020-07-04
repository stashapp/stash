package plugin

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin/common"
)

// PluginCache stores plugin details.
type PluginCache struct {
	path    string
	plugins []Config
}

// NewPluginCache returns a new PluginCache loading plugin configurations
// from the provided plugin path. It returns an new instance and an error
// if the plugin directory could not be loaded.
//
// Plugins configurations are loaded from yml files in the provided plugin
// directory and any subdirectories.
func NewPluginCache(pluginPath string) (*PluginCache, error) {
	plugins, err := loadPlugins(pluginPath)
	if err != nil {
		return nil, err
	}

	return &PluginCache{
		path:    pluginPath,
		plugins: plugins,
	}, nil
}

// ReloadPlugins clears the plugin cache and reloads from the plugin path.
// In the event of an error during loading, the cache will be left empty.
func (c *PluginCache) ReloadPlugins() error {
	c.plugins = nil
	plugins, err := loadPlugins(c.path)
	if err != nil {
		return err
	}

	c.plugins = plugins
	return nil
}

func loadPlugins(path string) ([]Config, error) {
	plugins := make([]Config, 0)

	logger.Debugf("Reading plugin configs from %s", path)
	pluginFiles := []string{}
	err := filepath.Walk(path, func(fp string, f os.FileInfo, err error) error {
		if filepath.Ext(fp) == ".yml" {
			pluginFiles = append(pluginFiles, fp)
		}
		return nil
	})

	if err != nil {

		return nil, err
	}

	for _, file := range pluginFiles {
		plugin, err := loadPluginFromYAMLFile(file)
		if err != nil {
			logger.Errorf("Error loading plugin %s: %s", file, err.Error())
		} else {
			plugins = append(plugins, *plugin)
		}
	}

	return plugins, nil
}

// ListPlugins returns plugin details for all of the loaded plugins.
func (c PluginCache) ListPlugins() []*models.Plugin {
	var ret []*models.Plugin
	for _, s := range c.plugins {
		ret = append(ret, s.toPlugin())
	}

	return ret
}

// ListPluginTasks returns all runnable plugin tasks in all loaded plugins.
func (c PluginCache) ListPluginTasks() []*models.PluginTask {
	var ret []*models.PluginTask
	for _, s := range c.plugins {
		ret = append(ret, s.getPluginTasks()...)
	}

	return ret
}

// RunPluginOperation runs the plugin operation for the pluginID and operation
// name provided. Returns an error if the plugin or the operation could not be
// resolved.
func (c PluginCache) RunPluginOperation(pluginID string, operationName string, serverConnection common.StashServerConnection, args []*models.PluginArgInput) (PluginTaskManager, error) {
	// find the plugin and operation
	plugin := c.getPlugin(pluginID)

	if plugin == nil {
		return nil, fmt.Errorf("no plugin with ID %s", pluginID)
	}

	operation := plugin.getTask(operationName)
	if operation == nil {
		return nil, fmt.Errorf("no task with name %s in plugin %s", operationName, plugin.getName())
	}

	ret := newRPCPluginTask(operation, args, serverConnection)
	if err := ret.Start(); err != nil {
		return nil, err
	}

	return ret, nil
}

func (c PluginCache) getPlugin(pluginID string) *Config {
	for _, s := range c.plugins {
		if s.id == pluginID {
			return &s
		}
	}

	return nil
}
