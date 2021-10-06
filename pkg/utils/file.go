package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
	"github.com/stashapp/stash/pkg/logger"
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
	}
	return false, err
}

// DirExists returns true if the given path exists and is a directory
func DirExists(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("path doesn't exist <%s>", path)
	}
	if !fileInfo.IsDir() {
		return false, fmt.Errorf("path is not a directory <%s>", path)
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

// ListDir will return the contents of a given directory path as a string slice
func ListDir(path string) ([]string, error) {
	var dirPaths []string
	files, err := ioutil.ReadDir(path)
	if err != nil {
		path = filepath.Dir(path)
		files, err = ioutil.ReadDir(path)
		if err != nil {
			return dirPaths, err
		}
	}
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		dirPaths = append(dirPaths, filepath.Join(path, file.Name()))
	}
	return dirPaths, nil
}

// GetHomeDirectory returns the path of the user's home directory.  ~ on Unix and C:\Users\UserName on Windows
func GetHomeDirectory() string {
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}
	return currentUser.HomeDir
}

func SafeMove(src, dst string) error {
	err := os.Rename(src, dst)

	if err != nil {
		logger.Errorf("[Util] unable to rename: \"%s\" due to %s. Falling back to copying.", src, err.Error())

		in, err := os.Open(src)
		if err != nil {
			return err
		}
		defer in.Close()

		out, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, in)
		if err != nil {
			return err
		}

		err = out.Close()
		if err != nil {
			return err
		}

		err = os.Remove(src)
		if err != nil {
			return err
		}
	}

	return nil
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

// WriteFile writes file to path creating parent directories if needed
func WriteFile(path string, file []byte) error {
	pathErr := EnsureDirAll(filepath.Dir(path))
	if pathErr != nil {
		return fmt.Errorf("cannot ensure path %s", pathErr)
	}

	err := ioutil.WriteFile(path, file, 0755)
	if err != nil {
		return fmt.Errorf("write error for thumbnail %s: %s ", path, err)
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

func GetDir(path string) string {
	if path == "" {
		path = GetHomeDirectory()
	}

	return path
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

// ServeFileNoCache serves the provided file, ensuring that the response
// contains headers to prevent caching.
func ServeFileNoCache(w http.ResponseWriter, r *http.Request, filepath string) {
	w.Header().Add("Cache-Control", "no-cache")

	http.ServeFile(w, r, filepath)
}

// MatchEntries returns a string slice of the entries in directory dir which
// match the regexp pattern. On error an empty slice is returned
// MatchEntries isn't recursive, only the specific 'dir' is searched
// without being expanded.
func MatchEntries(dir, pattern string) ([]string, error) {
	var res []string
	var err error

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	files, err := f.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if re.Match([]byte(file)) {
			res = append(res, filepath.Join(dir, file))
		}
	}
	return res, err
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

// GetNameFromPath returns the name of a file from its path
// if stripExtension is true the extension is omitted from the name
func GetNameFromPath(path string, stripExtension bool) string {
	fn := filepath.Base(path)
	if stripExtension {
		ext := filepath.Ext(fn)
		fn = strings.TrimSuffix(fn, ext)
	}
	return fn
}

// GetFunscriptPath returns the path of a file
// with the extension changed to .funscript
func GetFunscriptPath(path string) string {
	ext := filepath.Ext(path)
	fn := strings.TrimSuffix(path, ext)
	return fn + ".funscript"
}

// IsFsPathCaseSensitive checks the fs of the given path to see if it is case sensitive
// if the case sensitivity can not be determined false and an error != nil are returned
func IsFsPathCaseSensitive(path string) (bool, error) {
	// The case sensitivity of the fs of "path" is determined by case flipping
	// the first letter rune from the base string of the path
	// If the resulting flipped path exists then the fs should not be case sensitive
	// ( we check the file mod time to avoid matching an existing path )

	fi, err := os.Stat(path)
	if err != nil { // path cannot be stat'd
		return false, err
	}

	base := filepath.Base(path)
	fBase, err := FlipCaseSingle(base)
	if err != nil { // cannot be case flipped
		return false, err
	}
	i := strings.LastIndex(path, base)
	if i < 0 { // shouldn't happen
		return false, fmt.Errorf("could not case flip path %s", path)
	}

	flipped := []byte(path)
	for _, c := range []byte(fBase) { // replace base of path with the flipped one ( we need to flip the base or last dir part )
		flipped[i] = c
		i++
	}

	fiCase, err := os.Stat(string(flipped))
	if err != nil { // cannot stat the case flipped path
		return true, nil // fs of path should be case sensitive
	}

	if fiCase.ModTime() == fi.ModTime() { // file path exists and is the same
		return false, nil // fs of path is not case sensitive
	}
	return false, fmt.Errorf("can not determine case sensitivity of path %s", path)
}

func FindInPaths(paths []string, baseName string) string {
	for _, p := range paths {
		filePath := filepath.Join(p, baseName)
		if exists, _ := FileExists(filePath); exists {
			return filePath
		}
	}

	return ""
}

// MatchExtension returns true if the extension of the provided path
// matches any of the provided extensions.
func MatchExtension(path string, extensions []string) bool {
	ext := filepath.Ext(path)
	for _, e := range extensions {
		if strings.EqualFold(ext, "."+e) {
			return true
		}
	}

	return false
}
