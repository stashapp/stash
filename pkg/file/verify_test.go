package file

import (
	"context"
	"io/fs"
	"testing"
	"time"

	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/mock"
)

// verifyTestFS reports only the given paths as present.
type verifyTestFS struct {
	scanTestFS
	present map[string]bool
}

func (f verifyTestFS) Lstat(name string) (fs.FileInfo, error) {
	if f.present[name] {
		return scanTestFileInfo{name: name}, nil
	}
	return nil, fs.ErrNotExist
}

type acceptAllPathFilter struct{}

func (acceptAllPathFilter) Accept(ctx context.Context, path string, info fs.FileInfo, zipFilePath string) bool {
	return true
}

type verifyFixture struct {
	db *mocks.Database

	present        *models.BaseFile
	missing        *models.BaseFile
	alreadyMissing *models.BaseFile
	folder         *models.Folder
}

func newVerifyFixture() *verifyFixture {
	since := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	f := &verifyFixture{
		db:             mocks.NewDatabase(),
		present:        &models.BaseFile{ID: 1, Path: "/stash/present.mp4"},
		missing:        &models.BaseFile{ID: 2, Path: "/stash/missing.mp4"},
		alreadyMissing: &models.BaseFile{ID: 3, Path: "/stash/gone.mp4", DirEntry: models.DirEntry{MissingSince: &since}},
		folder:         &models.Folder{ID: 1, Path: "/stash"},
	}

	db := f.db
	db.File.On("CountAllInPaths", mock.Anything, mock.Anything).Return(3, nil)
	db.Folder.On("CountAllInPaths", mock.Anything, mock.Anything).Return(1, nil)
	db.Folder.On("FindAllInPaths", mock.Anything, mock.Anything, mock.Anything, purgeBatchSize, 0).
		Return([]*models.Folder{f.folder}, nil).Once()
	db.File.On("FindAllInPaths", mock.Anything, mock.Anything, mock.Anything, purgeBatchSize, 0).
		Return([]models.File{f.present, f.missing, f.alreadyMissing}, nil).Once()
	db.File.On("FindByZipFileID", mock.Anything, f.missing.ID).Return([]models.File{}, nil)
	db.Folder.On("FindByZipFileID", mock.Anything, f.missing.ID).Return([]*models.Folder{}, nil)

	return f
}

func (f *verifyFixture) run(t *testing.T, options VerifyOptions) {
	v := &Verifier{
		FS:         verifyTestFS{present: map[string]bool{f.present.Path: true, f.folder.Path: true}},
		Repository: NewRepository(f.db.Repository()),
	}
	options.PathFilter = acceptAllPathFilter{}

	runInJob(t, func(ctx context.Context, progress *job.Progress) {
		v.Verify(ctx, options, progress)
	})
}

func TestVerifyMarksOnlyNewlyMissingFiles(t *testing.T) {
	f := newVerifyFixture()
	f.db.File.On("SetMissing", mock.Anything, f.missing.ID, mock.MatchedBy(func(t *time.Time) bool { return t != nil })).
		Return(nil).Once()

	f.run(t, VerifyOptions{})

	f.db.File.AssertExpectations(t)
	f.db.Folder.AssertNotCalled(t, "SetMissing", mock.Anything, mock.Anything, mock.Anything)
	f.db.File.AssertNotCalled(t, "Destroy", mock.Anything, mock.Anything)
}

func TestVerifyDryRunChangesNothing(t *testing.T) {
	f := newVerifyFixture()

	// purge is ignored in dry run
	f.run(t, VerifyOptions{DryRun: true, PurgeMissing: true})

	f.db.File.AssertNotCalled(t, "SetMissing", mock.Anything, mock.Anything, mock.Anything)
	f.db.File.AssertNotCalled(t, "Destroy", mock.Anything, mock.Anything)
}

func TestVerifyPurgeDeletesNewlyMissingFiles(t *testing.T) {
	f := newVerifyFixture()
	f.db.File.On("Destroy", mock.Anything, f.missing.ID).Return(nil).Once()

	f.run(t, VerifyOptions{PurgeMissing: true})

	f.db.File.AssertExpectations(t)
	f.db.File.AssertNotCalled(t, "SetMissing", mock.Anything, mock.Anything, mock.Anything)
}
