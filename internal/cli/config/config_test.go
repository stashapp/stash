package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfigWithDefaults(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")
	dbPath := filepath.Join(dir, "stash.sqlite")

	err := os.WriteFile(configPath, []byte(`database_path = "`+dbPath+`"`), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.DatabasePath != dbPath {
		t.Fatalf("DatabasePath = %q, want %q", cfg.DatabasePath, dbPath)
	}
	if cfg.GraphicsMode != GraphicsAuto {
		t.Fatalf("GraphicsMode = %q, want %q", cfg.GraphicsMode, GraphicsAuto)
	}
	if got := strings.Join(cfg.DisplayFields, ","); got != "name,duration,date" {
		t.Fatalf("DisplayFields = %q", got)
	}
	if cfg.CacheDir == "" {
		t.Fatal("CacheDir should have a default")
	}
	if cfg.LogFile == "" {
		t.Fatal("LogFile should have a default")
	}
	if !strings.HasSuffix(cfg.LogFile, filepath.Join("stash-cli", "stash-cli.log")) {
		t.Fatalf("LogFile = %q, want stash-cli/stash-cli.log suffix", cfg.LogFile)
	}
	if cfg.LogLevel != "info" {
		t.Fatalf("LogLevel = %q, want info", cfg.LogLevel)
	}
	if cfg.LogStdout {
		t.Fatal("LogStdout should default to false")
	}
	if cfg.Blobs.Storage != "database" {
		t.Fatalf("Blobs.Storage = %q, want database", cfg.Blobs.Storage)
	}
	if cfg.FFplayPath == "" {
		t.Fatal("FFplayPath should have a default")
	}
	if got := strings.Join(cfg.FFplayArgs, ","); got != "-autoexit,-hide_banner,-loglevel,warning" {
		t.Fatalf("FFplayArgs = %q, want default ffplay args", got)
	}
}

func TestLoadConfigCanSetFFplayCommand(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.toml")
	err := os.WriteFile(configPath, []byte(`
database_path = "/tmp/stash.sqlite"
ffplay_path = "/usr/local/bin/ffplay"
ffplay_args = ["-autoexit", "-fs"]
`), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.FFplayPath != "/usr/local/bin/ffplay" {
		t.Fatalf("FFplayPath = %q, want configured path", cfg.FFplayPath)
	}
	if got := strings.Join(cfg.FFplayArgs, ","); got != "-autoexit,-fs" {
		t.Fatalf("FFplayArgs = %q, want configured args", got)
	}
}

func TestLoadConfigIgnoresRemovedCoverFetchSettings(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.toml")
	err := os.WriteFile(configPath, []byte(`
database_path = "/tmp/stash.sqlite"
ffmpeg_path = "/usr/bin/ffmpeg"

[cover_fallback]
generate_with_ffmpeg = true

[stash_box]
endpoint = "https://stashdb.org/graphql"
api_key = "secret"
max_requests_per_minute = 240
`), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := Load(configPath); err != nil {
		t.Fatal(err)
	}
}

func TestLoadConfigCanUseFilesystemBlobs(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.toml")
	blobsPath := filepath.Join(dir, "blobs")
	err := os.WriteFile(configPath, []byte(`
database_path = "/tmp/stash.sqlite"

[blobs]
storage = "filesystem"
path = "`+blobsPath+`"
`), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Blobs.Storage != "filesystem" {
		t.Fatalf("Blobs.Storage = %q, want filesystem", cfg.Blobs.Storage)
	}
	if cfg.Blobs.Path != blobsPath {
		t.Fatalf("Blobs.Path = %q, want %q", cfg.Blobs.Path, blobsPath)
	}
}

func TestLoadConfigValidatesFilesystemBlobsPath(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.toml")
	err := os.WriteFile(configPath, []byte(`
database_path = "/tmp/stash.sqlite"

[blobs]
storage = "filesystem"
`), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	_, err = Load(configPath)
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), "blobs.path") {
		t.Fatalf("error = %v, want blobs.path validation", err)
	}
}

func TestLoadConfigValidatesStartupScanDirs(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.toml")
	err := os.WriteFile(configPath, []byte(`
database_path = "/tmp/stash.sqlite"
scan_on_startup = true
`), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	_, err = Load(configPath)
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), "media_dirs") {
		t.Fatalf("error = %v, want media_dirs validation", err)
	}
}

func TestLoadConfigValidatesGraphicsMode(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.toml")
	err := os.WriteFile(configPath, []byte(`
database_path = "/tmp/stash.sqlite"
graphics_mode = "sixel"
`), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	_, err = Load(configPath)
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !strings.Contains(err.Error(), "graphics_mode") {
		t.Fatalf("error = %v, want graphics_mode validation", err)
	}
}

func TestExampleConfigIsLoadable(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(configPath, []byte(Example()), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.ScanOnStartup {
		t.Fatal("example should enable startup scan")
	}
	if len(cfg.MediaDirs) == 0 {
		t.Fatal("example should include media directories")
	}
}
