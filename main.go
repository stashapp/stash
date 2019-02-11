package main

import (
	"github.com/stashapp/stash/api"
	"github.com/stashapp/stash/database"
	"github.com/stashapp/stash/manager"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	managerInstance := manager.Initialize()
	database.Initialize(managerInstance.StaticPaths.DatabaseFile)

	api.Start()
}
