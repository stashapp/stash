package plugin

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/stashapp/stash/pkg/models"
	"gopkg.in/yaml.v2"
)

type PluginConfig struct {
	ID         string
	Name       string                   `yaml:"name"`
	Operations []*PluginOperationConfig `yaml:"operations"`
}

func (c PluginConfig) getPluginOperations() []*models.PluginOperation {
	var ret []*models.PluginOperation

	for _, o := range c.Operations {
		ret = append(ret, &models.PluginOperation{
			PluginID:      c.ID,
			OperationName: o.Name,
		})
	}

	return ret
}

func (c PluginConfig) getName() string {
	if c.Name != "" {
		return c.Name
	}

	return c.ID
}

func (c PluginConfig) toPlugin() *models.Plugin {
	return &models.Plugin{
		ID:   c.ID,
		Name: c.getName(),
	}
}

func (c PluginConfig) getOperation(name string) *PluginOperationConfig {
	for _, o := range c.Operations {
		if o.Name == name {
			return o
		}
	}

	return nil
}

type PluginOperationConfig struct {
	Name string   `yaml:"name"`
	Exec []string `yaml:"exec,flow"`

	// task, query, mutation
	Type string `yaml:"type"`

	// communication interface used when communicating with the spawned plugin process
	Interface string `yaml:"interface"`
}

func loadPluginFromYAML(id string, reader io.Reader) (*PluginConfig, error) {
	ret := &PluginConfig{}

	parser := yaml.NewDecoder(reader)
	parser.SetStrict(true)
	err := parser.Decode(&ret)
	if err != nil {
		return nil, err
	}

	ret.ID = id

	return ret, nil
}

func loadPluginFromYAMLFile(path string) (*PluginConfig, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	// set id to the filename
	id := filepath.Base(path)
	id = id[:strings.LastIndex(id, ".")]

	return loadPluginFromYAML(id, file)
}
