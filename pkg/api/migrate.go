package api

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager"
)

type migrateData struct {
	ExistingVersion uint
	MigrateVersion  uint
	BackupPath      string
}

func getMigrateData() migrateData {
	return migrateData{
		ExistingVersion: database.Version(),
		MigrateVersion:  database.AppSchemaVersion(),
		BackupPath:      database.DatabaseBackupPath(),
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

	formBackupPath := r.Form.Get("backuppath")

	// always backup so that we can roll back to the previous version if
	// migration fails
	backupPath := formBackupPath
	if formBackupPath == "" {
		backupPath = database.DatabaseBackupPath()
	}

	// perform database backup
	if err = database.Backup(backupPath); err != nil {
		http.Error(w, fmt.Sprintf("error backing up database: %s", err), 500)
		return
	}

	err = database.RunMigrations()
	if err != nil {
		errStr := fmt.Sprintf("error performing migration: %s", err)

		// roll back to the backed up version
		restoreErr := database.RestoreFromBackup(backupPath)
		if restoreErr != nil {
			errStr = fmt.Sprintf("ERROR: unable to restore database from backup after migration failure: %s\n%s", restoreErr.Error(), errStr)
		} else {
			errStr = "An error occurred migrating the database to the latest schema version. The backup database file was automatically renamed to restore the database.\n" + errStr
		}

		http.Error(w, errStr, 500)
		return
	}

	// perform post-migration operations
	manager.GetInstance().PostMigrate()

	// if no backup path was provided, then delete the created backup
	if formBackupPath == "" {
		err = os.Remove(backupPath)
		if err != nil {
			logger.Warnf("error removing unwanted database backup (%s): %s", backupPath, err.Error())
		}
	}

	http.Redirect(w, r, "/", 301)
}
