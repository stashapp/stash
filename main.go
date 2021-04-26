//go:generate go run -mod=vendor github.com/99designs/gqlgen
package main

import (
	"github.com/stashapp/stash/pkg/api"
	"github.com/stashapp/stash/pkg/dlna"
	"github.com/stashapp/stash/pkg/manager"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	manager.Initialize()
	api.Start()
	dlna.Start()
	blockForever()
}

func blockForever() {
	select {}
}
