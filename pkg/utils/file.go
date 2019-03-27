package utils

import (
	"fmt"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

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
