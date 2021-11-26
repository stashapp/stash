//go:generate go run -mod=vendor github.com/99designs/gqlgen
package main

import (
	"embed"
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stashapp/stash/cmd"
)

//go:embed ui/v2.5/build
var uiBox embed.FS

//go:embed ui/login
var loginUIBox embed.FS

func main() {
	err := cmd.Execute(uiBox, loginUIBox)
	if err != nil {
		log.Fatalf("Execution failure: %v", err)
	}
}
