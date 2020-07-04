package plugin

import (
	"fmt"
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

	// The communication interface used when communicating with the spawned
	// plugin process. Defaults to 'raw' if not provided.
	Interface interfaceEnum `yaml:"interface"`

	// The command to execute for the operations in this plugin. The first
	// element should be the program name, and subsequent elements are passed
	// as arguments.
	//
	// Note: the execution process will search the path for the program,
	// then will attempt to find the program in the plugins
	// directory. The exe extension is not necessary on Windows platforms.
	// The current working directory is set to that of the stash process.
	Exec []string `yaml:"exec,flow"`

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

func (c Config) getExecCommand(task *OperationConfig) []string {
	ret := c.Exec

	ret = append(ret, task.ExecArgs...)
	return ret
}

type interfaceEnum string

// Valid interfaceEnum values
const (
	// InterfaceEnumRPC indicates that the plugin uses the RPCRunner interface
	// declared in common/rpc.go.
	InterfaceEnumRPC interfaceEnum = "rpc"

	// InterfaceEnumRaw indidates that stdout will be treated as raw output.
	InterfaceEnumRaw interfaceEnum = "raw"
)

func (i interfaceEnum) Valid() bool {
	return i == InterfaceEnumRPC || i == InterfaceEnumRaw
}

func (i *interfaceEnum) getTaskBuilder() taskBuilder {
	if !i.Valid() || *i == InterfaceEnumRaw {
		// TODO
		return nil
	}
	if *i == InterfaceEnumRPC {
		return &rpcTaskBuilder{}
	}

	// shouldn't happen
	return nil
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

	// A list of arguments that will be appended to the plugin's Exec arguments
	// when executing this operation.
	ExecArgs []string

	// A map of argument keys to their default values. The default value is
	// used if the applicable argument is not provided during the operation
	// call.
	DefaultArgs map[string]string `yaml:"defaultArgs"`
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

	if ret.Interface == "" {
		ret.Interface = InterfaceEnumRaw
	}

	if !ret.Interface.Valid() {
		return nil, fmt.Errorf("invalid interface type %s", ret.Interface)
	}

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
