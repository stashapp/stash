package utils

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	_ "golang.org/x/image/webp"
	"image"
	"image/jpeg"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

//ListZipContents returns the images in a zip file using a zip.File slice and ordered by name using natural order
func ListZipContents(path string) ([]*zip.File, *zip.ReadCloser, error) {
	readCloser, err := zip.OpenReader(path)
	if err != nil {
		fmt.Printf("Failed to read zip file %s", path)
		return nil, nil, err
	}

	filteredFiles := make([]*zip.File, 0)
	for _, file := range readCloser.File {
		if file.FileInfo().IsDir() {
			continue
		}
		if !FilenameIsImage(file.Name) {
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
		return NaturalCompare(a.Name, b.Name)
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

//ListDirContents returns the images in a directory path using a zip.File slice and ordered by name using natural order
func ListDirContents(path string) ([]*zip.File, error) {
	images := ListImages(path)
	if images == nil {
		return nil, fmt.Errorf("error getting images from %s", path)
	}
	var filteredFiles []*zip.File
	for _, image := range images {
		el := zip.File{FileHeader: zip.FileHeader{Name: image}}
		filteredFiles = append(filteredFiles, &el)
	}
	sort.Slice(filteredFiles, func(i, j int) bool {
		a := filteredFiles[i]
		b := filteredFiles[j]
		return NaturalCompare(a.Name, b.Name)
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

//ChecksumFromDirPath calculates the xor of the md5s from images in the dir
//If only one image exists the md5 of the image is returned as is
func ChecksumFromDirPath(path string) (string, error) {
	_, err := DirExists(path)
	if err != nil {
		return "", err
	}
	images := ListImages(path)
	if images == nil {
		return "", fmt.Errorf("no images found in %s", path)
	}

	m := make(map[string]struct{}) // a map is used to filter out duplicate MD5s in a dir
	var empty struct{}

	for _, image := range images { // we only want to keep one md5 if duplicates exist
		md5, err := MD5FromFilePath(image) // because  ( a xor a ) == 0
		if err != nil {
			return "", err
		}
		m[md5] = empty
	}

	var md5s []string
	for k := range m {
		md5s = append(md5s, k)
	} // md5s is now a slice of unique MD5s

	return XorMD5Strings(md5s)

}

//IsPathArchive returns true if path seems like a zip archive
func IsPathArchive(path string) bool {
	ext := filepath.Ext(path)
	if strings.ToLower(ext) != ".zip" {
		return false
	}
	return true
}

//contains returns the index of the first occurrenece of string x ( case insensitive ) in name of zip contents, -1 otherwise
func contains(a []*zip.File, x string) int {
	for i, n := range a {
		if strings.Contains(strings.ToLower(n.Name), strings.ToLower(x)) {
			return i
		}
	}
	return -1
}

//reorder reorders a slice so that element with position toFirst gets at the start
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

//ReadZipFile reads a zip file given by path and returns the image located in index position
func ReadZipFile(index int, path string) ([]byte, error) {
	filteredFiles, readCloser, err := ListZipContents(path)
	if err != nil {
		return nil, err
	}
	defer readCloser.Close()

	zipFile := filteredFiles[index]
	zipFileReadCloser, err := zipFile.Open()
	if err != nil {
		fmt.Printf("failed to read file inside zip file %s\n", path)
		return nil, err
	}
	defer zipFileReadCloser.Close()

	return ioutil.ReadAll(zipFileReadCloser)
}

//ReadFile returns the image located in dir with index position
func ReadFile(index int, dir string) ([]byte, error) {
	filteredFiles, err := ListDirContents(dir)
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

//ImageCount returns the number of images in a dir or zip file located in path
func ImageCount(path string) int {
	var images []*zip.File
	if IsPathArchive(path) {
		images, _, _ = ListZipContents(path)
	} else {
		images, _ = ListDirContents(path)
	}
	if images == nil {
		return 0
	}
	return len(images)
}

//GetImage returns the image in a dir or zip file (path) with position index
func GetImage(index int, path string) []byte {
	var data []byte
	if IsPathArchive(path) {
		data, _ = ReadZipFile(index, path)
	} else {
		data, _ = ReadFile(index, path)
	}
	return data
}

//GetThumbnail returns the thumbnail ( or original image ) of an image in a path ( dir or zip file)
func GetThumbnail(index int, width int, path string) []byte {
	var data []byte
	if IsPathArchive(path) {
		data, _ = ReadZipFile(index, path)
	} else {
		data, _ = ReadFile(index, path)
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
