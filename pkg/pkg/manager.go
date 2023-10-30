package pkg

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"time"
)

// Manager manages the installation of paks.
type Manager struct {
	Local       LocalRepository
	remoteCache map[string]*remotePackageCache

	Client *http.Client

	// CacheTTL is the time to live for the index cache.
	// The index is cached for this duration. The first request after the cache
	// expires will cause the index to be reloaded.
	// This applies only to http remote indexes.
	CacheTTL time.Duration
}

func (m *Manager) remoteFromURL(path string) (*httpRepository, error) {
	u, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("parsing path: %w", err)
	}

	return newHttpRepository(*u, m.Client), nil
}

func (m *Manager) checkCacheExpired(path string) {
	if m.remoteCache == nil {
		m.remoteCache = make(map[string]*remotePackageCache)
	}

	cache := m.remoteCache[path]

	if cache == nil {
		return
	}

	if time.Since(cache.cacheTime) > m.CacheTTL {
		m.remoteCache[path] = nil
	}
}

func (m *Manager) ListInstalled(ctx context.Context) (LocalPackageIndex, error) {
	installedList, err := m.Local.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing local packages: %w", err)
	}

	return localPackageIndexFromList(installedList), nil
}

func (m *Manager) ListRemote(ctx context.Context, remoteURL string) (RemotePackageIndex, error) {
	m.checkCacheExpired(remoteURL)
	cache := m.remoteCache[remoteURL]

	if cache != nil {
		return cache.cachedIndex, nil
	}

	r, err := m.remoteFromURL(remoteURL)
	if err != nil {
		return nil, fmt.Errorf("creating remote repository: %w", err)
	}

	list, err := r.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing remote packages: %w", err)
	}

	// add link to RemotePackage
	for i := range list {
		list[i].Repository = r
	}

	ret := remotePackageIndexFromList(list)

	// only cache remote http results
	if r.packageListURL.Scheme != "file" {
		m.remoteCache[remoteURL] = &remotePackageCache{
			cachedIndex: ret,
			cacheTime:   time.Now(),
		}
	}

	return ret, nil
}

func (m *Manager) InstalledStatus(ctx context.Context) (PackageStatusIndex, error) {
	// get all installed packages
	installed, err := m.ListInstalled(ctx)
	if err != nil {
		return nil, err
	}

	// get remotes for all installed packages
	allRemoteList := make(RemotePackageIndex)

	remoteURLs := installed.remoteURLs()
	for _, remoteURL := range remoteURLs {
		remoteList, err := m.ListRemote(ctx, remoteURL)
		if err != nil {
			return nil, err
		}

		allRemoteList.merge(remoteList)
	}

	ret := make(PackageStatusIndex)
	ret.populateLocal(installed, allRemoteList)

	return ret, nil
}

func (m *Manager) packageByID(ctx context.Context, remoteURL string, id string) (*RemotePackage, error) {
	l, err := m.ListRemote(ctx, remoteURL)
	if err != nil {
		return nil, err
	}

	pkg, found := l[id]
	if !found {
		return nil, nil
	}

	return &pkg, nil
}

func (m *Manager) Install(ctx context.Context, remoteURL string, id string) error {
	remote, err := m.remoteFromURL(remoteURL)
	if err != nil {
		return fmt.Errorf("creating remote repository: %w", err)
	}

	pkg, err := m.packageByID(ctx, remoteURL, id)
	if err != nil {
		return fmt.Errorf("getting remote package: %w", err)
	}

	fromRemote, err := remote.GetPackageZip(ctx, *pkg)
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

	// uninstall existing package if present
	if _, err := m.Local.getManifest(ctx, pkg.ID); err == nil {
		if err := m.deletePackageFiles(ctx, pkg.ID); err != nil {
			return fmt.Errorf("uninstalling existing package: %w", err)
		}
	}

	if err := m.installPackage(*pkg, zr); err != nil {
		return fmt.Errorf("installing package: %w", err)
	}

	return nil
}

func (m *Manager) installPackage(pkg RemotePackage, zr *zip.Reader) error {
	manifest := Manifest{
		ID:              pkg.ID,
		Name:            pkg.Name,
		PackageMetadata: pkg.PackageMetadata,
		PackageVersion:  pkg.PackageVersion,
		RepositoryURL:   pkg.Repository.Path(),
	}

	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			continue
		}

		i, err := f.Open()
		if err != nil {
			return err
		}

		fn := filepath.Clean(f.Name)
		if err := m.Local.writeFile(pkg.ID, fn, f.Mode(), i); err != nil {
			i.Close()
			return fmt.Errorf("writing file %q: %w", fn, err)
		}

		i.Close()
		manifest.Files = append(manifest.Files, fn)
	}

	if err := m.Local.writeManifest(pkg.ID, manifest); err != nil {
		return fmt.Errorf("writing manifest: %w", err)
	}

	return nil
}

// Uninstall uninstalls the given package.
func (m *Manager) Uninstall(ctx context.Context, id string) error {
	if err := m.deletePackageFiles(ctx, id); err != nil {
		return fmt.Errorf("deleting local package: %w", err)
	}

	// also delete the directory
	// ignore errors
	_ = m.Local.deletePackageDir(id)

	return nil
}

func (m *Manager) deletePackageFiles(ctx context.Context, id string) error {
	manifest, err := m.Local.getManifest(ctx, id)
	if err != nil {
		return fmt.Errorf("getting manifest: %w", err)
	}

	for _, f := range manifest.Files {
		if err := m.Local.deleteFile(id, f); err != nil {
			// ignore
			continue
		}
	}

	if err := m.Local.deleteManifest(id); err != nil {
		return fmt.Errorf("deleting manifest: %w", err)
	}

	return nil
}
