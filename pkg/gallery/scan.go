package gallery

import (
	"archive/zip"
	"context"
	"database/sql"
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
}

func FileScanner(hasher file.Hasher) file.Scanner {
	return file.Scanner{
		Hasher:       hasher,
		CalculateMD5: true,
	}
}

func (scanner *Scanner) ScanExisting(existing file.FileBased, file file.SourceFile) (retGallery *models.Gallery, scanImages bool, err error) {
	scanned, err := scanner.Scanner.ScanExisting(existing, file)
	if err != nil {
		return nil, false, err
	}

	retGallery = existing.(*models.Gallery)

	path := scanned.New.Path

	changed := false

	if scanned.ContentsChanged() {
		logger.Infof("%s has been updated: rescanning", path)

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

		if err := scanner.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {

			// TODO - ensure no clashes of hashes

			retGallery, err = r.Gallery().Update(*retGallery)
			return err
		}); err != nil {
			return nil, false, err
		}

		scanner.PluginCache.ExecutePostHooks(scanner.Ctx, retGallery.ID, plugin.GalleryUpdatePost, nil, nil)
	}

	return
}

func (scanner *Scanner) ScanNew(file file.SourceFile) (retGallery *models.Gallery, scanImages bool, err error) {
	scanned, err := scanner.Scanner.ScanNew(file)
	if err != nil {
		return nil, false, err
	}

	path := file.Path()
	checksum := scanned.Checksum
	isNewGallery := false
	isUpdatedGallery := false
	var g *models.Gallery

	if err := scanner.TxnManager.WithTxn(context.TODO(), func(r models.Repository) error {
		qb := r.Gallery()

		g, _ = qb.FindByChecksum(checksum)
		if g != nil {
			exists, _ := utils.FileExists(g.Path.String)
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
				g, err = qb.Update(*g)
				if err != nil {
					return err
				}

				isUpdatedGallery = true
			}
		} else {
			// don't create gallery if it has no images
			if scanner.hasImages(path) {
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

				g.SetFile(*scanned)

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

				scanImages = true
				isNewGallery = true
			}
		}

		return nil
	}); err != nil {
		return nil, false, err
	}

	if isNewGallery {
		scanner.PluginCache.ExecutePostHooks(scanner.Ctx, g.ID, plugin.GalleryCreatePost, nil, nil)
	} else if isUpdatedGallery {
		scanner.PluginCache.ExecutePostHooks(scanner.Ctx, g.ID, plugin.GalleryUpdatePost, nil, nil)
	}

	scanImages = isNewGallery
	retGallery = g

	return
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
