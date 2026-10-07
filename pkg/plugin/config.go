package plugin

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/stashapp/stash/pkg/plugin/hook"
	"github.com/stashapp/stash/pkg/python"
	"github.com/stashapp/stash/pkg/utils"
	"gopkg.in/yaml.v2"
)

// Config describes configuration single plugin.
type Config struct {
	// path configuration file
	path string
	// name plugin. displayed UI.
	Name string `yaml:"name"`

	// optional description what plugin does.
	Description *string `yaml:"description"`

	// optional URL plugin.
	URL *string `yaml:"url"`

	// optional version string.
	Version *string `yaml:"version"`

	// communication interface when communicating spawned
	// plugin process. Defaults 'raw' not provided.
	Interface interfaceEnum `yaml:"interface"`

	// command execute operations plugin. first
	// element program name, subsequent elements passed
	// arguments.
	//
	// Note: execution process search path program,
	// then attempt find program plugins
	// directory. exe extension not necessary Windows platforms.
	// current working directory set stash process.
	Exec []string `yaml:"exec,flow"`

	// default log level output plugin process's stderr stream.
	// Only plugin not encode output using log level
	// control characters.
	// See package common/log valid values.
	// left unset, defaults log.ErrorLevel.
	PluginErrLogLevel string `yaml:"errLog"`

	// task configurations tasks provided plugin.
	Tasks []*OperationConfig `yaml:"tasks"`

	// hooks configurations hooks registered plugin.
	Hooks []*HookConfig `yaml:"hooks"`

	// Javascript files injected stash UI.
	UIConfig `yaml:"ui"`

	// Settings configure plugin.
	Settings map[string]SettingConfig `yaml:"settings"`
}

type PluginCSP struct {
	ScriptSrc  []string `json:"script-src" yaml:"script-src"`
	StyleSrc   []string `json:"style-src" yaml:"style-src"`
	ConnectSrc []string `json:"connect-src" yaml:"connect-src"`
}

type UIConfig struct {
	// Requires list plugin IDs plugin depends on.
	// plugins loaded before plugin.
	Requires []string `yaml:"requires"`

	// Content Security Policy configuration plugin.
	CSP PluginCSP `yaml:"csp"`

	// Javascript files injected stash UI.
	// URLs paths files relative plugin configuration file.
	Javascript []string `yaml:"javascript"`

	// CSS files injected stash UI.
	// URLs paths files relative plugin configuration file.
	CSS []string `yaml:"css"`

	// Assets map URL prefixes hosted directories.
	// allows plugins serve static assets URL path.
	// Plugin assets exposed /plugin/{pluginId}/assets path.
	// example, plugin configuration file contains:
	// /foobar
	// /barbaz
	// /: root
	// Then following requests mapped following files:
	// /plugin/{pluginId}/assets/foo/file.txt {pluginDir}/foo/file.txt
	// /plugin/{pluginId}/assets/bar/file.txt {pluginDir}/baz/file.txt
	// /plugin/{pluginId}/assets/file.txt {pluginDir}/root/file.txt
	Assets utils.URLMap `yaml:"assets"`
}

func isURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

func (c UIConfig) getCSSFiles(parent Config) []string {
	var ret []string
	for _, v := range c.CSS {
		if !isURL(v) {
			ret = append(ret, filepath.Join(parent.getConfigPath(), v))
		}
	}
	return ret
}

func (c UIConfig) getExternalCSS() []string {
	var ret []string
	for _, v := range c.CSS {
		if isURL(v) {
			ret = append(ret, v)
		}
	}
	return ret
}

func (c UIConfig) getJavascriptFiles(parent Config) []string {
	var ret []string
	for _, v := range c.Javascript {
		if !isURL(v) {
			ret = append(ret, filepath.Join(parent.getConfigPath(), v))
		}
	}
	return ret
}

func (c UIConfig) getExternalScripts() []string {
	var ret []string
	for _, v := range c.Javascript {
		if isURL(v) {
			ret = append(ret, v)
		}
	}
	return ret
}

type SettingConfig struct {
	// defaults string
	Type PluginSettingTypeEnum `yaml:"type"`
	// defaults key name
	DisplayName string `yaml:"displayName"`
	Description string `yaml:"description"`
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

func convertHooks(hooks []hook.TriggerEnum) []string {
	var ret []string
	for _, h := range hooks {
		ret = append(ret, h.String())
	}
	return ret
}

func (c Config) getPluginSettings() []PluginSetting {
	ret := []PluginSetting{}
	var keys []string
	for k := range c.Settings {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		o := c.Settings[k]
		t := o.Type
		if t == "" {
			t = PluginSettingTypeEnumString
		}
		ret = append(ret, PluginSetting{
			Name:        k,
			DisplayName: o.DisplayName,
			Description: o.Description,
			Type:        t,
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

func (c Config) toPlugin() *Plugin {
	return &Plugin{
		ID:          c.id,
		Name:        c.getName(),
		Description: c.Description,
		URL:         c.URL,
		Version:     c.Version,
		Tasks:       c.getPluginTasks(false),
		Hooks:       c.getPluginHooks(false),
		UI: PluginUI{
			Requires:       c.UI.Requires,
			ExternalScript: c.UI.getExternalScripts(),
			ExternalCSS:    c.UI.getExternalCSS(),
			Javascript:     c.UI.getJavascriptFiles(c),
			CSS:            c.UI.getCSSFiles(c),
			CSP:            c.UI.CSP,
			Assets:         c.UI.Assets,
			Settings:       c.getPluginSettings(),
			ConfigPath:     c.path,
		},
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

func (c Config) getHooks(hookType hook.TriggerEnum) []*HookConfig {
	var ret []*HookConfig
	for _, h := range c.Hooks {
		for _, t := range h.TriggeredBy {
			if t == hookType {
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
	// #4859 don't modify original exec command
	ret := append([]string{}, c.Exec...)

	if task != nil {
		ret = append(ret, task.ExecArgs...)
	}

	// #4859 don't use plugin path exec command python command
	if len(ret) > 0 && !python.IsPythonCommand(ret[0]) {
		_, err := exec.LookPath(ret[0])
		if err != nil {
			// change command run plugin path
			pluginPath := filepath.Dir(c.path)
			ret[0] = filepath.Join(pluginPath, ret[0])
		}
	}

	// replace {pluginDir} arguments plugin directory
	dir := c.getConfigPath()
	for i, arg := range ret {
		if i == 0 {
			continue
		}
		ret[i] = strings.ReplaceAll(arg, "{pluginDir}", dir)
	}

	return ret
}

func (c Config) valid() error {
	if c.Interface != "" && !c.Interface.Valid() {
		return fmt.Errorf("invalid interface type %s", c.Interface)
	}

	for _, o := range c.Settings {
		if o.Type != "" && !o.Type.IsValid() {
			return fmt.Errorf("invalid type setting %s", o.Type)
		}
	}

	// Validate exec commands to prevent command injection
	if err := c.validateExec(); err != nil {
		return err
	}

	return nil
}

// validateExec validates that exec commands are safe and don't contain shell metacharacters
func (c Config) validateExec() error {
	// List of allowed base commands for plugins
	allowedCommands := map[string]bool{
		"python3":  true,
		"python":   true,
		"node":     true,
		"npm":      true,
		"npx":      true,
		"bash":     true,
		"sh":       true,
		"ffmpeg":   true,
		"ffprobe":  true,
		"magick":   true,
		"convert":  true,
		"identify": true,
		"exiftool": true,
		"jq":       true,
		"yq":       true,
		"sqlite3":  true,
		"git":      true,
		"curl":     true,
		"wget":     true,
	}

	// Check both plugin-level exec and task-level execArgs
	allExecs := make([][]string, 0)
	if len(c.Exec) > 0 {
		allExecs = append(allExecs, c.Exec)
	}
	for _, task := range c.Tasks {
		if len(task.ExecArgs) > 0 {
			allExecs = append(allExecs, task.ExecArgs)
		}
	}

	for _, exec := range allExecs {
		if len(exec) == 0 {
			continue
		}

		// Check for shell metacharacters in any argument
		for _, arg := range exec {
			if strings.ContainsAny(arg, ";|&$()`") {
				return fmt.Errorf("shell metacharacters not allowed in exec: %s", arg)
			}
		}

		// Validate the base command is in allowed list
		baseCmd := filepath.Base(exec[0])
		if !allowedCommands[baseCmd] {
			// Allow commands that are relative paths within the plugin directory
			if !strings.HasPrefix(exec[0], "./") && !strings.HasPrefix(exec[0], "../") {
				return fmt.Errorf("command not allowed: %s (not in allowlist)", exec[0])
			}
		}
	}

	return nil
}

type interfaceEnum string

// Valid interfaceEnum values
const (
	// InterfaceEnumRPC indicates plugin uses RPCRunner interface
	// declared common/rpc.go.
	InterfaceEnumRPC interfaceEnum = "rpc"

	// InterfaceEnumRaw interfaces common.PluginInput encoded
	// json (but may be ignored), output decoded
	// common.PluginOutput. decoding fails, then raw output
	// treated output.
	InterfaceEnumRaw interfaceEnum = "raw"

	InterfaceEnumJS interfaceEnum = "js"
)

func (i interfaceEnum) Valid() bool {
	return i == InterfaceEnumRPC || i == InterfaceEnumRaw || i == InterfaceEnumJS
}

func (i *interfaceEnum) getTaskBuilder() taskBuilder {
	switch *i {
	case InterfaceEnumRaw:
		return &rawTaskBuilder{}
	case InterfaceEnumRPC:
		return &rpcTaskBuilder{}
	case InterfaceEnumJS:
		return &jsTaskBuilder{}
	}
	// shouldn't happen
	return nil
}

// OperationConfig describes configuration single plugin operation
// provided plugin.
type OperationConfig struct {
	// identify operation. Must unique within plugin
	// configuration. name shown button operation
	// UI.
	Name string `yaml:"name"`

	// short description operation. description shown below
	// button UI.
	Description string `yaml:"description"`

	// list arguments appended plugin's Exec arguments
	// when executing operation.
	ExecArgs []string `yaml:"execArgs"`

	// map argument keys default values. default value
	// applicable argument not provided during operation
	// call.
	DefaultArgs map[string]string `yaml:"defaultArgs"`
}

type HookConfig struct {
	OperationConfig `yaml:",inline"`

	// list stash operations trigger hook operation.
	TriggeredBy []hook.TriggerEnum `yaml:"triggeredBy"`
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

	err = ret.valid()
	if err != nil {
		return nil, err
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

	// set filename
	id := filepath.Base(path)
	ret.id = id[:strings.LastIndex(id, ".")]
	ret.path = path
	return ret, nil
}
