// Package http provides a repository implementation for HTTP.
package pkg

import (
	"context"
	"fmt"
	"io"
	"io/fs"

	"gopkg.in/yaml.v2"
)

// FSRepository is a HTTP based repository.
// It is configured with a package list URL. Packages are located from the Path field of the package.
//
// The index is cached for the duration of CacheTTL. The first request after the cache expires will cause the index to be reloaded.
type FSRepository struct {
	Root                fs.FS
	PackageListFilename string
}

func (r *FSRepository) List(ctx context.Context) ([]RemotePackage, error) {
	f, err := r.getFile(ctx, r.PackageListFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to get package list: %w", err)
	}

	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("reading package list: %w", err)
	}

	var index []RemotePackage
	if err := yaml.Unmarshal(data, &index); err != nil {
		return nil, fmt.Errorf("reading package list: %w", err)
	}

	return index, nil
}

func (r *FSRepository) GetPackage(ctx context.Context, pkg RemotePackage) (io.ReadCloser, error) {
	f, err := r.getFile(ctx, pkg.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to get package list file: %w", err)
	}

	return f, nil
}

func (r *FSRepository) getFile(ctx context.Context, path string) (fs.File, error) {
	ret, err := r.Root.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get file %q: %w", path, err)
	}

	return ret, nil
}

var _ = RemoteRepository(&FSRepository{})
