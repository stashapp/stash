//go:generate go run -mod=vendor github.com/99designs/gqlgen
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/stashapp/stash/internal/api"
	"github.com/stashapp/stash/internal/desktop"
	"github.com/stashapp/stash/internal/manager"
	"github.com/stashapp/stash/ui"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	defer recoverPanic()

	_, err := manager.Initialize()
	if err != nil {
		panic(err)
	}

	go func() {
		defer recoverPanic()
		if err := api.Start(); err != nil {
			handleError(err)
		} else {
			manager.GetInstance().Shutdown(0)
		}
	}()

	go handleSignals()
	desktop.Start(manager.GetInstance(), &manager.FaviconProvider{UIBox: ui.UIBox})

	blockForever()
}

func recoverPanic() {
	if p := recover(); p != nil {
		handleError(fmt.Errorf("Panic: %v", p))
	}
}

func handleError(err error) {
	if desktop.IsDesktop() {
		desktop.FatalError(err)
		manager.GetInstance().Shutdown(0)
	} else {
		panic(err)
	}
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
