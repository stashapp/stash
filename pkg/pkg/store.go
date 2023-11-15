package pkg

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// ManifestFile is the default filename for the package manifest.
const ManifestFile = "manifest"

// Store is a folder-based local repository.
// Packages are installed in their own directory under BaseDir.
// The package details are stored in a file named based on PackageFile.
type Store struct {
	BaseDir string
	// ManifestFile is the filename of the package file.
	ManifestFile string
}

// sub returns a new Store with the given path appended to the BaseDir.
func (r *Store) sub(path string) *Store {
	if path == "" || path == "." {
		return r
	}

	return &Store{
		BaseDir:      filepath.Join(r.BaseDir, path),
		ManifestFile: r.ManifestFile,
	}
}

func (r *Store) List(ctx context.Context) ([]Manifest, error) {
	e, err := os.ReadDir(r.BaseDir)
	// ignore if directory cannot be read
	if err != nil {
		return nil, nil
	}

	var ret []Manifest

	for _, ee := range e {
		if !ee.IsDir() {
			// ignore non-directories
			continue
		}

		pkg, err := r.getManifest(ctx, ee.Name())
		if err != nil {
			// ignore if manifest does not exist
			if errors.Is(err, os.ErrNotExist) {
				continue
			}

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

func (r *Store) getManifest(ctx context.Context, packageID string) (*Manifest, error) {
	pfp := r.manifestPath(packageID)

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

func (r *Store) ensurePackageExists(packageID string) error {
	// ensure the manifest file exists
	if _, err := os.Stat(r.manifestPath(packageID)); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("package %q does not exist", packageID)
		}
	}

	return nil
}

func (r *Store) writeFile(packageID string, name string, mode fs.FileMode, i io.Reader) error {
	fn := filepath.Join(r.packageDir(packageID), name)

	if err := os.MkdirAll(filepath.Dir(fn), os.ModePerm); err != nil {
		return fmt.Errorf("creating directory %v: %w", fn, err)
	}

	o, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}

	defer o.Close()

	if _, err := io.Copy(o, i); err != nil {
		return err
	}

	return nil
}

func (r *Store) writeManifest(packageID string, m Manifest) error {
	pfp := r.manifestPath(packageID)
	data, err := yaml.Marshal(m)
	if err != nil {
		return fmt.Errorf("marshaling manifest: %w", err)
	}

	if err := os.WriteFile(pfp, data, os.ModePerm); err != nil {
		return fmt.Errorf("writing manifest file %q: %w", pfp, err)
	}

	return nil
}

func (r *Store) deleteFile(packageID string, name string) error {
	// ensure the package exists
	if err := r.ensurePackageExists(packageID); err != nil {
		return err
	}

	pkgDir := r.packageDir(packageID)
	fp := filepath.Join(pkgDir, name)

	return os.Remove(fp)
}

func (r *Store) deleteManifest(packageID string) error {
	return r.deleteFile(packageID, r.ManifestFile)
}

func (r *Store) deletePackageDir(packageID string) error {
	return os.Remove(r.packageDir(packageID))
}
