// Package http provides a repository implementation for HTTP.
package pkg

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"gopkg.in/yaml.v2"
)

// DefaultCacheTTL is the default time to live for the index cache.
const DefaultCacheTTL = 5 * time.Minute

// HttpRepository is a HTTP based repository.
// It is configured with a package list URL. Packages are located from the Path field of the package.
//
// The index is cached for the duration of CacheTTL. The first request after the cache expires will cause the index to be reloaded.
type HttpRepository struct {
	PackageListURL url.URL
	Client         *http.Client
}

// NewHttpRepository creates a new Repository. If client is nil then http.DefaultClient is used.
func NewHttpRepository(packageListURL url.URL, client *http.Client) *HttpRepository {
	if client == nil {
		client = http.DefaultClient
	}
	return &HttpRepository{
		PackageListURL: packageListURL,
		Client:         client,
	}
}

func (r *HttpRepository) Path() string {
	return r.PackageListURL.String()
}

func (r *HttpRepository) List(ctx context.Context) ([]RemotePackage, error) {
	u := r.PackageListURL

	f, err := r.getFile(ctx, u)
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

func (r *HttpRepository) GetPackageZip(ctx context.Context, pkg RemotePackage) (io.ReadCloser, error) {
	path := pkg.Path

	u := r.PackageListURL
	u.Path = path

	f, err := r.getFile(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("failed to get package list file: %w", err)
	}

	return f, nil
}

func (r *HttpRepository) getFile(ctx context.Context, u url.URL) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		// shouldn't happen
		return nil, err
	}

	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get remote file: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("failed to get remote file: %s", resp.Status)
	}

	return resp.Body, nil
}

var _ = RemoteRepository(&HttpRepository{})
