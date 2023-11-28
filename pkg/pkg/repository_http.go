// Package http provides a repository implementation for HTTP.
package pkg

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/stashapp/stash/pkg/logger"
	"gopkg.in/yaml.v2"
)

// httpRepository is a HTTP based repository.
// It is configured with a package list URL. Packages are located from the Path field of the package.
//
// The index is cached for the duration of CacheTTL. The first request after the cache expires will cause the index to be reloaded.
type httpRepository struct {
	packageListURL url.URL
	client         *http.Client

	cache *repositoryCache
}

// newHttpRepository creates a new Repository. If client is nil then http.DefaultClient is used.
func newHttpRepository(packageListURL url.URL, client *http.Client, cache *repositoryCache) *httpRepository {
	if client == nil {
		client = http.DefaultClient
	}
	return &httpRepository{
		packageListURL: packageListURL,
		client:         client,
		cache:          cache,
	}
}

func (r *httpRepository) Path() string {
	return r.packageListURL.String()
}

func (r *httpRepository) List(ctx context.Context) ([]RemotePackage, error) {
	u := r.packageListURL

	// the package list URL may be file://, in which case we need to use the local file system
	var (
		f       io.ReadCloser
		modTime *time.Time
		err     error
	)

	isLocal := u.Scheme == "file"

	if isLocal {
		f, err = r.getLocalFile(ctx, u.Path)
	} else {
		// try to get the cached list first
		var cachedList []RemotePackage
		cachedList, err = r.getCachedList(ctx, u)
		if err != nil {
			return nil, fmt.Errorf("failed to get cached package list: %w", err)
		}

		if cachedList != nil {
			return cachedList, nil
		}

		f, modTime, err = r.getFile(ctx, u)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get package list: %w", err)
	}

	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read package list: %w", err)
	}

	var index []RemotePackage
	if err := yaml.Unmarshal(data, &index); err != nil {
		return nil, fmt.Errorf("reading package list: %w", err)
	}

	// cache if not local file
	if !isLocal {
		r.cache.cacheList(u.String(), *modTime, index)
	}

	return index, nil
}

func isURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && u.Scheme != "" && (u.Scheme == "file" || u.Host != "")
}

func (r *httpRepository) resolvePath(p string) url.URL {
	// if the path can be resolved to a URL, then use that
	if isURL(p) {
		// isURL ensures URL is valid
		u, _ := url.Parse(p)
		return *u
	}

	// otherwise, determine if the path is relative or absolute
	// if it's relative, then join it with the package list URL
	u := r.packageListURL

	if path.IsAbs(p) {
		u.Path = p
	} else {
		u.Path = path.Join(path.Dir(u.Path), p)
	}

	return u
}

func (r *httpRepository) GetPackageZip(ctx context.Context, pkg RemotePackage) (io.ReadCloser, error) {
	p := pkg.Path

	u := r.resolvePath(p)

	var (
		f   io.ReadCloser
		err error
	)

	// the package list URL may be file://, in which case we need to use the local file system
	// the package zip path may be a URL. A remotely hosted list may _not_ use local files.
	if u.Scheme == "file" {
		if r.packageListURL.Scheme != "file" {
			return nil, fmt.Errorf("%s is invalid for a remotely hosted package list", u.String())
		}

		f, err = r.getLocalFile(ctx, u.Path)
	} else {
		f, _, err = r.getFile(ctx, u)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get package file: %w", err)
	}

	return f, nil
}

// getFileCached tries to get the list from the local cache.
// If it is not found or is stale, then nil is returned.
func (r *httpRepository) getCachedList(ctx context.Context, u url.URL) ([]RemotePackage, error) {
	// check if the file is in the cache first
	localModTime := r.cache.lastModified(u.String())

	if localModTime != nil {
		// get the update time of the file
		req, err := http.NewRequestWithContext(ctx, http.MethodHead, u.String(), nil)
		if err != nil {
			// shouldn't happen
			return nil, err
		}

		resp, err := r.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to get remote file: %w", err)
		}

		if resp.StatusCode >= 400 {
			return nil, fmt.Errorf("failed to get remote file: %s", resp.Status)
		}

		lastModified := resp.Header.Get("Last-Modified")
		if lastModified != "" {
			remoteModTime, _ := time.Parse(http.TimeFormat, lastModified)

			if !remoteModTime.After(*localModTime) {
				logger.Debugf("cached version of %s is equal or newer than remote", u.String())
				return r.cache.getPackageList(u.String()), nil
			}
		}

		logger.Debugf("cached version of %s is older than remote", u.String())
	}

	return nil, nil
}

// getFile gets the file from the remote server. Returns the file and the last modified time.
func (r *httpRepository) getFile(ctx context.Context, u url.URL) (io.ReadCloser, *time.Time, error) {
	logger.Debugf("fetching %s", u.String())
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		// shouldn't happen
		return nil, nil, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get remote file: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, nil, fmt.Errorf("failed to get remote file: %s", resp.Status)
	}

	lastModified := resp.Header.Get("Last-Modified")
	var remoteModTime time.Time
	if lastModified != "" {
		remoteModTime, _ = time.Parse(http.TimeFormat, lastModified)
	}

	return resp.Body, &remoteModTime, nil
}

func (r *httpRepository) getLocalFile(ctx context.Context, path string) (fs.File, error) {
	ret, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get file %q: %w", path, err)
	}

	return ret, nil
}

var _ = remoteRepository(&httpRepository{})
