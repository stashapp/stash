package manager

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
)

type databaseBackupZip struct {
	*zip.Writer
}

func (z *databaseBackupZip) zipFileRename(fn, outDir, outFn string) error {
	p := filepath.Join(outDir, outFn)
	p = filepath.ToSlash(p)

	f, err := z.Create(p)
	if err != nil {
		return fmt.Errorf("error creating zip entry for %s: %v", fn, err)
	}

	i, err := os.Open(fn)
	if err != nil {
		return fmt.Errorf("error opening %s: %v", fn, err)
	}

	defer i.Close()

	if _, err := io.Copy(f, i); err != nil {
		return fmt.Errorf("error writing %s to zip: %v", fn, err)
	}

	return nil
}

func (z *databaseBackupZip) zipFile(fn, outDir string) error {
	return z.zipFileRename(fn, outDir, filepath.Base(fn))
}

func (s *Manager) BackupDatabase(download bool, includeBlobs bool) (string, string, error) {
	var backupPath string
	var backupName string

	// if we include blobs, then the output is a zip file
	// if not, using the same backup logic as before, which creates a sqlite file
	if !includeBlobs || s.Config.GetBlobsStorage() != config.BlobStorageTypeFilesystem {
		return s.backupDatabaseOnly(download)
	}

	// use tmp directory for the backup
	backupDir := s.Paths.Generated.Tmp
	if err := fsutil.EnsureDir(backupDir); err != nil {
		return "", "", fmt.Errorf("could not create backup directory %v: %w", backupDir, err)
	}
	f, err := os.CreateTemp(backupDir, "backup*.sqlite")
	if err != nil {
		return "", "", err
	}

	backupPath = f.Name()
	backupName = s.Database.DatabaseBackupPath("")
	f.Close()

	// delete the temp file so that the backup operation can create it
	if err := os.Remove(backupPath); err != nil {
		return "", "", fmt.Errorf("could not remove temporary backup file %v: %w", backupPath, err)
	}

	if err := s.Database.Backup(backupPath); err != nil {
		return "", "", err
	}

	// create a zip file
	zipFileDir := s.Paths.Generated.Downloads
	if !download {
		zipFileDir = s.Config.GetBackupDirectoryPathOrDefault()
		if zipFileDir != "" {
			if err := fsutil.EnsureDir(zipFileDir); err != nil {
				return "", "", fmt.Errorf("could not create backup directory %v: %w", zipFileDir, err)
			}
		}
	}

	zipFileName := backupName + ".zip"
	zipFilePath := filepath.Join(zipFileDir, zipFileName)

	logger.Debugf("Preparing zip file for database backup at %v", zipFilePath)

	zf, err := os.Create(zipFilePath)
	if err != nil {
		return "", "", fmt.Errorf("could not create zip file %v: %w", zipFilePath, err)
	}
	defer zf.Close()

	z := databaseBackupZip{
		Writer: zip.NewWriter(zf),
	}

	defer z.Close()

	// move the database file into the zip
	dbFn := filepath.Base(s.Config.GetDatabasePath())
	if err := z.zipFileRename(backupPath, "", dbFn); err != nil {
		return "", "", fmt.Errorf("could not add database backup to zip file: %w", err)
	}

	if err := os.Remove(backupPath); err != nil {
		return "", "", fmt.Errorf("could not remove temporary backup file %v: %w", backupPath, err)
	}

	// walk the blobs directory and add files to the zip
	blobsDir := s.Config.GetBlobsPath()
	err = filepath.WalkDir(blobsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// calculate out dir by removing the blobsDir prefix from the path
		outDir := filepath.Join("blobs", strings.TrimPrefix(filepath.Dir(path), blobsDir))
		if err := z.zipFile(path, outDir); err != nil {
			return fmt.Errorf("could not add blob %v to zip file: %w", path, err)
		}

		return nil
	})

	if err != nil {
		return "", "", fmt.Errorf("error walking blobs directory: %w", err)
	}

	return zipFilePath, zipFileName, nil
}

func (s *Manager) backupDatabaseOnly(download bool) (string, string, error) {
	var backupPath string
	var backupName string

	if download {
		backupDir := s.Paths.Generated.Downloads
		if err := fsutil.EnsureDir(backupDir); err != nil {
			return "", "", fmt.Errorf("could not create backup directory %v: %w", backupDir, err)
		}
		f, err := os.CreateTemp(backupDir, "backup*.sqlite")
		if err != nil {
			return "", "", err
		}

		backupPath = f.Name()
		backupName = s.Database.DatabaseBackupPath("")
		f.Close()

		// delete the temp file so that the backup operation can create it
		if err := os.Remove(backupPath); err != nil {
			return "", "", fmt.Errorf("could not remove temporary backup file %v: %w", backupPath, err)
		}
	} else {
		backupDir := s.Config.GetBackupDirectoryPathOrDefault()
		if backupDir != "" {
			if err := fsutil.EnsureDir(backupDir); err != nil {
				return "", "", fmt.Errorf("could not create backup directory %v: %w", backupDir, err)
			}
		}
		backupPath = s.Database.DatabaseBackupPath(backupDir)
		backupName = filepath.Base(backupPath)
	}

	err := s.Database.Backup(backupPath)
	if err != nil {
		return "", "", err
	}

	return backupPath, backupName, nil
}
