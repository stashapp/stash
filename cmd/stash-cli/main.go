package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"

	"github.com/stashapp/stash/internal/cli/browse"
	"github.com/stashapp/stash/internal/cli/config"
	"github.com/stashapp/stash/internal/cli/cover"
	"github.com/stashapp/stash/internal/cli/database"
	"github.com/stashapp/stash/internal/cli/edit"
	"github.com/stashapp/stash/internal/cli/player"
	"github.com/stashapp/stash/internal/cli/scanner"
	"github.com/stashapp/stash/internal/cli/tui"
	"github.com/stashapp/stash/internal/log"
	"github.com/stashapp/stash/pkg/logger"
)

func main() {
	initLog("", false, "error")

	configPath := pflag.StringP("config", "c", "", "path to config.toml")
	printConfig := pflag.Bool("print-config", false, "print an example config.toml")
	check := pflag.Bool("check", false, "validate config and database, then exit")
	help := pflag.BoolP("help", "h", false, "show this help text and exit")
	pflag.Parse()

	if *help {
		pflag.Usage()
		return
	}

	if *printConfig {
		fmt.Print(config.Example())
		return
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		exitError(err)
		return
	}
	initLog(cfg.LogFile, cfg.LogStdout, cfg.LogLevel)

	store, err := database.Open(cfg.DatabasePath, database.Options{
		BlobStorage: database.BlobStorage(cfg.Blobs.Storage),
		BlobsPath:   cfg.Blobs.Path,
	})
	if err != nil {
		exitError(err)
		return
	}
	defer store.Close()

	if cfg.ScanOnStartup {
		s, err := scanner.New(cfg, store.Repo)
		if err != nil {
			exitError(err)
			return
		}

		result := s.Scan(context.Background())
		for _, err := range result.Errors {
			fmt.Fprintf(os.Stderr, "scan error: %v\n", err)
		}
		fmt.Printf("scan: %d files, %d directories, %d errors\n", result.FilesScanned, result.Directories, len(result.Errors))
	}

	if *check {
		fmt.Printf("config ok\n")
		fmt.Printf("database: %s (schema %d)\n", cfg.DatabasePath, store.Schema)
		return
	}

	mode := tui.ViewGrid
	if cfg.GraphicsMode == config.GraphicsListOnly {
		mode = tui.ViewList
	}

	s, err := scanner.New(cfg, store.Repo)
	if err != nil {
		exitError(err)
		return
	}

	if err := tui.Run(context.Background(), tui.Deps{
		Browser:    browse.New(store.Repo),
		Editor:     edit.New(store.Repo),
		Covers:     cover.New(store.Repo, cfg.CacheDir),
		Player:     player.New(store.Repo, cfg.FFplayPath, cfg.FFplayArgs),
		Scanner:    s,
		ForceKitty: cfg.GraphicsMode == config.GraphicsKitty,
	}, mode); err != nil {
		exitError(err)
		return
	}
}

func initLog(logFile string, logStdout bool, logLevel string) {
	if logFile != "" {
		if err := os.MkdirAll(filepath.Dir(logFile), 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "unable to create log directory for %s: %v\n", logFile, err)
		}
	}

	l := log.NewLogger()
	l.Init(logFile, logStdout, logLevel, 0)
	logger.Logger = l
}

func exitError(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
