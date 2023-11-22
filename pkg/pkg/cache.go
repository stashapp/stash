package pkg

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/stashapp/stash/pkg/hash/md5"
	"github.com/stashapp/stash/pkg/logger"
)

const cacheSubDir = "package_lists"

type repositoryCache struct {
	cachePath string
}

func (c *repositoryCache) path(url string) string {
	// convert the url to md5
	hash := md5.FromString(url)

	return filepath.Join(c.cachePath, cacheSubDir, hash)
}

func (c *repositoryCache) lastModified(url string) *time.Time {
	if c == nil {
		return nil
	}

	path := c.path(url)
	s, err := os.Stat(path)
	if err != nil {
		// ignore
		logger.Debugf("error getting cached file %s: %v", path, err)
		return nil
	}

	ret := s.ModTime()
	return &ret
}

func (c *repositoryCache) getPackageList(url string) (io.ReadCloser, error) {
	path := c.path(url)
	ret, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get file %q: %w", path, err)
	}

	return ret, nil
}

func (c *repositoryCache) cacheFile(url string, data io.ReadCloser) (io.ReadCloser, error) {
	if c == nil {
		return data, nil
	}

	path := c.path(url)

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		// ignore, just return the original file
		logger.Debugf("error creating cache path %s: %v", filepath.Dir(path), err)
		return data, nil
	}

	f, err := os.Create(path)
	if err != nil {
		// ignore, just return the original file
		logger.Debugf("error creating cached file %s: %v", path, err)
		return data, nil
	}

	defer data.Close()
	if _, err := io.Copy(f, data); err != nil {
		_ = f.Close()
		return nil, fmt.Errorf("writing to cache file %s - %w", path, err)
	}

	_ = f.Close()
	return c.getPackageList(url)
}
