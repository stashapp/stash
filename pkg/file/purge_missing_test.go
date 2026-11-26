package file

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/models/mocks"
	"github.com/stretchr/testify/mock"
)

type failingPurgeHandler struct{}

func (failingPurgeHandler) HandleFile(ctx context.Context, fileDeleter *Deleter, fileID models.FileID) error {
	return errors.New("delete failed")
}

func (failingPurgeHandler) HandleFolder(ctx context.Context, fileDeleter *Deleter, folderID models.FolderID) error {
	return errors.New("delete failed")
}

// runInJob runs fn with a real job progress, failing the test if it does not finish.
func runInJob(t *testing.T, fn func(ctx context.Context, progress *job.Progress)) {
	t.Helper()

	done := make(chan struct{})
	m := job.NewManager()
	defer m.Stop()
	m.Add(context.Background(), "test", job.MakeJobExec(func(ctx context.Context, progress *job.Progress) error {
		defer close(done)
		fn(ctx, progress)
		return nil
	}))

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("job did not finish")
	}
}

func makeMissingFiles(n int) []models.File {
	files := make([]models.File, n)
	for i := range files {
		files[i] = &models.BaseFile{ID: models.FileID(i + 1), Path: "missing.mp4"}
	}
	return files
}

func expectMissingCounts(db *mocks.Database, files, folders int) {
	db.File.On("CountMissingInPaths", mock.Anything, mock.Anything, mock.Anything).Return(files, nil)
	db.Folder.On("CountMissingInPaths", mock.Anything, mock.Anything, mock.Anything).Return(folders, nil)
}

// A full batch of files that fail to delete must not be refetched forever.
func TestPurgeMissingSkipsFailedDeletes(t *testing.T) {
	files := makeMissingFiles(purgeBatchSize)

	db := mocks.NewDatabase()
	expectMissingCounts(db, purgeBatchSize, 0)
	db.File.On("FindMissingInPaths", mock.Anything, mock.Anything, mock.Anything, purgeBatchSize, 0).Return(files, nil).Once()
	db.File.On("FindMissingInPaths", mock.Anything, mock.Anything, mock.Anything, purgeBatchSize, purgeBatchSize).Return([]models.File{}, nil).Once()
	db.Folder.On("FindMissingInPaths", mock.Anything, mock.Anything, mock.Anything, purgeBatchSize, 0).Return([]*models.Folder{}, nil).Once()

	p := &MissingPurger{
		Repository:    NewRepository(db.Repository()),
		PurgeHandlers: []PurgeHandler{failingPurgeHandler{}},
	}
	runInJob(t, func(ctx context.Context, progress *job.Progress) {
		p.PurgeMissing(ctx, PurgeMissingOptions{}, progress)
	})

	db.File.AssertExpectations(t)
	db.Folder.AssertExpectations(t)
}

func TestPurgeMissingDeletesMarkedFilesAndFolders(t *testing.T) {
	file := &models.BaseFile{ID: 1, Path: "/stash/missing.mp4"}
	folder := &models.Folder{ID: 2, Path: "/stash/missing"}
	cutoff := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	paths := []string{"/stash"}

	db := mocks.NewDatabase()
	expectMissingCounts(db, 1, 1)
	db.File.On("FindMissingInPaths", mock.Anything, paths, &cutoff, purgeBatchSize, 0).Return([]models.File{file}, nil).Once()
	db.Folder.On("FindMissingInPaths", mock.Anything, paths, &cutoff, purgeBatchSize, 0).Return([]*models.Folder{folder}, nil).Once()
	db.File.On("Destroy", mock.Anything, file.ID).Return(nil).Once()
	db.Folder.On("Destroy", mock.Anything, folder.ID).Return(nil).Once()

	p := &MissingPurger{Repository: NewRepository(db.Repository())}
	runInJob(t, func(ctx context.Context, progress *job.Progress) {
		p.PurgeMissing(ctx, PurgeMissingOptions{Paths: paths, MissingSinceBefore: &cutoff}, progress)
	})

	db.File.AssertExpectations(t)
	db.Folder.AssertExpectations(t)
}

// Dry run deletes nothing and pages through full batches.
func TestPurgeMissingDryRunDeletesNothing(t *testing.T) {
	db := mocks.NewDatabase()
	expectMissingCounts(db, purgeBatchSize+1, 0)
	db.File.On("FindMissingInPaths", mock.Anything, mock.Anything, mock.Anything, purgeBatchSize, 0).Return(makeMissingFiles(purgeBatchSize), nil).Once()
	db.File.On("FindMissingInPaths", mock.Anything, mock.Anything, mock.Anything, purgeBatchSize, purgeBatchSize).Return(makeMissingFiles(1), nil).Once()
	db.Folder.On("FindMissingInPaths", mock.Anything, mock.Anything, mock.Anything, purgeBatchSize, 0).Return([]*models.Folder{}, nil).Once()

	p := &MissingPurger{Repository: NewRepository(db.Repository())}
	runInJob(t, func(ctx context.Context, progress *job.Progress) {
		p.PurgeMissing(ctx, PurgeMissingOptions{DryRun: true}, progress)
	})

	db.File.AssertExpectations(t)
	db.Folder.AssertExpectations(t)
	db.File.AssertNotCalled(t, "Destroy", mock.Anything, mock.Anything)
}
