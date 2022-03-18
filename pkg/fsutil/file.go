package fsutil

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// SafeMove attempts to move the file with path src to dest using os.Rename. If this fails, then it copies src to dest, then deletes src.
func SafeMove(src, dst string) error {
	err := os.Rename(src, dst)

	if err != nil {
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
		return fmt.Errorf("cannot ensure path %s", pathErr)
	}

	err := os.WriteFile(path, file, 0755)
	if err != nil {
		return fmt.Errorf("write error for thumbnail %s: %s ", path, err)
	}
	return nil
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
