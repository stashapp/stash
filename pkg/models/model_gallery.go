package models

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/stashapp/stash/pkg/api/urlbuilders"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/utils"
	_ "golang.org/x/image/webp"
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

const DefaultGthumbWidth int = 200

func (g *Gallery) GetFiles(baseURL string) []*GalleryFilesType {
	var galleryFiles []*GalleryFilesType
	var readCloser *zip.ReadCloser
	var filteredFiles []*zip.File
	var err error

	if IsPathArchive(g.Path) {
		filteredFiles, readCloser, err = g.listZipContents()
		if err != nil {
			return nil
		}
		defer readCloser.Close()
	} else {
		filteredFiles, err = g.listDirContents()
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

func (g *Gallery) GetImage(index int) []byte {
	var data []byte
	if IsPathArchive(g.Path) {
		data, _ = g.readZipFile(index)
	} else {
		data, _ = g.readFile(index)
	}
	return data
}

func (g *Gallery) GetThumbnail(index int, width int) []byte {
	var data []byte
	if IsPathArchive(g.Path) {
		data, _ = g.readZipFile(index)
	} else {
		data, _ = g.readFile(index)
	}

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

func (g *Gallery) readFile(index int) ([]byte, error) {
	filteredFiles, err := g.listDirContents()
	if err != nil {
		return nil, err
	}
	path := filteredFiles[index].Name
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read Error for file %s : %s", path, err)
	}
	return data, nil
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

func (g *Gallery) listDirContents() ([]*zip.File, error) {
	images := utils.ListImages(g.Path)
	if images == nil {
		return nil, fmt.Errorf("error getting images from %s", g.Path)
	}
	var filteredFiles []*zip.File
	for _, image := range images {
		el := zip.File{FileHeader: zip.FileHeader{Name: image}}
		filteredFiles = append(filteredFiles, &el)
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
			return reorderedFiles, nil
		}
	}

	return filteredFiles, nil
}

// calculates the xor of the md5s from images in the dir
// if only one image exists the md5 is returned as is
func ChecksumFromDirPath(path string) (string, error) {
	_, err := utils.DirExists(path)
	if err != nil {
		return "", err
	}
	images := utils.ListImages(path)
	if images == nil {
		return "", fmt.Errorf("no images found in %s", path)
	}

	m := make(map[string]struct{}) // a map is used to filter out duplicate MD5s in a dir
	var empty struct{}

	for _, image := range images { // we only want to keep one md5 if duplicates exist
		md5, err := utils.MD5FromFilePath(image) // because  a xor a = 0
		if err != nil {
			return "", err
		}
		m[md5] = empty
	}

	var md5s []string
	for k := range m {
		md5s = append(md5s, k)
	} // md5s is now a slice of unique MD5s

	return utils.XorMD5Strings(md5s)

}

func IsPathArchive(path string) bool {
	ext := filepath.Ext(path)
	if strings.ToLower(ext) != ".zip" {
		return false
	}
	return true
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
	var images []*zip.File
	if IsPathArchive(g.Path) {
		images, _, _ = g.listZipContents()
	} else {
		images, _ = g.listDirContents()
	}
	if images == nil {
		return 0
	}
	return len(images)
}
