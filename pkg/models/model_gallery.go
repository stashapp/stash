package models

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"github.com/disintegration/imaging"
	"github.com/stashapp/stash/pkg/api/urlbuilders"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/utils"
	"image"
	"image/jpeg"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

type Gallery struct {
	ID        int             `db:"id" json:"id"`
	Path      string          `db:"path" json:"path"`
	Checksum  string          `db:"checksum" json:"checksum"`
	SceneID   sql.NullInt64   `db:"scene_id,omitempty" json:"scene_id"`
	CreatedAt SQLiteTimestamp `db:"created_at" json:"created_at"`
	UpdatedAt SQLiteTimestamp `db:"updated_at" json:"updated_at"`
}

func (g *Gallery) GetFiles(baseURL string) []*GalleryFilesType {
	var galleryFiles []*GalleryFilesType
	filteredFiles, readCloser, err := g.listZipContents()
	if err != nil {
		return nil
	}
	defer readCloser.Close()

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

func (g *Gallery) GetImage(index int) []byte {
	data, _ := g.readZipFile(index)
	return data
}

func (g *Gallery) GetThumbnail(index int, height int) []byte {
	data, _ := g.readZipFile(index)
	srcImage, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return data
	}
	resizedImage := imaging.Resize(srcImage, height, 0, imaging.Box)
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, resizedImage, nil)
	if err != nil {
		return data
	}
	return buf.Bytes()
}

func (g *Gallery) readZipFile(index int) ([]byte, error) {
	filteredFiles, readCloser, err := g.listZipContents()
	if err != nil {
		return nil, err
	}
	defer readCloser.Close()

	zipFile := filteredFiles[index]
	zipFileReadCloser, err := zipFile.Open()
	if err != nil {
		logger.Warn("failed to read file inside zip file")
		return nil, err
	}
	defer zipFileReadCloser.Close()

	return ioutil.ReadAll(zipFileReadCloser)
}

func (g *Gallery) listZipContents() ([]*zip.File, *zip.ReadCloser, error) {
	readCloser, err := zip.OpenReader(g.Path)
	if err != nil {
		logger.Warn("failed to read zip file")
		return nil, nil, err
	}

	filteredFiles := make([]*zip.File, 0)
	for _, file := range readCloser.File {
		if file.FileInfo().IsDir() {
			continue
		}
		ext := filepath.Ext(file.Name)
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
			continue
		}
		if strings.Contains(file.Name, "__MACOSX") {
			continue
		}
		filteredFiles = append(filteredFiles, file)
	}
	sort.Slice(filteredFiles, func(i, j int) bool {
		a := filteredFiles[i]
		b := filteredFiles[j]
		return utils.NaturalCompare(a.Name, b.Name)
	})

	return filteredFiles, readCloser, nil
}
