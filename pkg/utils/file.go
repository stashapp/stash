package utils

import (
	"archive/zip"
	"fmt"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"time"
)

type DuDetails struct {
	path  string
	mtime time.Time
	size  int64
}

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

// EnsudirAll will create a directory at the given path along with any necessary parents if they don't already exist
func EnsureDirAll(path string) error {
	return os.MkdirAll(path, 0755)
}

// RemoveDir removes the given file path along with all of its contents
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

// ListDir will return the contents of a given directory path as a string slice
func ListDir(path string) []string {
	if path == "" {
		path = GetHomeDirectory()
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

func PrintDuDetails(files []DuDetails) {
	for i, file := range files {
		fmt.Printf("%d: %s size: %d mod date: %s\n", i, file.path, file.size, file.mtime.Format("2006-01-02T15:04:05"))
	}
}

// Sort DuDetails slice by mod time, older first
func SortDuDetailsByMtime(files []DuDetails) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].mtime.After(files[j].mtime)
	})
}

// DuDir returns the size of the directory (only actual files are counted)
// The filesInfo slice is populated during the recursion
func DuDir(path string, info os.FileInfo, filesInfo *[]DuDetails) int64 {
	var details DuDetails
	var err error

	path, err = filepath.Abs(path)

	if err != nil {
		panic(err)
	}
	if !info.IsDir() {
		details.size = info.Size()
		details.mtime = info.ModTime()
		details.path = path
		*filesInfo = append(*filesInfo, details)
		return details.size
	}
	size := info.Size()
	dir, err := os.Open(path)
	if err != nil {
		fmt.Errorf("Error openig path %s: %s\n", path, err)
		return size
	}
	defer dir.Close()

	fis, err := dir.Readdir(-1)
	if err != nil {
		panic(err)
	}
	for _, fi := range fis {
		if fi.Name() == "." || fi.Name() == ".." {
			continue
		}

		size += DuDir(filepath.Join(path, fi.Name()), fi, filesInfo)
	}

	return size
}

// Reduce dir by at least size bytes
// Returns the size of removed files
func ReduceDir(files []DuDetails, size int64) int64 {
	var rmSize int64 = 0
	for _, file := range files {
		if rmSize <= size {
			err := os.Remove(file.path)

			if err != nil {
				fmt.Printf("Error removing file: %s\n", file.path)
				continue
			}
			rmSize += file.size
		} else {
			break
		}
	}
	return rmSize
}

// Return true if zip file is using 0 compression level
func IsZipFileUncompressed(path string) (bool, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		fmt.Printf("Error reading zip file %s: %s\n", path, err)
		return false, err
	} else {
		if r.File[0].Method == 0 { // for performance reasons we only check the compression
			return true, nil // level of the first file in the zip
		}
		r.Close()
	}
	return false, nil
}
