//go:generate go run -mod=vendor github.com/99designs/gqlgen
package main

import (
	"embed"
	"fmt"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	"github.com/apenwarr/fixconsole"
	"github.com/stashapp/stash/pkg/api"
	"github.com/stashapp/stash/pkg/manager"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

//go:embed ui/v2.5/build
var uiBox embed.FS

//go:embed ui/login
var loginUIBox embed.FS

func init() {
	// On Windows, attach to parent shell
	err := fixconsole.FixConsoleIfNeeded()
	if err != nil {
		fmt.Printf("FixConsoleOutput: %v\n", err)
	}
}

func main() {
	manager.Initialize()
	api.Start(uiBox, loginUIBox)

	blockForever()

	// stop any profiling at exit
	pprof.StopCPUProfile()

	manager.GetInstance().Shutdown(0)
}

func blockForever() {
	// handle signals
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals
}
