//go:generate go run -mod=vendor github.com/99designs/gqlgen
package main

import (
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	"github.com/stashapp/stash/pkg/api"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	manager.Initialize()
	api.Start()

	// stop any profiling at exit
	defer pprof.StopCPUProfile()
	blockForever()

	err := manager.GetInstance().Shutdown()
	if err != nil {
		logger.Error("Error when closing: %v", err)
	}
}

func blockForever() {
	// handle signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals
}
