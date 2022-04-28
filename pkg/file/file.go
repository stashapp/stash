package file

import (
	"context"
	"strconv"
	"time"
)

// ID represents an ID of a file.
type ID int32

func (i ID) String() string {
	return strconv.Itoa(int(i))
}

// DirEntry represents a file or directory in the file system.
type DirEntry struct {
	ZipFileID *ID `json:"zip_file_id"`

	ModTime      time.Time  `json:"mod_time"`
	MissingSince *time.Time `json:"missing_since"`

	LastScanned time.Time `json:"last_scanned"`
}

func (e *DirEntry) scanned() {
	e.LastScanned = time.Now()
	e.MissingSince = nil
}

// File represents a file in the file system.
type File interface {
	Base() *BaseFile
	SetFingerprints(fp []Fingerprint)
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
func (f *BaseFile) SetFingerprints(fp []Fingerprint) {
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

// Getter provides methods to find Files.
type Getter interface {
	FindByPath(ctx context.Context, path string) (File, error)
	FindByFingerprint(ctx context.Context, fp Fingerprint) ([]File, error)
}

// Creator provides methods to create Files.
type Creator interface {
	Create(ctx context.Context, f File) error
}

// Updater provides methods to update Files.
type Updater interface {
	Update(ctx context.Context, f File) error
}

// Store provides methods to find, create and update Files.
type Store interface {
	Getter
	Creator
	Updater
	MissingMarker
}

// MissingMarker wraps the MarkMissing method.
type MissingMarker interface {
	MarkMissing(ctx context.Context, scanStartTime time.Time, scanPaths []string) (int, error)
}

// Decorator wraps the Decorate method to add additional functionality while scanning files.
type Decorator interface {
	Decorate(ctx context.Context, fs FS, f File) (File, error)
}

type FilteredDecorator struct {
	Decorator
	Filter
}

// Decorate runs the decorator if the filter accepts the file.
func (d *FilteredDecorator) Decorate(ctx context.Context, fs FS, f File) (File, error) {
	if d.Accept(f) {
		return d.Decorator.Decorate(ctx, fs, f)
	}
	return f, nil
}
