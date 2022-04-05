package fsutil

import (
	"fmt"
	"io/fs"
	"math"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// DirExists returns true if the given path exists and is a directory
func DirExists(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, fs.ErrNotExist
	}
	if !fileInfo.IsDir() {
		return false, fmt.Errorf("path is not a directory <%s>", path)
	}
	return true, nil
}

// IsPathInDir returns true if pathToCheck is within dir.
func IsPathInDir(dir, pathToCheck string) bool {
	rel, err := filepath.Rel(dir, pathToCheck)

	if err == nil {
		if !strings.HasPrefix(rel, "..") {
			return true
		}
	}

	return false
}

// GetHomeDirectory returns the path of the user's home directory.  ~ on Unix and C:\Users\UserName on Windows
func GetHomeDirectory() string {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	return currentUser.HomeDir
}

// EnsureDir will create a directory at the given path if it doesn't already exist
func EnsureDir(path string) error {
	exists, err := DirExists(path)
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

// GetIntraDir returns a string that can be added to filepath.Join to implement directory depth, "" on error
// eg given a pattern of 0af63ce3c99162e9df23a997f62621c5 and a depth of 2 length of 3
// returns 0af/63c or 0af\63c ( dependin on os)  that can be later used like this  filepath.Join(directory, intradir, basename)
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

// GetIntraDirID returns a string that can be added to filepath.Join to implement directory depth with
// sequential numbers.
// Where the pattern length is more than provided length, the maxID for length is prepended to the result, then
// GetIntraDir is called to break the pattern into smaller chunks of the provided length.
// eg: for an id of 1234 and a length of 3, the result is 1000/001
// for an id of 12345678 and length 3, the result is 1000/1000000/012/345
func GetIntraDirID(id int, length int) string {
	ret := ""
	thisLen := length
	max := int(math.Pow10(thisLen))
	for id >= max {
		ret = filepath.Join(ret, fmt.Sprintf("%d", max))

		thisLen += length
		max = int(math.Pow10(thisLen))
	}

	// convert the id to a string, padding to the required length
	idStr := fmt.Sprintf("%d", id)
	idStr = strings.Repeat("0", thisLen-len(idStr)) + idStr

	return filepath.Join(ret, GetIntraDir(idStr, (thisLen/length)-1, length))
}
