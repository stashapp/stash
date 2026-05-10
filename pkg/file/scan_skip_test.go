package file_test

import (
	"context"
	"io/fs"
	"testing"
	"time"

	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/models"
)

// fakeTxnManager runs the function directly without a real database.
type fakeTxnManager struct{}

func (m *fakeTxnManager) Begin(ctx context.Context, _ bool) (context.Context, error) {
	return ctx, nil
}
func (m *fakeTxnManager) Commit(_ context.Context) error    { return nil }
func (m *fakeTxnManager) Rollback(_ context.Context) error  { return nil }
func (m *fakeTxnManager) IsLocked(_ error) bool             { return false }
func (m *fakeTxnManager) WithDatabase(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

// fakeFolderStore records calls and returns configurable results.
type fakeFolderStore struct {
	existing *models.Folder
	created  []*models.Folder
	updated  []*models.Folder
}

func (s *fakeFolderStore) Find(_ context.Context, _ models.FolderID) (*models.Folder, error) {
	return nil, nil
}
func (s *fakeFolderStore) FindMany(_ context.Context, _ []models.FolderID) ([]*models.Folder, error) {
	return nil, nil
}
func (s *fakeFolderStore) FindAllInPaths(_ context.Context, _ []string, _ bool, _, _ int) ([]*models.Folder, error) {
	return nil, nil
}
func (s *fakeFolderStore) FindByPath(_ context.Context, _ string, _ bool) (*models.Folder, error) {
	return s.existing, nil
}
func (s *fakeFolderStore) FindByZipFileID(_ context.Context, _ models.FileID) ([]*models.Folder, error) {
	return nil, nil
}
func (s *fakeFolderStore) FindByParentFolderID(_ context.Context, _ models.FolderID) ([]*models.Folder, error) {
	return nil, nil
}
func (s *fakeFolderStore) GetManyParentFolderIDs(_ context.Context, _ []models.FolderID) ([][]models.FolderID, error) {
	return nil, nil
}
func (s *fakeFolderStore) GetManySubFolderIDs(_ context.Context, _ []models.FolderID) ([][]models.FolderID, error) {
	return nil, nil
}
func (s *fakeFolderStore) Query(_ context.Context, _ models.FolderQueryOptions) (*models.FolderQueryResult, error) {
	return nil, nil
}
func (s *fakeFolderStore) CountAllInPaths(_ context.Context, _ []string) (int, error) {
	return 0, nil
}
func (s *fakeFolderStore) Create(_ context.Context, f *models.Folder) error {
	s.created = append(s.created, f)
	f.ID = 1
	return nil
}
func (s *fakeFolderStore) Update(_ context.Context, f *models.Folder) error {
	s.updated = append(s.updated, f)
	return nil
}
func (s *fakeFolderStore) Destroy(_ context.Context, _ models.FolderID) error { return nil }

// fakeFS is a minimal FS that reports paths as case-sensitive.
// Lstat always returns os.ErrNotExist so detectFolderMove aborts gracefully.
type fakeFS struct{}

func (f *fakeFS) Lstat(_ string) (fs.FileInfo, error)              { return nil, fs.ErrNotExist }
func (f *fakeFS) Stat(_ string) (fs.FileInfo, error)               { return nil, fs.ErrNotExist }
func (f *fakeFS) Open(_ string) (fs.ReadDirFile, error)            { return nil, fs.ErrNotExist }
func (f *fakeFS) OpenZip(_ string, _ int64) (models.ZipFS, error)  { return nil, nil }
func (f *fakeFS) IsPathCaseSensitive(_ string) (bool, error)       { return true, nil }

func makeTestScanner(store *fakeFolderStore) *file.Scanner {
	return &file.Scanner{
		FS: &fakeFS{},
		Repository: file.Repository{
			TxnManager: &fakeTxnManager{},
			Folder:     store,
		},
		RootPaths: []string{"/root"},
	}
}

func makeScannedFolder(path string, modTime time.Time) file.ScannedFile {
	return file.ScannedFile{
		BaseFile: &models.BaseFile{
			DirEntry: models.DirEntry{
				ModTime: modTime,
			},
			Path:     path,
			Basename: "test",
		},
		FS: &fakeFS{},
	}
}

func TestCheckFolder_ExistingUnchanged_ReturnsUnchangedTrue(t *testing.T) {
	modTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	parentID := models.FolderID(99)

	store := &fakeFolderStore{
		existing: &models.Folder{
			ID:             models.FolderID(1),
			DirEntry:       models.DirEntry{ModTime: modTime},
			Path:           "/root/sub",
			ParentFolderID: &parentID,
		},
	}
	scanner := makeTestScanner(store)

	scanned := makeScannedFolder("/root/sub", modTime)

	_, unchanged, err := scanner.CheckFolder(context.Background(), scanned)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !unchanged {
		t.Error("expected unchanged=true for matching modtime, got unchanged=false")
	}
}

func TestCheckFolder_ExistingChanged_ReturnsUnchangedFalse(t *testing.T) {
	oldTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	newTime := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	parentID := models.FolderID(99)

	store := &fakeFolderStore{
		existing: &models.Folder{
			ID:             models.FolderID(1),
			DirEntry:       models.DirEntry{ModTime: oldTime},
			Path:           "/root/sub",
			ParentFolderID: &parentID,
		},
	}
	scanner := makeTestScanner(store)

	scanned := makeScannedFolder("/root/sub", newTime)

	_, unchanged, err := scanner.CheckFolder(context.Background(), scanned)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if unchanged {
		t.Error("expected unchanged=false for changed modtime, got unchanged=true")
	}
}

func TestCheckFolder_NotFound_ReturnsNilUnchangedFalse(t *testing.T) {
	store := &fakeFolderStore{existing: nil}
	scanner := makeTestScanner(store)

	scanned := makeScannedFolder("/root/new", time.Now())

	existing, unchanged, err := scanner.CheckFolder(context.Background(), scanned)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if existing != nil {
		t.Error("expected existing=nil for unknown folder")
	}
	if unchanged {
		t.Error("expected unchanged=false for unknown folder")
	}
}

func TestScanFolder_ExistingUnchanged_ReturnsChangedFalse(t *testing.T) {
	modTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	parentID := models.FolderID(99)

	store := &fakeFolderStore{
		existing: &models.Folder{
			DirEntry:       models.DirEntry{ModTime: modTime},
			Path:           "/root/sub",
			ParentFolderID: &parentID,
		},
	}
	scanner := makeTestScanner(store)

	scanned := makeScannedFolder("/root/sub", modTime)

	_, changed, err := scanner.ScanFolder(context.Background(), scanned)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if changed {
		t.Error("expected changed=false for unchanged folder, got changed=true")
	}
	if len(store.updated) != 0 {
		t.Errorf("expected no updates, got %d", len(store.updated))
	}
}

func TestScanFolder_ExistingChanged_ReturnsChangedTrue(t *testing.T) {
	oldTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	newTime := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	parentID := models.FolderID(99)

	store := &fakeFolderStore{
		existing: &models.Folder{
			DirEntry:       models.DirEntry{ModTime: oldTime},
			Path:           "/root/sub",
			ParentFolderID: &parentID,
		},
	}
	scanner := makeTestScanner(store)

	scanned := makeScannedFolder("/root/sub", newTime)

	_, changed, err := scanner.ScanFolder(context.Background(), scanned)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !changed {
		t.Error("expected changed=true for modified folder, got changed=false")
	}
}

func TestScanFolder_NewFolder_ReturnsChangedTrue(t *testing.T) {
	store := &fakeFolderStore{existing: nil}
	scanner := makeTestScanner(store)

	scanned := makeScannedFolder("/root/brandnew", time.Now())

	_, changed, err := scanner.ScanFolder(context.Background(), scanned)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !changed {
		t.Error("expected changed=true for new folder, got changed=false")
	}
}
