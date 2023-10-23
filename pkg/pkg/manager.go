package pkg

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
)

// Manager manages the installation of paks.
type Manager struct {
	Local   LocalRepository
	Remotes []RemoteRepository
}

func (m *Manager) ListInstalled(ctx context.Context) (LocalPackageIndex, error) {
	installedList, err := m.Local.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing local packages: %w", err)
	}

	return localPackageIndexFromList(installedList), nil
}

func (m *Manager) ListRemote(ctx context.Context) (RemotePackageIndex, error) {
	var retList []RemotePackage

	for _, r := range m.Remotes {
		list, err := r.List(ctx)
		if err != nil {
			return nil, fmt.Errorf("listing remote packages: %w", err)
		}

		// add link to RemotePackage
		for i := range list {
			list[i].Repository = r
		}

		retList = append(retList, list...)
	}

	return remotePackageIndexFromList(retList), nil
}

func (m *Manager) InstalledStatus(ctx context.Context) (PackageStatusIndex, error) {
	// get all installed packages
	installed, err := m.ListInstalled(ctx)
	if err != nil {
		return nil, err
	}

	remoteList, err := m.ListRemote(ctx)
	if err != nil {
		return nil, err
	}

	ret := make(PackageStatusIndex)
	ret.populateLocal(installed, remoteList)

	return ret, nil
}

func (m *Manager) List(ctx context.Context) (PackageStatusIndex, error) {
	// get all installed packages
	installed, err := m.ListInstalled(ctx)
	if err != nil {
		return nil, err
	}

	remoteList, err := m.ListRemote(ctx)
	if err != nil {
		return nil, err
	}

	ret := make(PackageStatusIndex)
	ret.populateLocal(installed, remoteList)
	ret.populateRemote(remoteList)

	return ret, nil
}

func (m *Manager) Install(ctx context.Context, pkg RemotePackage) error {
	fromRemote, err := pkg.Repository.GetPackage(ctx, pkg)
	if err != nil {
		return fmt.Errorf("getting remote package: %w", err)
	}

	defer fromRemote.Close()

	d, err := io.ReadAll(fromRemote)
	if err != nil {
		return fmt.Errorf("reading package data: %w", err)
	}

	sha := fmt.Sprintf("%x", sha256.Sum256(d))
	if sha != pkg.Sha256 {
		return fmt.Errorf("package data (%s) does not match expected SHA256 (%s)", sha, pkg.Sha256)
	}

	zr, err := zip.NewReader(bytes.NewReader(d), int64(len(d)))
	if err != nil {
		return fmt.Errorf("reading zip data: %w", err)
	}

	if err := m.Local.InstallPackage(ctx, pkg, zr); err != nil {
		return fmt.Errorf("installing package: %w", err)
	}

	return nil
}

// Uninstall uninstalls the given package.
func (m *Manager) Uninstall(ctx context.Context, name string) error {
	if err := m.Local.DeletePackage(ctx, name); err != nil {
		return fmt.Errorf("deleting local package: %w", err)
	}

	return nil
}
