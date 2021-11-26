package gallery

import (
	"archive/zip"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/manager/paths"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/plugin"
	"github.com/stashapp/stash/pkg/utils"
)

type Scanner struct {
	file.Scanner

	ImageExtensions    []string
	StripFileExtension bool
	Ctx                context.Context
	CaseSensitiveFs    bool
	TxnManager         models.TransactionManager
	Paths              *paths.Paths
	PluginCache        *plugin.Cache

	Gallery    *models.Gallery
	ScanImages bool
}

func FileScanner(hasher file.Hasher, statter file.Statter) file.Scanner {
	return file.Scanner{
		Hasher:       hasher,
		Statter:      statter,
		CalculateMD5: true,
		Done:         make(chan struct{}),
	}
}

func (scanner *Scanner) PostScan(scanned file.Scanned) error {
	if scanned.Old != nil {
		// should be an existing gallery
		var gallery *models.Gallery
		if err := scanner.TxnManager.WithReadTxn(scanner.Ctx, func(r models.ReaderRepository) error {
			galleries, err := r.Gallery().FindByFileID(scanned.Old.ID)
			if err != nil {
				return err
			}

			// assume only one gallery for now
			if len(galleries) > 0 {
				gallery = galleries[0]
			}
			return err
		}); err != nil {
			logger.Error(err.Error())
			return nil
		}

		if gallery != nil {
			return scanner.ScanExisting(gallery, scanned)
		}

		// we shouldn't be able to have an existing file without a gallery, but
		// assuming that it's happened, treat it as a new gallery
	}

	// assume a new file/scene
	return scanner.ScanNew(scanned.New)
}

func (scanner *Scanner) ScanExisting(existing *models.Gallery, scanned file.Scanned) error {
	retGallery := existing

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
		scanner.ScanImages = true

		logger.Infof("%s has been updated: rescanning", path)

		retGallery.UpdatedAt = models.SQLiteTimestamp{Timestamp: time.Now()}

		if err := scanner.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
			if err := scanner.Scanner.ApplyChanges(r.File(), &scanned); err != nil {
				return err
			}

			// ensure no clashes of hashes
			if scanned.New.Checksum != "" && scanned.Old.Checksum != scanned.New.Checksum {
				dupe, _ := r.Gallery().FindByChecksum(retGallery.Checksum)
				if dupe != nil {
					return fmt.Errorf("MD5 for file %s is the same as that of %s", path, dupe.Path.String)
				}
			}

			var err error
			retGallery, err = r.Gallery().Update(*retGallery)
			return err
		}); err != nil {
			return err
		}

		scanner.PluginCache.ExecutePostHooks(scanner.Ctx, retGallery.ID, plugin.GalleryUpdatePost, nil, nil)

		scanner.Gallery = retGallery
	}

	return nil
}

func (scanner *Scanner) ScanNew(f *models.File) error {
	path := f.Path
	checksum := f.Checksum
	isNewGallery := false
	isUpdatedGallery := false
	var g *models.Gallery

	if err := scanner.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
		if err := scanner.Scanner.ApplyChanges(r.File(), &file.Scanned{
			New: f,
		}); err != nil {
			return err
		}

		qb := r.Gallery()

		g, _ = qb.FindByChecksum(checksum)
		if g != nil {
			logger.Infof("%s already exists.  Duplicate of %s ", path, g.Path.String)

			// link gallery to file
			return addFile(r.Gallery(), g, f)
		} else if scanner.hasImages(path) { // don't create gallery if it has no images
			currentTime := time.Now()

			g = &models.Gallery{
				Zip: true,
				Title: sql.NullString{
					String: utils.GetNameFromPath(path, scanner.StripFileExtension),
					Valid:  true,
				},
				CreatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
				UpdatedAt: models.SQLiteTimestamp{Timestamp: currentTime},
			}

			g.SetFile(*f)

			// only warn when creating the gallery
			ok, err := utils.IsZipFileUncompressed(path)
			if err == nil && !ok {
				logger.Warnf("%s is using above store (0) level compression.", path)
			}

			logger.Infof("%s doesn't exist.  Creating new item...", path)
			g, err = qb.Create(*g)
			if err != nil {
				return err
			}

			isNewGallery = true
			scanner.ScanImages = true

			// link scene to file
			return addFile(r.Gallery(), g, f)
		}

		return nil
	}); err != nil {
		return err
	}

	if isNewGallery {
		scanner.PluginCache.ExecutePostHooks(scanner.Ctx, g.ID, plugin.GalleryCreatePost, nil, nil)
		scanner.Gallery = g
	} else if isUpdatedGallery {
		scanner.PluginCache.ExecutePostHooks(scanner.Ctx, g.ID, plugin.GalleryUpdatePost, nil, nil)
		// file was not updated, don't rescan
	}

	return nil
}

func addFile(rw models.GalleryReaderWriter, g *models.Gallery, f *models.File) error {
	ids, err := rw.GetFileIDs(g.ID)
	if err != nil {
		return err
	}

	ids = utils.IntAppendUnique(ids, f.ID)
	return rw.UpdateFiles(g.ID, ids)
}

func (scanner *Scanner) isImage(pathname string) bool {
	return utils.MatchExtension(pathname, scanner.ImageExtensions)
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
