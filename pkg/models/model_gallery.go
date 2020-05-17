package models

import (
	"archive/zip"
	"database/sql"
	"github.com/stashapp/stash/pkg/api/urlbuilders"
	"github.com/stashapp/stash/pkg/utils"
)

type Gallery struct {
	ID        int             `db:"id" json:"id"`
	Path      string          `db:"path" json:"path"`
	Checksum  string          `db:"checksum" json:"checksum"`
	SceneID   sql.NullInt64   `db:"scene_id,omitempty" json:"scene_id"`
	CreatedAt SQLiteTimestamp `db:"created_at" json:"created_at"`
	UpdatedAt SQLiteTimestamp `db:"updated_at" json:"updated_at"`
}

const DefaultGthumbWidth int = 200

func (g *Gallery) GetFiles(baseURL string) []*GalleryFilesType {
	var galleryFiles []*GalleryFilesType
	var readCloser *zip.ReadCloser
	var filteredFiles []*zip.File
	var err error

	if utils.IsPathArchive(g.Path) {
		filteredFiles, readCloser, err = utils.ListZipContents(g.Path)
		if err != nil {
			return nil
		}
		defer readCloser.Close()
	} else {
		filteredFiles, err = utils.ListDirContents(g.Path)
		if err != nil {
			return nil
		}

	}

	builder := urlbuilders.NewGalleryURLBuilder(baseURL, g.ID)
	for i, file := range filteredFiles {
		galleryURL := builder.GetGalleryImageURL(i)
		galleryFile := GalleryFilesType{
			Index: i,
			Name:  &file.Name,
			Path:  &galleryURL,
		}
		galleryFiles = append(galleryFiles, &galleryFile)
	}

	return galleryFiles
}
