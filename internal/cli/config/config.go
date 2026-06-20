package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

const (
	GraphicsAuto     = "auto"
	GraphicsKitty    = "kitty"
	GraphicsListOnly = "list-only"
)

var (
	defaultDisplayFields = []string{"name", "duration", "date"}
	validGraphicsModes   = []string{GraphicsAuto, GraphicsKitty, GraphicsListOnly}
	validDisplayFields   = []string{"name", "title", "duration", "date", "rating", "organized", "path", "performers", "tags"}
)

type Config struct {
	DatabasePath  string   `toml:"database_path"`
	MediaDirs     []string `toml:"media_dirs"`
	ScanOnStartup bool     `toml:"scan_on_startup"`
	DisplayFields []string `toml:"display_fields"`
	GraphicsMode  string   `toml:"graphics_mode"`
	CacheDir      string   `toml:"cache_dir"`
	LogFile       string   `toml:"log_file"`
	LogLevel      string   `toml:"log_level"`
	LogStdout     bool     `toml:"log_stdout"`
	FFprobePath   string   `toml:"ffprobe_path"`
	FFplayPath    string   `toml:"ffplay_path"`
	FFplayArgs    []string `toml:"ffplay_args"`
	Blobs         Blobs    `toml:"blobs"`
}

type Blobs struct {
	Storage string `toml:"storage"`
	Path    string `toml:"path"`
}

func Default() Config {
	return Config{
		DisplayFields: append([]string(nil), defaultDisplayFields...),
		GraphicsMode:  GraphicsAuto,
		CacheDir:      defaultCacheDir(),
		LogFile:       defaultLogFile(),
		LogLevel:      "info",
		FFprobePath:   lookupPath("ffprobe"),
		FFplayPath:    lookupPathOrName("ffplay"),
		FFplayArgs:    []string{"-autoexit", "-hide_banner", "-loglevel", "warning"},
		Blobs: Blobs{
			Storage: "database",
		},
	}
}

func DefaultPath() string {
	if configDir, err := os.UserConfigDir(); err == nil && configDir != "" {
		return filepath.Join(configDir, "stash-cli", "config.toml")
	}

	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, ".config", "stash-cli", "config.toml")
	}

	return "config.toml"
}

func Load(path string) (Config, error) {
	if path == "" {
		path = DefaultPath()
	}

	cfg := Default()
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config file %q: %w", path, err)
	}

	if err := toml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config file %q: %w", path, err)
	}

	if err := cfg.Normalize(); err != nil {
		return Config{}, err
	}
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c *Config) Normalize() error {
	var err error
	c.DatabasePath, err = expandPath(c.DatabasePath)
	if err != nil {
		return fmt.Errorf("normalize database_path: %w", err)
	}

	c.CacheDir, err = expandPath(c.CacheDir)
	if err != nil {
		return fmt.Errorf("normalize cache_dir: %w", err)
	}

	c.LogFile, err = expandPath(c.LogFile)
	if err != nil {
		return fmt.Errorf("normalize log_file: %w", err)
	}

	for i := range c.MediaDirs {
		c.MediaDirs[i], err = expandPath(c.MediaDirs[i])
		if err != nil {
			return fmt.Errorf("normalize media_dirs[%d]: %w", i, err)
		}
	}

	c.GraphicsMode = strings.ToLower(strings.TrimSpace(c.GraphicsMode))
	if c.GraphicsMode == "" {
		c.GraphicsMode = GraphicsAuto
	}

	if len(c.DisplayFields) == 0 {
		c.DisplayFields = append([]string(nil), defaultDisplayFields...)
	}

	if c.CacheDir == "" {
		c.CacheDir = defaultCacheDir()
	}

	if c.LogFile == "" {
		c.LogFile = defaultLogFile()
	}

	c.LogLevel = strings.ToLower(strings.TrimSpace(c.LogLevel))
	if c.LogLevel == "" {
		c.LogLevel = "info"
	}

	c.Blobs.Storage = strings.ToLower(strings.TrimSpace(c.Blobs.Storage))
	if c.Blobs.Storage == "" {
		c.Blobs.Storage = "database"
	}
	c.Blobs.Path, err = expandPath(c.Blobs.Path)
	if err != nil {
		return fmt.Errorf("normalize blobs.path: %w", err)
	}

	c.FFplayPath = strings.TrimSpace(c.FFplayPath)
	if len(c.FFplayArgs) == 0 {
		c.FFplayArgs = []string{"-autoexit", "-hide_banner", "-loglevel", "warning"}
	}

	return nil
}

func (c Config) Validate() error {
	var errs []error

	if strings.TrimSpace(c.DatabasePath) == "" {
		errs = append(errs, errors.New("database_path cannot be empty"))
	}

	if !slices.Contains(validGraphicsModes, c.GraphicsMode) {
		errs = append(errs, fmt.Errorf("graphics_mode must be one of %s", strings.Join(validGraphicsModes, ", ")))
	}

	for _, field := range c.DisplayFields {
		if !slices.Contains(validDisplayFields, field) {
			errs = append(errs, fmt.Errorf("display_fields contains unknown field %q", field))
		}
	}

	if c.ScanOnStartup && len(c.MediaDirs) == 0 {
		errs = append(errs, errors.New("media_dirs cannot be empty when scan_on_startup=true"))
	}

	for _, dir := range c.MediaDirs {
		if strings.TrimSpace(dir) == "" {
			errs = append(errs, errors.New("media_dirs cannot contain empty paths"))
		}
	}

	switch c.Blobs.Storage {
	case "database":
	case "filesystem":
		if strings.TrimSpace(c.Blobs.Path) == "" {
			errs = append(errs, errors.New("blobs.path cannot be empty when blobs.storage=filesystem"))
		}
	default:
		errs = append(errs, errors.New("blobs.storage must be one of database, filesystem"))
	}

	return errors.Join(errs...)
}

func Example() string {
	return `# Point this at an existing Stash server sqlite database. stash-cli will not create it.
database_path = '~/.stash/stash-go.sqlite'
media_dirs = ['/mnt/media/videos', '/run/media/tan/remote/videos']
scan_on_startup = true
display_fields = ['name', 'duration', 'date']
graphics_mode = 'auto'
cache_dir = '~/.cache/stash-cli'
log_file = '~/.local/state/stash-cli/stash-cli.log'
log_level = 'info'
log_stdout = false
ffprobe_path = 'ffprobe'
ffplay_path = 'ffplay'
ffplay_args = ['-autoexit', '-hide_banner', '-loglevel', 'warning']

[blobs]
# Match the Stash server config. Use filesystem when blobs_storage: FILESYSTEM.
storage = 'database'
path = ''
`
}

func expandPath(path string) (string, error) {
	path = strings.TrimSpace(os.ExpandEnv(path))
	if path == "" {
		return "", nil
	}

	if path == "~" || strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		if path == "~" {
			path = home
		} else {
			path = filepath.Join(home, path[2:])
		}
	}

	return filepath.Abs(path)
}

func defaultCacheDir() string {
	if cacheDir, err := os.UserCacheDir(); err == nil && cacheDir != "" {
		return filepath.Join(cacheDir, "stash-cli")
	}

	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, ".cache", "stash-cli")
	}

	return ".stash-cli-cache"
}

func defaultLogFile() string {
	if stateDir := os.Getenv("XDG_STATE_HOME"); stateDir != "" {
		return filepath.Join(stateDir, "stash-cli", "stash-cli.log")
	}

	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, ".local", "state", "stash-cli", "stash-cli.log")
	}

	return "stash-cli.log"
}

func lookupPath(name string) string {
	path, _ := exec.LookPath(name)
	return path
}

func lookupPathOrName(name string) string {
	if path := lookupPath(name); path != "" {
		return path
	}
	return name
}
