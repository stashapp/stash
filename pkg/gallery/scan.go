package gallery

import (
	"archive/zip"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/fsutil"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/paths"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/txn"
	"github.com/stashapp/stash/pkg/utils"
)

const mutexType = "gallery"

type FinderCreatorUpdater interface {
	FindByChecksum(ctx context.Context, checksum string) (*models.Gallery, error)
	Create(ctx context.Context, newGallery models.Gallery) (*models.Gallery, error)
	Update(ctx context.Context, updatedGallery models.Gallery) (*models.Gallery, error)
}

type Scanner struct {
	file.Scanner

	ImageExtensions    []string
	StripFileExtension bool
	CaseSensitiveFs    bool
	TxnManager         txn.Manager
	CreatorUpdater     FinderCreatorUpdater
	Paths              *paths.Paths
	PluginCache        *plugin.Cache
	MutexManager       *utils.MutexManager
}

func FileScanner(hasher file.Hasher) file.Scanner {
	return file.Scanner{
		Hasher:       hasher,
		CalculateMD5: true,
	}
}

func (scanner *Scanner) ScanExisting(ctx context.Context, existing file.FileBased, file file.SourceFile) (retGallery *models.Gallery, scanImages bool, err error) {
	scanned, err := scanner.Scanner.ScanExisting(existing, file)
	if err != nil {
		return nil, false, err
	}

	// we don't currently store sizes for gallery files
	// clear the file size so that we don't incorrectly detect a
	// change
	scanned.New.Size = ""

	retGallery = existing.(*models.Gallery)

	path := scanned.New.Path

	changed := false

	if scanned.ContentsChanged() {
		retGallery.SetFile(*scanned.New)
		changed = true
	} else if scanned.FileUpdated() {
		logger.Infof("Updated gallery file %s", path)

		retGallery.SetFile(*scanned.New)
		changed = true
	}

	if changed {
		scanImages = true
		logger.Infof("%s has been updated: rescanning", path)

		retGallery.UpdatedAt = models.SQLiteTimestamp{Timestamp: time.Now()}

		// we are operating on a checksum now, so grab a mutex on the checksum
		done := make(chan struct{})
		scanner.MutexManager.Claim(mutexType, scanned.New.Checksum, done)

		if err := txn.WithTxn(ctx, scanner.TxnManager, func(ctx context.Context) error {
			// free the mutex once transaction is complete
			defer close(done)

			// ensure no clashes of hashes
			if scanned.New.Checksum != "" && scanned.Old.Checksum != scanned.New.Checksum {
				dupe, _ := scanner.CreatorUpdater.FindByChecksum(ctx, retGallery.Checksum)
				if dupe != nil {
					return fmt.Errorf("MD5 for file %s is the same as that of %s", path, dupe.Path.String)
				}
			}

			retGallery, err = scanner.CreatorUpdater.Update(ctx, *retGallery)
			return err
		}); err != nil {
			return nil, false, err
		}

		scanner.PluginCache.ExecutePostHooks(ctx, retGallery.ID, plugin.GalleryUpdatePost, nil, nil)
	}

	return
}

func (scanner *Scanner) ScanNew(ctx context.Context, file file.SourceFile) (retGallery *models.Gallery, scanImages bool, err error) {
	scanned, err := scanner.Scanner.ScanNew(file)
	if err != nil {
		return nil, false, err
	}

	path := file.Path()
	checksum := scanned.Checksum
	isNewGallery := false
	isUpdatedGallery := false
	var g *models.Gallery

	// grab a mutex on the checksum
	done := make(chan struct{})
	scanner.MutexManager.Claim(mutexType, checksum, done)
	defer close(done)

	if err := txn.WithTxn(ctx, scanner.TxnManager, func(ctx context.Context) error {
		qb := scanner.CreatorUpdater

		g, _ = qb.FindByChecksum(ctx, checksum)
		if g != nil {
			exists, _ := fsutil.FileExists(g.Path.String)
			if !scanner.CaseSensitiveFs {
				// #1426 - if file exists but is a case-insensitive match for the
				// original filename, then treat it as a move
				if exists && strings.EqualFold(path, g.Path.String) {
					exists = false
				}
			}

			if exists {
				logger.Infof("%s already exists.  Duplicate of %s ", path, g.Path.String)
			} else {
				logger.Infof("%s already exists.  Updating path...", path)
				g.Path = sql.NullString{
					String: path,
					Valid:  true,
				}
				g, err = qb.Update(ctx, *g)
				if err != nil {
					return err
				}

				isUpdatedGallery = true
			}
		} else if scanner.hasImages(path) { // don't create gallery if it has no images
			currentTime := time.Now()

			g = &models.Gallery{
				Zip: true,
				Title: sql.NullString{
					String: fsutil.GetNameFromPath(path, scanner.StripFileExtension),
					Valid:  true,
				},
				CreatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
				UpdatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
			}

			g.SetFile(*scanned)

			// only warn when creating the gallery
			ok, err := isZipFileUncompressed(path)
			if err == nil && !ok {
				logger.Warnf("%s is using above store (0) level compression.", path)
			}

			logger.Infof("%s doesn't exist.  Creating new item...", path)
			g, err = qb.Create(ctx, *g)
			if err != nil {
				return err
			}

			scanImages = true
			isNewGallery = true
		}

		return nil
	}); err != nil {
		return nil, false, err
	}

	if isNewGallery {
		scanner.PluginCache.ExecutePostHooks(ctx, g.ID, plugin.GalleryCreatePost, nil, nil)
	} else if isUpdatedGallery {
		scanner.PluginCache.ExecutePostHooks(ctx, g.ID, plugin.GalleryUpdatePost, nil, nil)
	}

	// Also scan images if zip file has been moved (ie updated) as the image paths are no longer valid
	scanImages = isNewGallery || isUpdatedGallery
	retGallery = g

	return
}

// IsZipFileUnmcompressed returns true if zip file in path is using 0 compression level
func isZipFileUncompressed(path string) (bool, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		fmt.Printf("Error reading zip file %s: %s\n", path, err)
		return false, err
	} else {
		defer r.Close()
		for _, f := range r.File {
			if f.FileInfo().IsDir() { // skip dirs, they always get store level compression
				continue
			}
			return f.Method == 0, nil // check compression level of first actual  file
		}
	}
	return false, nil
}

func (scanner *Scanner) isImage(pathname string) bool {
	return fsutil.MatchExtension(pathname, scanner.ImageExtensions)
}

func (scanner *Scanner) hasImages(path string) bool {
	readCloser, err := zip.OpenReader(path)
	if err != nil {
		logger.Warnf("Error while walking gallery zip: %v", err)
		return false
	}
	defer readCloser.Close()

	for _, file := range readCloser.File {
		if file.FileInfo().IsDir() {
			continue
		}

		if strings.Contains(file.Name, "__MACOSX") {
			continue
		}

		if !scanner.isImage(file.Name) {
			continue
		}

		return true
	}

	return false
}
