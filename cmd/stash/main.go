//go:generate go run github.com/99designs/gqlgen
package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"runtime/pprof"
	"syscall"

	"github.com/spf13/pflag"

	"github.com/stashapp/stash/internal/api"
	"github.com/stashapp/stash/internal/build"
	"github.com/stashapp/stash/internal/desktop"
	"github.com/stashapp/stash/internal/log"
	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/ui"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var exitCode = 0

func main() {
	defer func() {
		if exitCode != 0 {
			os.Exit(exitCode)
		}
	}()

	defer recoverPanic()

	initLogTemp()

	helpFlag := false
	pflag.BoolVarP(&helpFlag, "help", "h", false, "show this help text and exit")

	versionFlag := false
	pflag.BoolVarP(&versionFlag, "version", "v", false, "show version number and exit")

	cpuProfilePath := ""
	pflag.StringVar(&cpuProfilePath, "cpuprofile", "", "write cpu profile to file")

	pflag.Parse()

	if helpFlag {
		pflag.Usage()
		return
	}

	if versionFlag {
		fmt.Println(build.VersionString())
		return
	}

	cfg, err := config.Initialize()
	if err != nil {
		exitError(fmt.Errorf("config initialization error: %w", err))
		return
	}

	l := initLog(cfg)

	if cpuProfilePath != "" {
		if err := initProfiling(cpuProfilePath); err != nil {
			exitError(err)
			return
		}
		defer pprof.StopCPUProfile()
	}

	// initialise desktop.IsDesktop here so that it doesn't get affected by
	// ffmpeg hardware checks later on
	desktop.InitIsDesktop()

	mgr, err := manager.Initialize(cfg, l)
	if err != nil {
		exitError(fmt.Errorf("manager initialization error: %w", err))
		return
	}
	defer mgr.Shutdown()

	server, err := api.Initialize()
	if err != nil {
		exitError(fmt.Errorf("api initialization error: %w", err))
		return
	}
	defer server.Shutdown()

	exit := make(chan int)

	go func() {
		err := server.Start()
		if !errors.Is(err, http.ErrServerClosed) {
			exitError(fmt.Errorf("http server error: %w", err))
			exit <- 1
		}
	}()

	go handleSignals(exit)
	desktop.Start(exit, &ui.FaviconProvider)

	exitCode = <-exit
}

// initLogTemp initializes a temporary logger for use before the config is loaded.
// Logs only error level message to stderr.
func initLogTemp() *log.Logger {
	l := log.NewLogger()
	l.Init("", true, "Error", 0)
	logger.Logger = l

	return l
}

func initLog(cfg *config.Config) *log.Logger {
	l := log.NewLogger()
	l.Init(cfg.GetLogFile(), cfg.GetLogOut(), cfg.GetLogLevel(), cfg.GetLogFileMaxSize())
	logger.Logger = l

	return l
}

func initProfiling(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("unable to create CPU profile file: %v", err)
	}

	if err = pprof.StartCPUProfile(f); err != nil {
		return fmt.Errorf("could not start CPU profiling: %v", err)
	}

	logger.Infof("profiling to %s", path)

	return nil
}

func recoverPanic() {
	if err := recover(); err != nil {
		exitCode = 1
		logger.Errorf("panic: %v\n%s", err, debug.Stack())
		if desktop.IsDesktop() {
			desktop.FatalError(fmt.Errorf("Panic: %v", err))
		}
	}
}

func exitError(err error) {
	exitCode = 1
	logger.Error(err)
	// #5784 - log to stdout as well as the logger
	// this does mean that it will log twice if the logger is set to stdout
	fmt.Println(err)
	if desktop.IsDesktop() {
		desktop.FatalError(err)
	}
}

func handleSignals(exit chan<- int) {
	// handle signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals
	exit <- 0
}
