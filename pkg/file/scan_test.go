package file

import (
	"context"
	"errors"
	"io/fs"
	"testing"
	"time"

	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type scanTestFS struct {
	caseSensitive bool
}

func (f scanTestFS) Stat(name string) (fs.FileInfo, error) {
	return nil, fs.ErrNotExist
}

func (f scanTestFS) Lstat(name string) (fs.FileInfo, error) {
	return nil, fs.ErrNotExist
}

func (f scanTestFS) Open(name string) (fs.ReadDirFile, error) {
	return nil, fs.ErrNotExist
}

func (f scanTestFS) OpenZip(name string, size int64) (models.ZipFS, error) {
	return nil, errors.New("zip files not supported")
}

func (f scanTestFS) IsPathCaseSensitive(path string) (bool, error) {
	return f.caseSensitive, nil
}

type scanTestFileInfo struct {
	name    string
	size    int64
	modTime time.Time
}

func (i scanTestFileInfo) Name() string       { return i.name }
func (i scanTestFileInfo) Size() int64        { return i.size }
func (i scanTestFileInfo) Mode() fs.FileMode  { return 0644 }
func (i scanTestFileInfo) ModTime() time.Time { return i.modTime }
func (i scanTestFileInfo) IsDir() bool        { return false }
func (i scanTestFileInfo) Sys() interface{}   { return nil }

type scanTestFingerprintCalculator struct {
	fingerprints []models.Fingerprint
}

func (c scanTestFingerprintCalculator) CalculateFingerprints(f *models.BaseFile, o Opener, useExisting bool) ([]models.Fingerprint, error) {
	return c.fingerprints, nil
}

func TestScannerScanFileMatchesCaseOnlyRenameOnCaseInsensitiveFilesystem(t *testing.T) {
	ctx := context.Background()
	db := mocks.NewDatabase()

	oldModTime := time.Date(2026, 7, 2, 12, 0, 0, 0, time.UTC)
	newModTime := oldModTime.Add(time.Hour)

	existing := &models.VideoFile{
		BaseFile: &models.BaseFile{
			ID:             1,
			Path:           `C:\stash\ALLCAPS.mp4`,
			Basename:       "ALLCAPS.mp4",
			ParentFolderID: 2,
			Size:           100,
			DirEntry: models.DirEntry{
				ModTime: oldModTime,
			},
			Fingerprints: []models.Fingerprint{
				{Type: models.FingerprintTypeOshash, Fingerprint: "old-oshash"},
				{Type: models.FingerprintTypeMD5, Fingerprint: "old-md5"},
				{Type: models.FingerprintTypePhash, Fingerprint: int64(1234)},
			},
		},
		Format:     "mp4",
		VideoCodec: "h264",
		AudioCodec: "aac",
	}

	newFingerprints := []models.Fingerprint{
		{Type: models.FingerprintTypeOshash, Fingerprint: "new-oshash"},
		{Type: models.FingerprintTypeMD5, Fingerprint: "new-md5"},
	}

	scanner := Scanner{
		FS: scanTestFS{caseSensitive: false},
		Repository: Repository{
			TxnManager: db,
			File:       db.File,
			Folder:     db.Folder,
		},
		FingerprintCalculator: scanTestFingerprintCalculator{fingerprints: newFingerprints},
	}

	scanned := ScannedFile{
		BaseFile: &models.BaseFile{
			Path:     `C:\stash\Allcaps.mp4`,
			Basename: "Allcaps.mp4",
			Size:     90,
			DirEntry: models.DirEntry{
				ModTime: newModTime,
			},
		},
		FS: scanTestFS{caseSensitive: false},
		Info: scanTestFileInfo{
			name:    "Allcaps.mp4",
			size:    90,
			modTime: newModTime,
		},
	}

	db.File.On("FindByPath", mock.Anything, scanned.Path, true).Return(nil, nil).Once()
	db.File.On("FindByPath", mock.Anything, scanned.Path, false).Return(existing, nil).Once()
	db.File.On("Update", mock.Anything, mock.MatchedBy(func(f models.File) bool {
		base := f.Base()
		return base.ID == existing.ID &&
			base.Basename == scanned.Basename &&
			base.Size == scanned.Size &&
			base.Fingerprints.Get(models.FingerprintTypeOshash) == "new-oshash" &&
			base.Fingerprints.Get(models.FingerprintTypeMD5) == "new-md5" &&
			base.Fingerprints.Get(models.FingerprintTypePhash) == nil
	})).Return(nil).Once()

	result, err := scanner.ScanFile(ctx, scanned)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.False(t, result.New)
	assert.True(t, result.Updated)
	assert.Equal(t, existing.ID, result.File.Base().ID)

	db.File.AssertExpectations(t)
}

// rejectAllFilter makes unchanged files skip their handlers.
type rejectAllFilter struct{}

func (rejectAllFilter) Accept(ctx context.Context, f models.File) bool { return false }

func TestScannerScanFileClearsMissingSinceOnUnchangedFile(t *testing.T) {
	modTime := time.Date(2026, 7, 2, 12, 0, 0, 0, time.UTC)
	missingSince := modTime.Add(time.Hour)
	fingerprints := []models.Fingerprint{{Type: models.FingerprintTypeOshash, Fingerprint: "oshash"}}

	for _, tc := range []struct {
		name         string
		missingSince *time.Time
	}{
		{"previously missing", &missingSince},
		{"not missing", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			db := mocks.NewDatabase()

			existing := &models.BaseFile{
				ID:           1,
				Path:         "/stash/archive.zip",
				Basename:     "archive.zip",
				Size:         100,
				DirEntry:     models.DirEntry{ModTime: modTime, MissingSince: tc.missingSince},
				Fingerprints: fingerprints,
			}

			scanner := Scanner{
				FS:                     scanTestFS{caseSensitive: true},
				Repository:             Repository{TxnManager: db, File: db.File, Folder: db.Folder},
				FingerprintCalculator:  scanTestFingerprintCalculator{fingerprints: fingerprints},
				HandlerRequiredFilters: []Filter{rejectAllFilter{}},
			}

			scanned := ScannedFile{
				BaseFile: &models.BaseFile{
					Path:     existing.Path,
					Basename: existing.Basename,
					Size:     existing.Size,
					DirEntry: models.DirEntry{ModTime: modTime},
				},
				FS:   scanTestFS{caseSensitive: true},
				Info: scanTestFileInfo{name: existing.Basename, size: existing.Size, modTime: modTime},
			}

			db.File.On("FindByPath", mock.Anything, scanned.Path, true).Return(existing, nil).Once()
			if tc.missingSince != nil {
				db.File.On("SetMissing", mock.Anything, existing.ID, (*time.Time)(nil)).Return(nil).Once()
			}

			result, err := scanner.ScanFile(ctx, scanned)
			require.NoError(t, err)
			require.NotNil(t, result)
			// signals the scan task to rescan zip contents, which were marked missing with the zip
			assert.Equal(t, tc.missingSince != nil, result.NoLongerMissing)
			assert.False(t, result.Updated)

			db.File.AssertExpectations(t)
			if tc.missingSince == nil {
				db.File.AssertNotCalled(t, "SetMissing", mock.Anything, mock.Anything, mock.Anything)
			}
		})
	}
}
