package fsutil

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/stashapp/stash/pkg/models"
)

func GetFileSize(f models.FS, path string, info fs.FileInfo) (int64, error) {
	// #2196/#3042 - replace size with target size if file is a symlink
	if info.Mode()&os.ModeSymlink == os.ModeSymlink {
		targetInfo, err := f.Stat(path)
		if err != nil {
			return 0, fmt.Errorf("reading info for symlink %q: %w", path, err)
		}
		return targetInfo.Size(), nil
	}

	return info.Size(), nil
}
