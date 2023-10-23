package pkg

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Store is a folder-based local repository.
// Packages are installed in their own directory under BaseDir.
// The package details are stored in a file named based on PackageFile.
type Store struct {
	BaseDir string
	// ManifestFile is the filename of the package file.
	ManifestFile string
}

func (r *Store) List(ctx context.Context) ([]Manifest, error) {
	e, err := os.ReadDir(r.BaseDir)
	if err != nil {
		return nil, fmt.Errorf("listing directory %q: %w", r.BaseDir, err)
	}

	var ret []Manifest

	for _, ee := range e {
		if !ee.IsDir() {
			// ignore non-directories
			continue
		}

		pkg, err := r.readManifest(ee)
		if err != nil {
			return nil, err
		}

		ret = append(ret, *pkg)
	}

	return ret, nil
}

func (r *Store) packageDir(id string) string {
	return filepath.Join(r.BaseDir, id)
}

func (r *Store) manifestPath(id string) string {
	return filepath.Join(r.packageDir(id), r.ManifestFile)
}

func (r *Store) readManifest(e fs.DirEntry) (*Manifest, error) {
	pfp := r.manifestPath(e.Name())
	data, err := os.ReadFile(pfp)
	if err != nil {
		return nil, fmt.Errorf("reading manifest file %q: %w", pfp, err)
	}

	var manifest Manifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("reading manifest file %q: %w", pfp, err)
	}

	return &manifest, nil
}

func (r *Store) InstallPackage(ctx context.Context, pkg RemotePackage, zr *zip.Reader) error {
	// assume data is zip encoded
	// assume zip contains the package file

	// create the directory for the package
	pkgDir := r.packageDir(pkg.ID)
	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		return fmt.Errorf("creating package directory %q: %w", pkgDir, err)
	}

	// copy the contents of the zip into the package directory, overwriting existing
	if err := unzipFile(pkgDir, zr); err != nil {
		return fmt.Errorf("unzipping package data: %w", err)
	}

	return nil
}

func unzipFile(dest string, zr *zip.Reader) error {
	for _, f := range zr.File {
		fn := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fn, os.ModePerm); err != nil {
				return fmt.Errorf("creating directory %v: %w", fn, err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fn), os.ModePerm); err != nil {
			return fmt.Errorf("creating directory %v: %w", fn, err)
		}

		o, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		i, err := f.Open()
		if err != nil {
			o.Close()
			return err
		}

		if _, err := io.Copy(o, i); err != nil {
			o.Close()
			i.Close()
			return err
		}

		o.Close()
		i.Close()
	}

	return nil
}

func (r *Store) DeletePackage(ctx context.Context, id string) error {
	// ensure the manifest file exists
	if _, err := os.Stat(r.manifestPath(id)); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("package %q does not exist", id)
		}
	}

	pkgDir := r.packageDir(id)
	return os.RemoveAll(pkgDir)
}

// ensure LocalRepository implements LocalRepository
var _ = LocalRepository(&Store{})
