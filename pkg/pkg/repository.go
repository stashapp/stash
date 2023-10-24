package pkg

import (
	"context"
	"io"
	"io/fs"
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
	GetPackageZip(ctx context.Context, pkg RemotePackage) (io.ReadCloser, error)
}

// LocalRepository is a repository that can be used to store paks in.
type LocalRepository interface {
	LocalPackageLister
	PackageInstaller
	PackageDeleter
}

type PackageInstaller interface {
	writeFile(packageID string, name string, mode fs.FileMode, i io.Reader) error
	writeManifest(packageID string, m Manifest) error
}

type PackageDeleter interface {
	deleteFile(packageID string, name string) error
	deleteManifest(packageID string) error
	deletePackageDir(packageID string) error
}

type LocalPackageLister interface {
	// List returns all specs in the repository.
	List(ctx context.Context) ([]Manifest, error)
	getManifest(ctx context.Context, id string) (*Manifest, error)
}
