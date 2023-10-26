// Package http provides a repository implementation for HTTP.
package pkg

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// FSRepository is a file based repository.
// It is configured with a package list path. Packages are located in the same directory as the package list.
type FSRepository struct {
	PackageListPath string
}

func (r *FSRepository) Path() string {
	return r.PackageListPath
}

func (r *FSRepository) dir() string {
	return filepath.Dir(r.PackageListPath)
}

func (r *FSRepository) List(ctx context.Context) ([]RemotePackage, error) {
	f, err := r.getFile(ctx, r.PackageListPath)
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

func (r *FSRepository) getPackagePath(pkg RemotePackage) string {
	return filepath.Join(r.dir(), pkg.Path)
}

func (r *FSRepository) PackageByID(ctx context.Context, id string) (*RemotePackage, error) {
	list, err := r.List(ctx)
	if err != nil {
		return nil, err
	}

	for i := range list {
		if list[i].ID == id {
			return &list[i], nil
		}
	}

	return nil, fmt.Errorf("package not found")
}

func (r *FSRepository) GetPackageZip(ctx context.Context, pkg RemotePackage) (io.ReadCloser, error) {
	f, err := r.getFile(ctx, r.getPackagePath(pkg))
	if err != nil {
		return nil, fmt.Errorf("failed to get package list file: %w", err)
	}

	return f, nil
}

func (r *FSRepository) getFile(ctx context.Context, path string) (fs.File, error) {
	ret, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get file %q: %w", path, err)
	}

	return ret, nil
}

var _ = remoteRepository(&FSRepository{})
