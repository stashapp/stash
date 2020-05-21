package utils

import (
	"archive/zip"
	"fmt"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
	"io/ioutil"
	"math"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

var supportedImages = []string{".jpg", ".jpeg", ".png", ".webp", ".gif"}

// FileType uses the filetype package to determine the given file path's type
func FileType(filePath string) (types.Type, error) {
	file, _ := os.Open(filePath)

	// We only have to pass the file header = first 261 bytes
	head := make([]byte, 261)
	_, _ = file.Read(head)

	return filetype.Match(head)
}

// FileExists returns true if the given path exists
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, err
	} else {
		panic(err)
	}
}

// DirExists returns true if the given path exists and is a directory
func DirExists(path string) (bool, error) {
	exists, _ := FileExists(path)
	fileInfo, _ := os.Stat(path)
	if !exists || !fileInfo.IsDir() {
		return false, fmt.Errorf("path either doesn't exist, or is not a directory <%s>", path)
	}
	return true, nil
}

// Touch creates an empty file at the given path if it doesn't already exist
func Touch(path string) error {
	var _, err = os.Stat(path)
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	return nil
}

// EnsureDir will create a directory at the given path if it doesn't already exist
func EnsureDir(path string) error {
	exists, err := FileExists(path)
	if !exists {
		err = os.Mkdir(path, 0755)
		return err
	}
	return err
}

// EnsureDirAll will create a directory at the given path along with any necessary parents if they don't already exist
func EnsureDirAll(path string) error {
	return os.MkdirAll(path, 0755)
}

// RemoveDir removes the given dir (if it exists) along with all of its contents
func RemoveDir(path string) error {
	return os.RemoveAll(path)
}

// EmptyDir will recursively remove the contents of a directory at the given path
func EmptyDir(path string) error {
	d, err := os.Open(path)
	if err != nil {
		return err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range names {
		err = os.RemoveAll(filepath.Join(path, name))
		if err != nil {
			return err
		}
	}

	return nil
}

// ListDir returns the contents of a given directory path as a string slice
func ListDir(path string) []string {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		path = filepath.Dir(path)
		files, err = ioutil.ReadDir(path)
	}

	var dirPaths []string
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		abs, err := filepath.Abs(path)
		if err != nil {
			continue
		}
		dirPaths = append(dirPaths, filepath.Join(abs, file.Name()))
	}
	return dirPaths
}

// GetHomeDirectory returns the path of the user's home directory.  ~ on Unix and C:\Users\UserName on Windows
func GetHomeDirectory() string {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	return currentUser.HomeDir
}

// FilenameIsImage returns true if filename extension is included in our supported images
func FilenameIsImage(path string) bool {
	ext := filepath.Ext(path)
	ext = strings.ToLower(ext)
	for _, image := range supportedImages {
		if strings.ToLower(ext) == strings.ToLower(image) {
			return true
		}
	}
	return false
}

// ListImages returns a slice of image paths present in the directory
// not recursive
func ListImages(path string) []string {
	if path == "" {
		return nil
	}

	absolutePath, err := filepath.Abs(path)
	if err == nil {
		path = absolutePath
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		path = filepath.Dir(path)
		files, err = ioutil.ReadDir(path)
	}

	if err != nil {
		return nil
	}

	var images []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		abs, err := filepath.Abs(path)
		if err != nil {
			continue
		}
		if FilenameIsImage(file.Name()) {
			images = append(images, filepath.Join(abs, file.Name()))
		}
	}

	return images
}

// GetImageDirs gets all directories in the path that contain at least minImg images ( recursive )
// imageDirs slice gets updated during the recursion
func GetImageDirs(path string, info os.FileInfo, imageDirs *[]string, minImg int) {
	var err error
	path, err = filepath.Abs(path)

	if err != nil || !info.IsDir() {
		return
	}

	dir, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error openig path %s: %s\n", path, err)
		return
	}
	defer dir.Close()

	fis, err := dir.Readdir(-1)
	if err != nil {
		return
	}
	imgCount := len(ListImages(path))
	if imgCount > 0 && imgCount >= minImg { // add dir to list only if it has >= minImg images (minImg > 0)
		*imageDirs = append(*imageDirs, path)
	}
	for _, fi := range fis {
		if fi.Name() == "." || fi.Name() == ".." {
			continue
		}
		GetImageDirs(filepath.Join(path, fi.Name()), fi, imageDirs, minImg)
	}
}

// IsZipFileUnmcompressed returns true if zip file in path is using 0 compression level
func IsZipFileUncompressed(path string) (bool, error) {
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

// humanize code taken from https://github.com/dustin/go-humanize and adjusted

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

// HumanizeBytes returns a human readable bytes string of a uint
func HumanizeBytes(s uint64) string {
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	if s < 10 {
		return fmt.Sprintf("%d B", s)
	}
	e := math.Floor(logn(float64(s), 1024))
	suffix := sizes[int(e)]
	val := math.Floor(float64(s)/math.Pow(1024, e)*10+0.5) / 10
	f := "%.0f %s"
	if val < 10 {
		f = "%.1f %s"
	}

	return fmt.Sprintf(f, val, suffix)
}

// WriteFile writes file to path creating parent directories if needed
func WriteFile(path string, file []byte) error {
	pathErr := EnsureDirAll(filepath.Dir(path))
	if pathErr != nil {
		return fmt.Errorf("Cannot ensure path %s", pathErr)
	}

	err := ioutil.WriteFile(path, file, 0755)
	if err != nil {
		return fmt.Errorf("Write error for thumbnail %s: %s ", path, err)
	}
	return nil
}

// GetIntraDir returns a string that can be added to filepath.Join to implement directory depth, "" on error
//eg given a pattern of 0af63ce3c99162e9df23a997f62621c5 and a depth of 2 length of 3
//returns 0af/63c or 0af\63c ( dependin on os)  that can be later used like this  filepath.Join(directory, intradir, basename)
func GetIntraDir(pattern string, depth, length int) string {
	if depth < 1 || length < 1 || (depth*length > len(pattern)) {
		return ""
	}
	intraDir := pattern[0:length] // depth 1 , get length number of characters from pattern
	for i := 1; i < depth; i++ {  // for every extra depth: move to the right of the pattern length positions, get length number of chars
		intraDir = filepath.Join(intraDir, pattern[length*i:length*(i+1)]) //  adding each time to intradir the extra characters with a filepath join
	}
	return intraDir
}

func GetDir(path string) string {
	if path == "" {
		path = GetHomeDirectory()
	}

	absolutePath, err := filepath.Abs(path)
	if err == nil {
		path = absolutePath
	}
	return absolutePath
}

func GetParent(path string) *string {
	isRoot := path[len(path)-1:] == "/"
	if isRoot {
		return nil
	} else {
		parentPath := filepath.Clean(path + "/..")
		return &parentPath
	}
}
