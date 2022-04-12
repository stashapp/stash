//go:generate go run -mod=vendor github.com/99designs/gqlgen
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/apenwarr/fixconsole"
	"github.com/stashapp/stash/internal/api"
	"github.com/stashapp/stash/internal/desktop"
	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/ui"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func init() {
	// On Windows, attach to parent shell
	err := fixconsole.FixConsoleIfNeeded()
	if err != nil {
		fmt.Printf("FixConsoleOutput: %v\n", err)
	}
}

func main() {
	manager.Initialize()
	api.Start()

	go handleSignals()
	desktop.Start(manager.GetInstance(), &manager.FaviconProvider{UIBox: ui.UIBox})

	blockForever()
}

func handleSignals() {
	// handle signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals
	manager.GetInstance().Shutdown(0)
}

func blockForever() {
	select {}
}
