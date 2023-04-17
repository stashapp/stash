package file

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"strconv"
	"syscall"
	"time"

	"github.com/stashapp/stash/pkg/logger"
)

// ID represents an ID of a file.
type ID int32

func (i ID) String() string {
	return strconv.Itoa(int(i))
}

// DirEntry represents a file or directory in the file system.
type DirEntry struct {
	ZipFileID *ID `json:"zip_file_id"`

	// transient - not persisted
	// only guaranteed to have id, path and basename set
	ZipFile File

	ModTime time.Time `json:"mod_time"`
}

func (e *DirEntry) info(fs FS, path string) (fs.FileInfo, error) {
	if e.ZipFile != nil {
		zipPath := e.ZipFile.Base().Path
		zfs, err := fs.OpenZip(zipPath)
		if err != nil {
			return nil, err
		}
		defer zfs.Close()
		fs = zfs
	}
	// else assume os file

	ret, err := fs.Lstat(path)
	return ret, err
}

// File represents a file in the file system.
type File interface {
	Base() *BaseFile
	SetFingerprints(fp Fingerprints)
	Open(fs FS) (io.ReadCloser, error)
}

// BaseFile represents a file in the file system.
type BaseFile struct {
	ID ID `json:"id"`

	DirEntry

	// resolved from parent folder and basename only - not stored in DB
	Path string `json:"path"`

	Basename       string   `json:"basename"`
	ParentFolderID FolderID `json:"parent_folder_id"`

	Fingerprints Fingerprints `json:"fingerprints"`

	Size int64 `json:"size"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SetFingerprints sets the fingerprints of the file.
// If a fingerprint of the same type already exists, it is overwritten.
func (f *BaseFile) SetFingerprints(fp Fingerprints) {
	for _, v := range fp {
		f.SetFingerprint(v)
	}
}

// SetFingerprint sets the fingerprint of the file.
// If a fingerprint of the same type already exists, it is overwritten.
func (f *BaseFile) SetFingerprint(fp Fingerprint) {
	for i, existing := range f.Fingerprints {
		if existing.Type == fp.Type {
			f.Fingerprints[i] = fp
			return
		}
	}

	f.Fingerprints = append(f.Fingerprints, fp)
}

// Base is used to fulfil the File interface.
func (f *BaseFile) Base() *BaseFile {
	return f
}

func (f *BaseFile) Open(fs FS) (io.ReadCloser, error) {
	if f.ZipFile != nil {
		zipPath := f.ZipFile.Base().Path
		zfs, err := fs.OpenZip(zipPath)
		if err != nil {
			return nil, err
		}

		return zfs.OpenOnly(f.Path)
	}

	return fs.Open(f.Path)
}

func (f *BaseFile) Info(fs FS) (fs.FileInfo, error) {
	return f.info(fs, f.Path)
}

func (f *BaseFile) Serve(fs FS, w http.ResponseWriter, r *http.Request) error {
	w.Header().Add("Cache-Control", "max-age=604800000") // 1 Week

	reader, err := f.Open(fs)
	if err != nil {
		return err
	}

	defer reader.Close()

	rsc, ok := reader.(io.ReadSeeker)
	if !ok {
		// fallback to direct copy
		data, err := io.ReadAll(reader)
		if err != nil {
			return err
		}

		k, err := w.Write(data)
		if err != nil && !errors.Is(err, syscall.EPIPE) {
			logger.Warnf("error serving file (wrote %v bytes out of %v): %v", k, len(data), err)
		}

		return nil
	}

	http.ServeContent(w, r, f.Basename, f.ModTime, rsc)
	return nil
}

type Finder interface {
	Find(ctx context.Context, id ...ID) ([]File, error)
}

// Getter provides methods to find Files.
type Getter interface {
	Finder
	FindByPath(ctx context.Context, path string) (File, error)
	FindAllByPath(ctx context.Context, path string) ([]File, error)
	FindByFingerprint(ctx context.Context, fp Fingerprint) ([]File, error)
	FindByZipFileID(ctx context.Context, zipFileID ID) ([]File, error)
	FindAllInPaths(ctx context.Context, p []string, limit, offset int) ([]File, error)
}

type Counter interface {
	CountAllInPaths(ctx context.Context, p []string) (int, error)
}

// Creator provides methods to create Files.
type Creator interface {
	Create(ctx context.Context, f File) error
}

// Updater provides methods to update Files.
type Updater interface {
	Update(ctx context.Context, f File) error
}

type Destroyer interface {
	Destroy(ctx context.Context, id ID) error
}

type GetterUpdater interface {
	Getter
	Updater
}

type GetterDestroyer interface {
	Getter
	Destroyer
}

// Store provides methods to find, create and update Files.
type Store interface {
	Getter
	Counter
	Creator
	Updater
	Destroyer

	IsPrimary(ctx context.Context, fileID ID) (bool, error)
}

// Decorator wraps the Decorate method to add additional functionality while scanning files.
type Decorator interface {
	Decorate(ctx context.Context, fs FS, f File) (File, error)
	IsMissingMetadata(ctx context.Context, fs FS, f File) bool
}

type FilteredDecorator struct {
	Decorator
	Filter
}

// Decorate runs the decorator if the filter accepts the file.
func (d *FilteredDecorator) Decorate(ctx context.Context, fs FS, f File) (File, error) {
	if d.Accept(ctx, f) {
		return d.Decorator.Decorate(ctx, fs, f)
	}
	return f, nil
}

func (d *FilteredDecorator) IsMissingMetadata(ctx context.Context, fs FS, f File) bool {
	if d.Accept(ctx, f) {
		return d.Decorator.IsMissingMetadata(ctx, fs, f)
	}

	return false
}
