package file

import (
	"context"
	"strconv"
	"time"
)

type FileID int32

func (i FileID) String() string {
	return strconv.Itoa(int(i))
}

type DirEntry struct {
	Path           string    `json:"path"`
	ParentFolderID *FolderID `json:"parent_folder_id"`
	ZipFileID      *FileID   `json:"zip_file_id"`

	ModTime      time.Time  `json:"mod_time"`
	MissingSince *time.Time `json:"missing_since"`

	LastScanned time.Time `json:"last_scanned"`
}

func (e *DirEntry) scanned() {
	e.LastScanned = time.Now()
	e.MissingSince = nil
}

type File interface {
	Basic() *BasicFile
	SetFingerprint(fp Fingerprint)
}

type BasicFile struct {
	ID FileID `json:"id"`

	DirEntry
	Basename string `json:"basename"`

	Fingerprints []Fingerprint `json:"fingerprints"`

	Size int64 `json:"size"`

	Title   string     `json:"title"`
	Details string     `json:"details"`
	Date    *time.Time `json:"date"`
	Rating  *int       `json:"rating"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (f *BasicFile) SetFingerprint(fp Fingerprint) {
	for i, existing := range f.Fingerprints {
		if existing.Type == fp.Type {
			f.Fingerprints[i] = fp
			return
		}
	}

	f.Fingerprints = append(f.Fingerprints, fp)
}

func (f *BasicFile) Basic() *BasicFile {
	return f
}

type Getter interface {
	GetByPath(ctx context.Context, path string) (File, error)
	GetByFingerprint(ctx context.Context, fp Fingerprint) ([]File, error)
}

type Creator interface {
	Create(ctx context.Context, f *BasicFile) error
}

type Updater interface {
	Update(ctx context.Context, f *BasicFile) error
}

type Store interface {
	Getter
	Creator
	Updater
}

type MissedMarker interface {
	MarkMissing(ctx context.Context, scanTime time.Time) error
}
