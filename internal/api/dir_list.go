package api

import (
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/text/collate"
)

type dirLister []fs.DirEntry

func (s dirLister) Len() int {
	return len(s)
}

func (s dirLister) Swap(i, j int) {
	s[j], s[i] = s[i], s[j]
}

func (s dirLister) Bytes(i int) []byte {
	return []byte(s[i].Name())
}

// listDir will return the contents of a given directory path as a string slice
func listDir(col *collate.Collator, path string) ([]string, error) {
	var dirPaths []string
	files, err := os.ReadDir(path)
	if err != nil {
		path = filepath.Dir(path)
		files, err = os.ReadDir(path)
		if err != nil {
			return dirPaths, err
		}
	}

	if col != nil {
		col.Sort(dirLister(files))
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		dirPaths = append(dirPaths, filepath.Join(path, file.Name()))
	}
	return dirPaths, nil
}
