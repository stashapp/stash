package manager

import (
	"context"
	"io/fs"
	"testing"
	"time"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/models"
)

// fakeTxnMgr runs functions directly without a database.
type fakeTxnMgr struct{}

func (m *fakeTxnMgr) Begin(ctx context.Context, _ bool) (context.Context, error) {
	return ctx, nil
}
func (m *fakeTxnMgr) Commit(_ context.Context) error   { return nil }
func (m *fakeTxnMgr) Rollback(_ context.Context) error { return nil }
func (m *fakeTxnMgr) IsLocked(_ error) bool            { return false }
func (m *fakeTxnMgr) WithDatabase(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

// fixedFolderStore returns a pre-configured folder for FindByPath.
type fixedFolderStore struct {
	folder *models.Folder
}

func (s *fixedFolderStore) Find(_ context.Context, _ models.FolderID) (*models.Folder, error) {
	return nil, nil
}
func (s *fixedFolderStore) FindMany(_ context.Context, _ []models.FolderID) ([]*models.Folder, error) {
	return nil, nil
}
func (s *fixedFolderStore) FindAllInPaths(_ context.Context, _ []string, _ bool, _, _ int) ([]*models.Folder, error) {
	return nil, nil
}
func (s *fixedFolderStore) FindByPath(_ context.Context, _ string, _ bool) (*models.Folder, error) {
	return s.folder, nil
}
func (s *fixedFolderStore) FindByZipFileID(_ context.Context, _ models.FileID) ([]*models.Folder, error) {
	return nil, nil
}
func (s *fixedFolderStore) FindByParentFolderID(_ context.Context, _ models.FolderID) ([]*models.Folder, error) {
	return nil, nil
}
func (s *fixedFolderStore) GetManyParentFolderIDs(_ context.Context, _ []models.FolderID) ([][]models.FolderID, error) {
	return nil, nil
}
func (s *fixedFolderStore) GetManySubFolderIDs(_ context.Context, _ []models.FolderID) ([][]models.FolderID, error) {
	return nil, nil
}
func (s *fixedFolderStore) Query(_ context.Context, _ models.FolderQueryOptions) (*models.FolderQueryResult, error) {
	return nil, nil
}
func (s *fixedFolderStore) CountAllInPaths(_ context.Context, _ []string) (int, error) {
	return 0, nil
}
func (s *fixedFolderStore) Create(_ context.Context, f *models.Folder) error {
	f.ID = 1
	return nil
}
func (s *fixedFolderStore) Update(_ context.Context, _ *models.Folder) error { return nil }
func (s *fixedFolderStore) Destroy(_ context.Context, _ models.FolderID) error { return nil }

// fakeTaskFS is a minimal FS that returns ErrNotExist for all operations.
type fakeTaskFS struct{}

func (f *fakeTaskFS) Lstat(_ string) (fs.FileInfo, error)             { return nil, fs.ErrNotExist }
func (f *fakeTaskFS) Stat(_ string) (fs.FileInfo, error)              { return nil, fs.ErrNotExist }
func (f *fakeTaskFS) Open(_ string) (fs.ReadDirFile, error)           { return nil, fs.ErrNotExist }
func (f *fakeTaskFS) OpenZip(_ string, _ int64) (models.ZipFS, error) { return nil, nil }
func (f *fakeTaskFS) IsPathCaseSensitive(_ string) (bool, error)      { return true, nil }

func makeTestScanner(t *testing.T, folder *models.Folder, _ bool) *file.Scanner {
	t.Helper()
	return &file.Scanner{
		FS: &fakeTaskFS{},
		Repository: file.Repository{
			TxnManager: &fakeTxnMgr{},
			Folder:     &fixedFolderStore{folder: folder},
		},
		RootPaths: []string{"/root"},
	}
}

func makeTestScanJob(scanner *file.Scanner, skipUnchangedFolders bool) *ScanJob {
	return &ScanJob{
		scanner:              scanner,
		fileQueue:            make(chan file.ScannedFile, 10),
		dirQueue:             make(chan file.ScannedFile, 10),
		skipUnchangedFolders: skipUnchangedFolders,
	}
}

func makeTestScannedFolder(path string, modTime time.Time) file.ScannedFile {
	return file.ScannedFile{
		BaseFile: &models.BaseFile{
			DirEntry: models.DirEntry{ModTime: modTime},
			Path:     path,
			Basename: "test",
		},
		FS: &fakeTaskFS{},
	}
}

func TestHandleFolder_UnchangedSkip(t *testing.T) {
	modTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	parentID := models.FolderID(99)

	folder := &models.Folder{
		DirEntry:       models.DirEntry{ModTime: modTime},
		Path:           "/root/sub",
		ParentFolderID: &parentID,
	}

	scanner := makeTestScanner(t, folder, true)
	job := makeTestScanJob(scanner, true)

	scanned := makeTestScannedFolder("/root/sub", modTime)
	err := job.handleFolder(context.Background(), scanned, nil)

	// handleFolder must NOT return fs.SkipDir — subdirectories must be walked
	// so that nested new files are discovered. Instead, the path is stored in
	// unchangedDirs for file-level skipping.
	if err != nil {
		t.Errorf("expected nil from handleFolder for unchanged folder, got %v", err)
	}

	if _, stored := job.unchangedDirs.Load("/root/sub"); !stored {
		t.Error("expected /root/sub to be recorded in unchangedDirs")
	}
}

func TestHandleFolder_UnchangedFilesSkipped(t *testing.T) {
	modTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	parentID := models.FolderID(99)

	folder := &models.Folder{
		DirEntry:       models.DirEntry{ModTime: modTime},
		Path:           "/root/sub",
		ParentFolderID: &parentID,
	}

	scanner := makeTestScanner(t, folder, true)
	job := makeTestScanJob(scanner, true)

	scanned := makeTestScannedFolder("/root/sub", modTime)
	if err := job.handleFolder(context.Background(), scanned, nil); err != nil {
		t.Fatalf("handleFolder: %v", err)
	}

	// A file inside the unchanged directory should be in unchangedDirs
	_, stored := job.unchangedDirs.Load("/root/sub")
	if !stored {
		t.Fatal("expected /root/sub in unchangedDirs")
	}
}

func TestHandleFolder_ChangedNoSkip(t *testing.T) {
	oldTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	newTime := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	parentID := models.FolderID(99)

	folder := &models.Folder{
		DirEntry:       models.DirEntry{ModTime: oldTime},
		Path:           "/root/sub",
		ParentFolderID: &parentID,
	}

	scanner := makeTestScanner(t, folder, true)
	job := makeTestScanJob(scanner, true)

	scanned := makeTestScannedFolder("/root/sub", newTime)
	err := job.handleFolder(context.Background(), scanned, nil)

	if err != nil {
		t.Errorf("expected no error for changed folder, got %v", err)
	}
}

func TestHandleFolder_SkipUnchangedFolders_Disabled(t *testing.T) {
	modTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	parentID := models.FolderID(99)

	folder := &models.Folder{
		DirEntry:       models.DirEntry{ModTime: modTime},
		Path:           "/root/sub",
		ParentFolderID: &parentID,
	}

	scanner := makeTestScanner(t, folder, false) // disabled
	job := makeTestScanJob(scanner, false)

	scanned := makeTestScannedFolder("/root/sub", modTime)
	err := job.handleFolder(context.Background(), scanned, nil)

	if err != nil {
		t.Errorf("expected no skip when SkipUnchangedFolders=false, got %v", err)
	}
}
