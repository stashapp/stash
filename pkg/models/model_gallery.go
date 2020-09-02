package models

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"image"
	"image/jpeg"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/stashapp/stash/pkg/api/urlbuilders"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/utils"
	_ "golang.org/x/image/webp"
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

func (g *Gallery) CountFiles() int {
	filteredFiles, readCloser, err := g.listZipContents()
	if err != nil {
		return 0
	}
	defer readCloser.Close()

	return len(filteredFiles)
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

func (g *Gallery) GetThumbnail(index int, width int) []byte {
	data, _ := g.readZipFile(index)
	srcImage, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return data
	}
	resizedImage := imaging.Resize(srcImage, width, 0, imaging.Box)
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
		logger.Warnf("failed to read zip file %s", g.Path)
		return nil, nil, err
	}

	filteredFiles := make([]*zip.File, 0)
	for _, file := range readCloser.File {
		if file.FileInfo().IsDir() {
			continue
		}
		ext := filepath.Ext(file.Name)
		ext = strings.ToLower(ext)
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" {
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

	cover := contains(filteredFiles, "cover.jpg") // first image with cover.jpg in the name
	if cover >= 0 {                               // will be moved to the start
		reorderedFiles := reorder(filteredFiles, cover)
		if reorderedFiles != nil {
			return reorderedFiles, readCloser, nil
		}
	}

	return filteredFiles, readCloser, nil
}

// return index of first occurrenece of string x ( case insensitive ) in name of zip contents, -1 otherwise
func contains(a []*zip.File, x string) int {
	for i, n := range a {
		if strings.Contains(strings.ToLower(n.Name), strings.ToLower(x)) {
			return i
		}
	}
	return -1
}

// reorder slice so that element with position toFirst gets at the start
func reorder(a []*zip.File, toFirst int) []*zip.File {
	var first *zip.File
	switch {
	case toFirst < 0 || toFirst >= len(a):
		return nil
	case toFirst == 0:
		return a
	default:
		first = a[toFirst]
		copy(a[toFirst:], a[toFirst+1:])     // Shift a[toFirst+1:] left one index removing a[toFirst] element
		a[len(a)-1] = nil                    // Nil now unused element for garbage collection
		a = a[:len(a)-1]                     // Truncate slice
		a = append([]*zip.File{first}, a...) // Push first to the start of the slice
	}
	return a
}

func (g *Gallery) ImageCount() int {
	images, _, _ := g.listZipContents()
	if images == nil {
		return 0
	}
	return len(images)
}
