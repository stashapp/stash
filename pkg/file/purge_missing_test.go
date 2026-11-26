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

// A full batch of files that fail to delete must not be refetched forever.
func TestPurgeMissingSkipsFailedDeletes(t *testing.T) {
	const batchSize = 1000

	files := make([]models.File, batchSize)
	for i := range files {
		files[i] = &models.BaseFile{ID: models.FileID(i + 1), Path: "missing.mp4"}
	}

	db := mocks.NewDatabase()
	db.File.On("CountMissingInPaths", mock.Anything, mock.Anything, mock.Anything).Return(batchSize, nil)
	db.Folder.On("CountMissingInPaths", mock.Anything, mock.Anything, mock.Anything).Return(0, nil)
	db.File.On("FindMissingInPaths", mock.Anything, mock.Anything, mock.Anything, batchSize, 0).Return(files, nil).Once()
	db.File.On("FindMissingInPaths", mock.Anything, mock.Anything, mock.Anything, batchSize, batchSize).Return([]models.File{}, nil).Once()
	db.Folder.On("FindMissingInPaths", mock.Anything, mock.Anything, mock.Anything, batchSize, 0).Return([]*models.Folder{}, nil).Once()

	p := &MissingPurger{
		Repository:    NewRepository(db.Repository()),
		PurgeHandlers: []PurgeHandler{failingPurgeHandler{}},
	}
	done := make(chan struct{})
	m := job.NewManager()
	defer m.Stop()
	m.Add(context.Background(), "purge", job.MakeJobExec(func(ctx context.Context, progress *job.Progress) error {
		p.PurgeMissing(ctx, PurgeMissingOptions{}, progress)
		close(done)
		return nil
	}))

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("purge did not finish; failed deletes are being refetched")
	}

	db.File.AssertExpectations(t)
	db.Folder.AssertExpectations(t)
}
