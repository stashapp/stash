//go:generate go run -mod=vendor github.com/99designs/gqlgen
package main

import (
	"github.com/stashapp/stash/pkg/api"
	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/manager/config"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	manager.Initialize()
	database.RunBackup(config.GetDatabasePath())
	database.Initialize(config.GetDatabasePath())
	api.Start()
	blockForever()
}

func blockForever() {
	select {}
}
