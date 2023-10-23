package pkg

import (
	"archive/zip"
	"context"
	"io"
)

// RemoteRepository is a repository that can be used to get paks from.
type RemoteRepository interface {
	RemotePackageLister
	RemotePackageGetter
}

type RemotePackageLister interface {
	// List returns all specs in the repository.
	List(ctx context.Context) ([]RemotePackage, error)
}

type RemotePackageGetter interface {
	GetPackage(ctx context.Context, pkg RemotePackage) (io.ReadCloser, error)
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
	DeletePackage(ctx context.Context, name string) error
}

type LocalPackageLister interface {
	// List returns all specs in the repository.
	List(ctx context.Context) ([]Manifest, error)
}
