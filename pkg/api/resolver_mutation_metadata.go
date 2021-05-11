package api

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/stashapp/stash/pkg/database"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager"
	"github.com/stashapp/stash/pkg/manager/config"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/utils"
)

func (r *mutationResolver) MetadataScan(ctx context.Context, input models.ScanMetadataInput) (string, error) {
	jobID := manager.GetInstance().Scan(input)
	return strconv.Itoa(jobID), nil
}

func (r *mutationResolver) MetadataImport(ctx context.Context) (string, error) {
	jobID, err := manager.GetInstance().Import()
	if err != nil {
		return "", err
	}

	return strconv.Itoa(jobID), nil
}

func (r *mutationResolver) ImportObjects(ctx context.Context, input models.ImportObjectsInput) (string, error) {
	t, err := manager.CreateImportTask(config.GetInstance().GetVideoFileNamingAlgorithm(), input)
	if err != nil {
		return "", err
	}

	jobID := manager.GetInstance().RunSingleTask(t)

	return strconv.Itoa(jobID), nil
}

func (r *mutationResolver) MetadataExport(ctx context.Context) (string, error) {
	jobID, err := manager.GetInstance().Export()
	if err != nil {
		return "", err
	}

	return strconv.Itoa(jobID), nil
}

func (r *mutationResolver) ExportObjects(ctx context.Context, input models.ExportObjectsInput) (*string, error) {
	t := manager.CreateExportTask(config.GetInstance().GetVideoFileNamingAlgorithm(), input)

	var wg sync.WaitGroup
	wg.Add(1)
	t.Start(&wg)

	if t.DownloadHash != "" {
		baseURL, _ := ctx.Value(BaseURLCtxKey).(string)

		// generate timestamp
		suffix := time.Now().Format("20060102-150405")
		ret := baseURL + "/downloads/" + t.DownloadHash + "/export" + suffix + ".zip"
		return &ret, nil
	}

	return nil, nil
}

func (r *mutationResolver) MetadataGenerate(ctx context.Context, input models.GenerateMetadataInput) (string, error) {
	jobID := manager.GetInstance().Generate(input)
	return strconv.Itoa(jobID), nil
}

func (r *mutationResolver) MetadataAutoTag(ctx context.Context, input models.AutoTagMetadataInput) (string, error) {
	jobID := manager.GetInstance().AutoTag(input)
	return strconv.Itoa(jobID), nil
}

func (r *mutationResolver) MetadataClean(ctx context.Context, input models.CleanMetadataInput) (string, error) {
	jobID := manager.GetInstance().Clean(input)
	return strconv.Itoa(jobID), nil
}

func (r *mutationResolver) MigrateHashNaming(ctx context.Context) (string, error) {
	jobID := manager.GetInstance().MigrateHash()
	return strconv.Itoa(jobID), nil
}

func (r *mutationResolver) BackupDatabase(ctx context.Context, input models.BackupDatabaseInput) (*string, error) {
	// if download is true, then backup to temporary file and return a link
	download := input.Download != nil && *input.Download
	mgr := manager.GetInstance()
	var backupPath string
	if download {
		utils.EnsureDir(mgr.Paths.Generated.Downloads)
		f, err := ioutil.TempFile(mgr.Paths.Generated.Downloads, "backup*.sqlite")
		if err != nil {
			return nil, err
		}

		backupPath = f.Name()
		f.Close()
	} else {
		backupPath = database.DatabaseBackupPath()
	}

	err := database.Backup(database.DB, backupPath)
	if err != nil {
		return nil, err
	}

	if download {
		downloadHash := mgr.DownloadStore.RegisterFile(backupPath, "", false)
		logger.Debugf("Generated backup file %s with hash %s", backupPath, downloadHash)

		baseURL, _ := ctx.Value(BaseURLCtxKey).(string)

		fn := filepath.Base(database.DatabaseBackupPath())
		ret := baseURL + "/downloads/" + downloadHash + "/" + fn
		return &ret, nil
	} else {
		logger.Infof("Successfully backed up database to: %s", backupPath)
	}

	return nil, nil
}
