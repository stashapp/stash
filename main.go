package main

import (
	"github.com/stashapp/stash/internal/api"
	"github.com/stashapp/stash/internal/database"
	"github.com/stashapp/stash/internal/manager"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

//var migrationsBox *packr.Box

func main() {
	//migrationsBox := packr.New("My Box", "./internal/database/migrations")
	//html, err := migrationsBox.FindString("1_initial.up.sql")
	//fmt.Println(html, err)
	//fmt.Println("hello world")

	managerInstance := manager.Initialize()
	database.Initialize(managerInstance.Paths.FixedPaths.DatabaseFile)

	api.Start()
}
