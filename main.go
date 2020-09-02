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

	// perform the post-migration for new databases
	if database.Initialize(config.GetDatabasePath()) {
		manager.GetInstance().PostMigrate()
	}

	api.Start()
	blockForever()
}

func blockForever() {
	select {}
}
