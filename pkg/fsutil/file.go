package fsutil

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// CopyFile copies the contents of the file at srcpath to a regular file at dstpath.
// It will copy the last modified timestamp
// If dstpath already exists the function will fail.
func CopyFile(srcpath, dstpath string) (err error) {
	r, err := os.Open(srcpath)
	if err != nil {
		return err
	}

	w, err := os.OpenFile(dstpath, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0666)
	if err != nil {
		r.Close() // We need to close the input file as the defer below would not be called.
		return err
	}

	defer func() {
		r.Close() // ok to ignore error: file was opened read-only.
		e := w.Close()
		// Report the error from w.Close, if any.
		// But do so only if there isn't already an outgoing error.
		if e != nil && err == nil {
			err = e
		}
		// Copy modified time
		if err == nil {
			// io.Copy succeeded, we should fix the dstpath timestamp
			srcFileInfo, e := os.Stat(srcpath)
			if e != nil {
				err = e
				return
			}

			e = os.Chtimes(dstpath, srcFileInfo.ModTime(), srcFileInfo.ModTime())
			if e != nil {
				err = e
			}
		}
	}()

	_, err = io.Copy(w, r)
	return err
}

// SafeMove attempts to move the file with path src to dest using os.Rename. If this fails, then it copies src to dest, then deletes src.
// If the copy fails, or the delete fails, the function will return an error.
func SafeMove(src, dst string) error {
	err := os.Rename(src, dst)

	if err != nil {
		copyErr := CopyFile(src, dst)
		if copyErr != nil {
			return fmt.Errorf("copying file during SaveMove failed with: '%w'; renaming file failed previously with: '%v'", copyErr, err)
		}

		removeErr := os.Remove(src)
		if removeErr != nil {
			// if we can't remove the old file, remove the new one and fail
			_ = os.Remove(dst)
			return fmt.Errorf("removing old file during SafeMove failed with: '%w'; renaming file failed previously with: '%v'", removeErr, err)
		}
	}

	return nil
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

// FindInPaths returns the path to baseName in the first path where it exists from paths.
func FindInPaths(paths []string, baseName string) string {
	for _, p := range paths {
		filePath := filepath.Join(p, baseName)
		if exists, _ := FileExists(filePath); exists {
			return filePath
		}
	}

	return ""
}

// FileExists returns true if the given path exists and is a file.
// This function returns false and the error encountered if the call to os.Stat fails.
func FileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return !info.IsDir(), nil
	}
	return false, err
}

// WriteFile writes file to path creating parent directories if needed
func WriteFile(path string, file []byte) error {
	pathErr := EnsureDirAll(filepath.Dir(path))
	if pathErr != nil {
		return fmt.Errorf("cannot ensure path exists: %w", pathErr)
	}

	return os.WriteFile(path, file, 0755)
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

var (
	replaceCharsRE = regexp.MustCompile(`[&=\\/:*"?_ ]`)
	removeCharsRE  = regexp.MustCompile(`[^[:alnum:]-.]`)
	multiHyphenRE  = regexp.MustCompile(`\-+`)
)

// SanitiseBasename returns a file basename removing any characters that are illegal or problematic to use in the filesystem.
// It appends a short hash of the original string to ensure uniqueness.
func SanitiseBasename(v string) string {
	// Generate a short hash for uniqueness
	hash := sha1.Sum([]byte(v))
	shortHash := hex.EncodeToString(hash[:4]) // Use the first 4 bytes of the hash

	v = strings.TrimSpace(v)

	// replace illegal filename characters with -
	v = replaceCharsRE.ReplaceAllString(v, "-")

	// remove other characters
	v = removeCharsRE.ReplaceAllString(v, "")

	// remove multiple hyphens
	v = multiHyphenRE.ReplaceAllString(v, "-")

	return strings.TrimSpace(v) + "-" + shortHash
}

// GetExeName returns the name of the given executable for the current platform.
// One windows it returns the name with the .exe extension.
func GetExeName(base string) string {
	if runtime.GOOS == "windows" {
		return base + ".exe"
	}
	return base
}
