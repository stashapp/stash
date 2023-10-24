package pkg

import (
	"archive/zip"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// RemoteRepository is a repository that can be used to get paks from.
type RemoteRepository interface {
	RemotePackageLister
	RemotePackageGetter
	Path() string
}

type RemotePackageLister interface {
	// List returns all specs in the repository.
	List(ctx context.Context) ([]RemotePackage, error)
}

type RemotePackageGetter interface {
	PackageByID(ctx context.Context, id string) (*RemotePackage, error)
	GetPackageZip(ctx context.Context, pkg RemotePackage) (io.ReadCloser, error)
}

// LocalRepository is a repository that can be used to store paks in.
type LocalRepository interface {
	LocalPackageLister
	PackageInstaller
	PackageDeleter
}

type PackageInstaller interface {
	InstallPackage(ctx context.Context, pkg RemotePackage, data *zip.Reader) error
}

type PackageDeleter interface {
	DeletePackage(ctx context.Context, id string) error
}

type LocalPackageLister interface {
	// List returns all specs in the repository.
	List(ctx context.Context) ([]Manifest, error)
}

func NewRemoteRepository(path string, httpClient *http.Client, httpTTL time.Duration) (RemoteRepository, error) {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		u, err := url.Parse(path)
		if err != nil {
			return nil, err
		}
		return NewHttpRepository(*u, httpClient, httpTTL), nil
	} else {
		return &FSRepository{
			PackageListPath: path,
		}, nil
	}
}
