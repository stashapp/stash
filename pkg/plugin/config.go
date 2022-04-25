package plugin

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// Config describes the configuration for a single plugin.
type Config struct {
	id string

	// path to the configuration file
	path string

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

	// The default log level to output the plugin process's stderr stream.
	// Only used if the plugin does not encode its output using log level
	// control characters.
	// See package common/log for valid values.
	// If left unset, defaults to log.ErrorLevel.
	PluginErrLogLevel string `yaml:"errLog"`

	// The task configurations for tasks provided by this plugin.
	Tasks []*OperationConfig `yaml:"tasks"`

	// The hooks configurations for hooks registered by this plugin.
	Hooks []*HookConfig `yaml:"hooks"`
}

func (c Config) getPluginTasks(includePlugin bool) []*PluginTask {
	var ret []*PluginTask

	for _, o := range c.Tasks {
		task := &PluginTask{
			Name:        o.Name,
			Description: &o.Description,
		}

		if includePlugin {
			task.Plugin = c.toPlugin()
		}
		ret = append(ret, task)
	}

	return ret
}

func (c Config) getPluginHooks(includePlugin bool) []*PluginHook {
	var ret []*PluginHook

	for _, o := range c.Hooks {
		hook := &PluginHook{
			Name:        o.Name,
			Description: &o.Description,
			Hooks:       convertHooks(o.TriggeredBy),
		}

		if includePlugin {
			hook.Plugin = c.toPlugin()
		}
		ret = append(ret, hook)
	}

	return ret
}

func convertHooks(hooks []HookTriggerEnum) []string {
	var ret []string
	for _, h := range hooks {
		ret = append(ret, h.String())
	}

	return ret
}

func (c Config) getName() string {
	if c.Name != "" {
		return c.Name
	}

	return c.id
}

func (c Config) toPlugin() *Plugin {
	return &Plugin{
		ID:          c.id,
		Name:        c.getName(),
		Description: c.Description,
		URL:         c.URL,
		Version:     c.Version,
		Tasks:       c.getPluginTasks(false),
		Hooks:       c.getPluginHooks(false),
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

func (c Config) getHooks(hookType HookTriggerEnum) []*HookConfig {
	var ret []*HookConfig
	for _, h := range c.Hooks {
		for _, t := range h.TriggeredBy {
			if hookType == t {
				ret = append(ret, h)
			}
		}
	}

	return ret
}

func (c Config) getConfigPath() string {
	return filepath.Dir(c.path)
}

func (c Config) getExecCommand(task *OperationConfig) []string {
	ret := c.Exec

	ret = append(ret, task.ExecArgs...)

	if len(ret) > 0 {
		_, err := exec.LookPath(ret[0])
		if err != nil {
			// change command to use absolute path
			pluginPath := filepath.Dir(c.path)
			ret[0] = filepath.Join(pluginPath, ret[0])
		}
	}

	// replace {pluginDir} in arguments with that of the plugin directory
	dir := c.getConfigPath()
	for i, arg := range ret {
		if i == 0 {
			continue
		}

		ret[i] = strings.ReplaceAll(arg, "{pluginDir}", dir)
	}

	return ret
}

type interfaceEnum string

// Valid interfaceEnum values
const (
	// InterfaceEnumRPC indicates that the plugin uses the RPCRunner interface
	// declared in common/rpc.go.
	InterfaceEnumRPC interfaceEnum = "rpc"

	// InterfaceEnumRaw interfaces will have the common.PluginInput encoded as
	// json (but may be ignored), and output will be decoded as
	// common.PluginOutput. If this decoding fails, then the raw output will be
	// treated as the output.
	InterfaceEnumRaw interfaceEnum = "raw"

	InterfaceEnumJS interfaceEnum = "js"
)

func (i interfaceEnum) Valid() bool {
	return i == InterfaceEnumRPC || i == InterfaceEnumRaw || i == InterfaceEnumJS
}

func (i *interfaceEnum) getTaskBuilder() taskBuilder {
	if *i == InterfaceEnumRaw {
		return &rawTaskBuilder{}
	}

	if *i == InterfaceEnumRPC {
		return &rpcTaskBuilder{}
	}

	if *i == InterfaceEnumJS {
		return &jsTaskBuilder{}
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
	ExecArgs []string `yaml:"execArgs"`

	// A map of argument keys to their default values. The default value is
	// used if the applicable argument is not provided during the operation
	// call.
	DefaultArgs map[string]string `yaml:"defaultArgs"`
}

type HookConfig struct {
	OperationConfig `yaml:",inline"`

	// A list of stash operations that will be used to trigger this hook operation.
	TriggeredBy []HookTriggerEnum `yaml:"triggeredBy"`
}

func loadPluginFromYAML(reader io.Reader) (*Config, error) {
	ret := &Config{}

	parser := yaml.NewDecoder(reader)
	parser.SetStrict(true)
	err := parser.Decode(&ret)
	if err != nil {
		return nil, err
	}

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
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ret, err := loadPluginFromYAML(file)
	if err != nil {
		return nil, err
	}

	// set id to the filename
	id := filepath.Base(path)
	ret.id = id[:strings.LastIndex(id, ".")]
	ret.path = path

	return ret, nil
}
