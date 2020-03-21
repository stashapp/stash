package api

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/stashapp/stash/pkg/database"
)

type migrateData struct {
	ExistingVersion uint
	MigrateVersion uint
	BackupPath string
}

func getMigrateData() migrateData {
	return migrateData{
		ExistingVersion: database.Version(),
		MigrateVersion: database.AppSchemaVersion(),
		BackupPath: "",
	}
}

func getMigrateHandler(w http.ResponseWriter, r *http.Request) {
	if !database.NeedsMigration() {
		http.Redirect(w, r, "/", 301)
		return
	}

	data, _ := setupUIBox.Find("migrate.html")
	templ, err := template.New("Migrate").Parse(string(data))
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %s", err), 500)
		return
	}

	err = templ.Execute(w, getMigrateData())
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %s", err), 500)
	}
}

func doMigrateHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %s", err), 500)
	}

	r.Form.Get("backuppath")

	// TODO - perform database backup

	err = database.RunMigrations()
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %s", err), 500)
	}
	http.Redirect(w, r, "/", 301)
}
