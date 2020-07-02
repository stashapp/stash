package plugin

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
)

var plugins []PluginConfig

func loadPlugins() ([]PluginConfig, error) {
	if plugins != nil {
		return plugins, nil
	}

	path := config.GetPluginsPath()
	plugins = make([]PluginConfig, 0)

	logger.Debugf("Reading plugin configs from %s", path)
	pluginFiles := []string{}
	err := filepath.Walk(path, func(fp string, f os.FileInfo, err error) error {
		if filepath.Ext(fp) == ".yml" {
			pluginFiles = append(pluginFiles, fp)
		}
		return nil
	})

	if err != nil {
		logger.Errorf("Error reading plugin configs: %s", err.Error())
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

func ReloadPlugins() error {
	plugins = nil
	_, err := loadPlugins()
	return err
}

func ListPlugins() ([]*models.Plugin, error) {
	// read plugin config files from the directory and cache
	plugins, err := loadPlugins()

	if err != nil {
		return nil, err
	}

	var ret []*models.Plugin
	for _, s := range plugins {
		ret = append(ret, s.toPlugin())
	}

	return ret, nil
}

func ListPluginTasks() ([]*models.PluginTask, error) {
	// read plugin config files from the directory and cache
	plugins, err := loadPlugins()

	if err != nil {
		return nil, err
	}

	var ret []*models.PluginTask
	for _, s := range plugins {
		ret = append(ret, s.getPluginTasks()...)
	}

	return ret, nil
}

func getPlugin(pluginID string) (*PluginConfig, error) {
	// read plugin config files from the directory and cache
	plugins, err := loadPlugins()

	if err != nil {
		return nil, err
	}

	for _, s := range plugins {
		if s.ID == pluginID {
			return &s, nil
		}
	}

	return nil, nil
}

func RunPluginOperation(pluginID string, operationName string, args []*models.PluginArgInput) error {
	// find the plugin and operation
	plugin, err := getPlugin(pluginID)

	if err != nil {
		return err
	}

	if plugin == nil {
		return fmt.Errorf("no plugin with ID %s", pluginID)
	}

	operation := plugin.getTask(operationName)
	if operation == nil {
		return fmt.Errorf("no task with name %s in plugin %s", operationName, plugin.getName())
	}

	_, err = executeRPC(operation, args)

	return err
}
