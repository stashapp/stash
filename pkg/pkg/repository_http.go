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

	"gopkg.in/yaml.v2"
)

// DefaultCacheTTL is the default time to live for the index cache.
const DefaultCacheTTL = 5 * time.Minute

// httpRepository is a HTTP based repository.
// It is configured with a package list URL. Packages are located from the Path field of the package.
//
// The index is cached for the duration of CacheTTL. The first request after the cache expires will cause the index to be reloaded.
type httpRepository struct {
	packageListURL url.URL
	client         *http.Client
}

// newHttpRepository creates a new Repository. If client is nil then http.DefaultClient is used.
func newHttpRepository(packageListURL url.URL, client *http.Client) *httpRepository {
	if client == nil {
		client = http.DefaultClient
	}
	return &httpRepository{
		packageListURL: packageListURL,
		client:         client,
	}
}

func (r *httpRepository) Path() string {
	return r.packageListURL.String()
}

func (r *httpRepository) List(ctx context.Context) ([]RemotePackage, error) {
	u := r.packageListURL

	// the package list URL may be file://, in which case we need to use the local file system
	var (
		f   io.ReadCloser
		err error
	)
	if u.Scheme == "file" {
		f, err = r.getLocalFile(ctx, u.Path)
	} else {
		f, err = r.getFile(ctx, u)
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
		f, err = r.getFile(ctx, u)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get package file: %w", err)
	}

	return f, nil
}

func (r *httpRepository) getFile(ctx context.Context, u url.URL) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
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

	return resp.Body, nil
}

func (r *httpRepository) getLocalFile(ctx context.Context, path string) (fs.File, error) {
	ret, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get file %q: %w", path, err)
	}

	return ret, nil
}

var _ = remoteRepository(&httpRepository{})
