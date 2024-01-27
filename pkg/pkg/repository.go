package pkg

import (
	"context"
	"io"
)

// remoteRepository is a repository that can be used to get paks from.
type remoteRepository interface {
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
