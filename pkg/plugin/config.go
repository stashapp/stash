package plugin

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/stashapp/stash/pkg/models"
	"gopkg.in/yaml.v2"
)

// Config describes the configuration for a single plugin.
type Config struct {
	id string

	// The name of the plugin. This will be displayed in the UI.
	Name string `yaml:"name"`

	// An optional description of what the plugin does.
	Description *string `yaml:"description"`

	// An optional URL for the plugin.
	URL *string `yaml:"url"`

	// An optional version string.
	Version *string `yaml:"version"`

	// The task configurations for tasks provided by this plugin.
	Tasks []*OperationConfig `yaml:"tasks"`
}

func (c Config) getPluginTasks() []*models.PluginTask {
	var ret []*models.PluginTask

	for _, o := range c.Tasks {
		ret = append(ret, &models.PluginTask{
			Name:        o.Name,
			Description: &o.Description,
			Plugin:      c.toPlugin(),
		})
	}

	return ret
}

func (c Config) getName() string {
	if c.Name != "" {
		return c.Name
	}

	return c.id
}

func (c Config) toPlugin() *models.Plugin {
	return &models.Plugin{
		ID:          c.id,
		Name:        c.getName(),
		Description: c.Description,
		URL:         c.URL,
		Version:     c.Version,
	}
}

func (c Config) getTask(name string) *OperationConfig {
	for _, o := range c.Tasks {
		if o.Name == name {
			return o
		}
	}

	return nil
}

type InterfaceEnum string

const (
	// Uses the RPCRunner interface declared in common/rpc.go
	InterfaceEnumRPC InterfaceEnum = "rpc"

	// Treats stdout as raw output
	InterfaceEnumRaw InterfaceEnum = "raw"
)

func (i InterfaceEnum) Valid() bool {
	return i == InterfaceEnumRPC || i == InterfaceEnumRaw
}

// OperationConfig describes the configuration for a single plugin operation
// provided by a plugin.
type OperationConfig struct {
	// Used to identify the operation. Must be unique within a plugin
	// configuration. This name is shown in the button for the operation
	// in the UI.
	Name string `yaml:"name"`

	// A short description of the operation. This description is shown below
	// the button in the UI.
	Description string `yaml:"description"`

	// The command to execute for the operation.
	//
	// Note: the execution process will search the path for the first element
	// in this list, then will attempt to find the element in the plugins
	// directory. The exe extension is not necessary on Windows platforms.
	// The current working directory is set to that of the stash process.
	Exec []string `yaml:"exec,flow"`

	// A map of argument keys to their default values. The default value is
	// used if the applicable argument is not provided during the operation
	// call.
	DefaultArgs map[string]string `yaml:"defaultArgs"`

	// The communication interface used when communicating with the spawned
	// plugin process.
	Interface InterfaceEnum `yaml:"interface"`
}

func loadPluginFromYAML(id string, reader io.Reader) (*Config, error) {
	ret := &Config{}

	parser := yaml.NewDecoder(reader)
	parser.SetStrict(true)
	err := parser.Decode(&ret)
	if err != nil {
		return nil, err
	}

	ret.id = id

	return ret, nil
}

func loadPluginFromYAMLFile(path string) (*Config, error) {
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
