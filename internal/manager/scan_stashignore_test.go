//go:build integration
// +build integration

package manager

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stashapp/stash/internal/manager/config"
	"github.com/stashapp/stash/pkg/file"
	"github.com/stashapp/stash/pkg/models"
	"github.com/stashapp/stash/pkg/sqlite"
	"github.com/stashapp/stash/pkg/txn"

	// Necessary to register custom migrations.
	_ "github.com/stashapp/stash/pkg/sqlite/migrations"
)

// mockFingerprintCalculator returns empty fingerprints.
type mockFingerprintCalculator struct{}

func (m *mockFingerprintCalculator) CalculateFingerprints(f *models.BaseFile, o file.Opener, useExisting bool) ([]models.Fingerprint, error) {
	// Return a simple fingerprint based on path for testing.
	return []models.Fingerprint{
		{
			Type:        models.FingerprintTypeMD5,
			Fingerprint: fmt.Sprintf("md5-%s", f.Basename),
		},
	}, nil
}

// mockProgressReporter does nothing.
type mockProgressReporter struct{}

func (m *mockProgressReporter) AddTotal(total int)                        {}
func (m *mockProgressReporter) Increment()                                {}
func (m *mockProgressReporter) Definite()                                 {}
func (m *mockProgressReporter) ExecuteTask(description string, fn func()) { fn() }

// stashIgnorePathFilter wraps StashIgnoreFilter to implement PathFilter for testing.
// It provides a fixed library root for the filter.
type stashIgnorePathFilter struct {
	filter      *file.StashIgnoreFilter
	libraryRoot string
}

func (f *stashIgnorePathFilter) Accept(ctx context.Context, path string, info fs.FileInfo) bool {
	return f.filter.Accept(ctx, path, info, f.libraryRoot)
}

// createTestFileOnDisk creates a file with some content.
func createTestFileOnDisk(t *testing.T, dir, name string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("failed to create directory for %s: %v", path, err)
	}
	// Write some content so the file has a non-zero size.
	if err := os.WriteFile(path, []byte("test content for "+name), 0644); err != nil {
		t.Fatalf("failed to create file %s: %v", path, err)
	}
	return path
}

// createStashIgnoreFile creates a .stashignore file with the given content.
func createStashIgnoreFile(t *testing.T, dir, content string) {
	t.Helper()
	path := filepath.Join(dir, ".stashignore")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create .stashignore: %v", err)
	}
}

// setupTestDatabase creates a temporary SQLite database for testing.
func setupScanTestDatabase(t *testing.T) (*sqlite.Database, func()) {
	t.Helper()

	// Initialise empty config - needed by some migrations.
	_ = config.InitializeEmpty()

	// Create temporary database file.
	f, err := os.CreateTemp("", "stash-scan-test-*.sqlite")
	if err != nil {
		t.Fatalf("failed to create temp database file: %v", err)
	}
	f.Close()
	dbFile := f.Name()

	db := sqlite.NewDatabase()
	db.SetBlobStoreOptions(sqlite.BlobStoreOptions{
		UseDatabase: true,
	})

	if err := db.Open(dbFile); err != nil {
		os.Remove(dbFile)
		t.Fatalf("failed to open database: %v", err)
	}

	cleanup := func() {
		db.Close()
		os.Remove(dbFile)
	}

	return db, cleanup
}

func TestScannerWithStashIgnore(t *testing.T) {
	// Setup test database.
	db, cleanup := setupScanTestDatabase(t)
	defer cleanup()

	// Create temp directory structure.
	tmpDir := t.TempDir()

	// Create test files.
	createTestFileOnDisk(t, tmpDir, "video1.mp4")
	createTestFileOnDisk(t, tmpDir, "video2.mp4")
	createTestFileOnDisk(t, tmpDir, "ignore_me.mp4")
	createTestFileOnDisk(t, tmpDir, "subdir/video3.mp4")
	createTestFileOnDisk(t, tmpDir, "subdir/skip_this.mp4")
	createTestFileOnDisk(t, tmpDir, "excluded_dir/video4.mp4")
	createTestFileOnDisk(t, tmpDir, "temp/processing.mp4")

	// Create .stashignore file.
	stashignore := `# Ignore specific files
ignore_me.mp4
subdir/skip_this.mp4

# Ignore directories
excluded_dir/
temp/
`
	createStashIgnoreFile(t, tmpDir, stashignore)

	// Create scanner.
	repo := file.NewRepository(db.Repository())
	scanner := &file.Scanner{
		FS:                    &file.OsFS{},
		Repository:            repo,
		FingerprintCalculator: &mockFingerprintCalculator{},
	}

	// Create stashignore filter with library root.
	stashIgnoreFilter := &stashIgnorePathFilter{
		filter:      file.NewStashIgnoreFilter(),
		libraryRoot: tmpDir,
	}

	// Run scan.
	ctx := context.Background()
	scanner.Scan(ctx, nil, file.ScanOptions{
		Paths:         []string{tmpDir},
		ScanFilters:   []file.PathFilter{stashIgnoreFilter},
		ParallelTasks: 1,
	}, &mockProgressReporter{})

	// Verify results by checking what's in the database.
	var scannedPaths []string
	err := txn.WithTxn(ctx, db, func(ctx context.Context) error {
		// Check folders by path.
		checkDirs := []string{
			filepath.Join(tmpDir, "subdir"),
			filepath.Join(tmpDir, "excluded_dir"),
			filepath.Join(tmpDir, "temp"),
		}

		for _, dir := range checkDirs {
			f, err := db.Folder.FindByPath(ctx, dir, true)
			if err != nil {
				return fmt.Errorf("checking folder %s: %w", dir, err)
			}
			if f != nil {
				relPath, _ := filepath.Rel(tmpDir, dir)
				scannedPaths = append(scannedPaths, "dir:"+relPath)
			}
		}

		// Check specific files.
		checkFiles := []string{
			filepath.Join(tmpDir, "video1.mp4"),
			filepath.Join(tmpDir, "video2.mp4"),
			filepath.Join(tmpDir, "ignore_me.mp4"),
			filepath.Join(tmpDir, "subdir/video3.mp4"),
			filepath.Join(tmpDir, "subdir/skip_this.mp4"),
			filepath.Join(tmpDir, "excluded_dir/video4.mp4"),
			filepath.Join(tmpDir, "temp/processing.mp4"),
		}

		for _, path := range checkFiles {
			f, err := db.File.FindByPath(ctx, path, true)
			if err != nil {
				return fmt.Errorf("checking file %s: %w", path, err)
			}
			if f != nil {
				relPath, _ := filepath.Rel(tmpDir, path)
				scannedPaths = append(scannedPaths, "file:"+relPath)
			}
		}

		return nil
	})

	if err != nil {
		t.Fatalf("failed to verify scan results: %v", err)
	}

	sort.Strings(scannedPaths)

	// Expected: video1.mp4, video2.mp4, subdir/video3.mp4, and their folders.
	// NOT expected: ignore_me.mp4, subdir/skip_this.mp4, excluded_dir/*, temp/*.
	expectedPaths := []string{
		"dir:subdir",
		"file:subdir/video3.mp4",
		"file:video1.mp4",
		"file:video2.mp4",
	}
	sort.Strings(expectedPaths)

	if len(scannedPaths) != len(expectedPaths) {
		t.Errorf("scanned path count mismatch:\nexpected %d: %v\nactual %d: %v",
			len(expectedPaths), expectedPaths, len(scannedPaths), scannedPaths)
		return
	}

	for i := range expectedPaths {
		if scannedPaths[i] != expectedPaths[i] {
			t.Errorf("path mismatch at index %d:\nexpected: %s\nactual: %s",
				i, expectedPaths[i], scannedPaths[i])
		}
	}
}

func TestScannerWithNestedStashIgnore(t *testing.T) {
	// Setup test database.
	db, cleanup := setupScanTestDatabase(t)
	defer cleanup()

	// Create temp directory structure.
	tmpDir := t.TempDir()

	// Create test files.
	createTestFileOnDisk(t, tmpDir, "root.mp4")
	createTestFileOnDisk(t, tmpDir, "root.tmp")
	createTestFileOnDisk(t, tmpDir, "subdir/sub.mp4")
	createTestFileOnDisk(t, tmpDir, "subdir/sub.log")
	createTestFileOnDisk(t, tmpDir, "subdir/sub.tmp")

	// Root .stashignore excludes *.tmp.
	createStashIgnoreFile(t, tmpDir, "*.tmp\n")

	// Subdir .stashignore excludes *.log.
	createStashIgnoreFile(t, filepath.Join(tmpDir, "subdir"), "*.log\n")

	// Create scanner.
	repo := file.NewRepository(db.Repository())
	scanner := &file.Scanner{
		FS:                    &file.OsFS{},
		Repository:            repo,
		FingerprintCalculator: &mockFingerprintCalculator{},
	}

	// Create stashignore filter with library root.
	stashIgnoreFilter := &stashIgnorePathFilter{
		filter:      file.NewStashIgnoreFilter(),
		libraryRoot: tmpDir,
	}

	// Run scan.
	ctx := context.Background()
	scanner.Scan(ctx, nil, file.ScanOptions{
		Paths:         []string{tmpDir},
		ScanFilters:   []file.PathFilter{stashIgnoreFilter},
		ParallelTasks: 1,
	}, &mockProgressReporter{})

	// Verify results.
	var scannedFiles []string
	err := txn.WithTxn(ctx, db, func(ctx context.Context) error {
		checkFiles := []string{
			filepath.Join(tmpDir, "root.mp4"),
			filepath.Join(tmpDir, "root.tmp"),
			filepath.Join(tmpDir, "subdir/sub.mp4"),
			filepath.Join(tmpDir, "subdir/sub.log"),
			filepath.Join(tmpDir, "subdir/sub.tmp"),
		}

		for _, path := range checkFiles {
			f, err := db.File.FindByPath(ctx, path, true)
			if err != nil {
				return fmt.Errorf("checking file %s: %w", path, err)
			}
			if f != nil {
				relPath, _ := filepath.Rel(tmpDir, path)
				scannedFiles = append(scannedFiles, relPath)
			}
		}

		return nil
	})

	if err != nil {
		t.Fatalf("failed to verify scan results: %v", err)
	}

	sort.Strings(scannedFiles)

	// Expected: root.mp4, subdir/sub.mp4.
	// NOT expected: root.tmp (root ignore), subdir/sub.log (subdir ignore), subdir/sub.tmp (root ignore).
	expectedFiles := []string{
		"root.mp4",
		"subdir/sub.mp4",
	}
	sort.Strings(expectedFiles)

	if len(scannedFiles) != len(expectedFiles) {
		t.Errorf("scanned file count mismatch:\nexpected %d: %v\nactual %d: %v",
			len(expectedFiles), expectedFiles, len(scannedFiles), scannedFiles)
		return
	}

	for i := range expectedFiles {
		if scannedFiles[i] != expectedFiles[i] {
			t.Errorf("file mismatch at index %d:\nexpected: %s\nactual: %s",
				i, expectedFiles[i], scannedFiles[i])
		}
	}
}

func TestScannerWithoutStashIgnore(t *testing.T) {
	// Setup test database.
	db, cleanup := setupScanTestDatabase(t)
	defer cleanup()

	// Create temp directory structure (no .stashignore).
	tmpDir := t.TempDir()

	// Create test files.
	createTestFileOnDisk(t, tmpDir, "video1.mp4")
	createTestFileOnDisk(t, tmpDir, "video2.mp4")
	createTestFileOnDisk(t, tmpDir, "subdir/video3.mp4")

	// Create scanner.
	repo := file.NewRepository(db.Repository())
	scanner := &file.Scanner{
		FS:                    &file.OsFS{},
		Repository:            repo,
		FingerprintCalculator: &mockFingerprintCalculator{},
	}

	// Create stashignore filter with library root (but no .stashignore file exists).
	stashIgnoreFilter := &stashIgnorePathFilter{
		filter:      file.NewStashIgnoreFilter(),
		libraryRoot: tmpDir,
	}

	// Run scan.
	ctx := context.Background()
	scanner.Scan(ctx, nil, file.ScanOptions{
		Paths:         []string{tmpDir},
		ScanFilters:   []file.PathFilter{stashIgnoreFilter},
		ParallelTasks: 1,
	}, &mockProgressReporter{})

	// Verify all files were scanned.
	var scannedFiles []string
	err := txn.WithTxn(ctx, db, func(ctx context.Context) error {
		checkFiles := []string{
			filepath.Join(tmpDir, "video1.mp4"),
			filepath.Join(tmpDir, "video2.mp4"),
			filepath.Join(tmpDir, "subdir/video3.mp4"),
		}

		for _, path := range checkFiles {
			f, err := db.File.FindByPath(ctx, path, true)
			if err != nil {
				return fmt.Errorf("checking file %s: %w", path, err)
			}
			if f != nil {
				relPath, _ := filepath.Rel(tmpDir, path)
				scannedFiles = append(scannedFiles, relPath)
			}
		}

		return nil
	})

	if err != nil {
		t.Fatalf("failed to verify scan results: %v", err)
	}

	sort.Strings(scannedFiles)

	// All files should be scanned.
	expectedFiles := []string{
		"subdir/video3.mp4",
		"video1.mp4",
		"video2.mp4",
	}
	sort.Strings(expectedFiles)

	if len(scannedFiles) != len(expectedFiles) {
		t.Errorf("scanned file count mismatch:\nexpected %d: %v\nactual %d: %v",
			len(expectedFiles), expectedFiles, len(scannedFiles), scannedFiles)
		return
	}

	for i := range expectedFiles {
		if scannedFiles[i] != expectedFiles[i] {
			t.Errorf("file mismatch at index %d:\nexpected: %s\nactual: %s",
				i, expectedFiles[i], scannedFiles[i])
		}
	}
}

func TestScannerWithNegationPattern(t *testing.T) {
	// Setup test database.
	db, cleanup := setupScanTestDatabase(t)
	defer cleanup()

	// Create temp directory structure.
	tmpDir := t.TempDir()

	// Create test files.
	createTestFileOnDisk(t, tmpDir, "file1.tmp")
	createTestFileOnDisk(t, tmpDir, "file2.tmp")
	createTestFileOnDisk(t, tmpDir, "keep_this.tmp")
	createTestFileOnDisk(t, tmpDir, "video.mp4")

	// Create .stashignore with negation.
	stashignore := `*.tmp
!keep_this.tmp
`
	createStashIgnoreFile(t, tmpDir, stashignore)

	// Create scanner.
	repo := file.NewRepository(db.Repository())
	scanner := &file.Scanner{
		FS:                    &file.OsFS{},
		Repository:            repo,
		FingerprintCalculator: &mockFingerprintCalculator{},
	}

	// Create stashignore filter with library root.
	stashIgnoreFilter := &stashIgnorePathFilter{
		filter:      file.NewStashIgnoreFilter(),
		libraryRoot: tmpDir,
	}

	// Run scan.
	ctx := context.Background()
	scanner.Scan(ctx, nil, file.ScanOptions{
		Paths:         []string{tmpDir},
		ScanFilters:   []file.PathFilter{stashIgnoreFilter},
		ParallelTasks: 1,
	}, &mockProgressReporter{})

	// Verify results.
	var scannedFiles []string
	err := txn.WithTxn(ctx, db, func(ctx context.Context) error {
		checkFiles := []string{
			filepath.Join(tmpDir, "file1.tmp"),
			filepath.Join(tmpDir, "file2.tmp"),
			filepath.Join(tmpDir, "keep_this.tmp"),
			filepath.Join(tmpDir, "video.mp4"),
		}

		for _, path := range checkFiles {
			f, err := db.File.FindByPath(ctx, path, true)
			if err != nil {
				return fmt.Errorf("checking file %s: %w", path, err)
			}
			if f != nil {
				relPath, _ := filepath.Rel(tmpDir, path)
				scannedFiles = append(scannedFiles, relPath)
			}
		}

		return nil
	})

	if err != nil {
		t.Fatalf("failed to verify scan results: %v", err)
	}

	sort.Strings(scannedFiles)

	// Expected: keep_this.tmp (negated), video.mp4.
	// NOT expected: file1.tmp, file2.tmp.
	expectedFiles := []string{
		"keep_this.tmp",
		"video.mp4",
	}
	sort.Strings(expectedFiles)

	if len(scannedFiles) != len(expectedFiles) {
		t.Errorf("scanned file count mismatch:\nexpected %d: %v\nactual %d: %v",
			len(expectedFiles), expectedFiles, len(scannedFiles), scannedFiles)
		return
	}

	for i := range expectedFiles {
		if scannedFiles[i] != expectedFiles[i] {
			t.Errorf("file mismatch at index %d:\nexpected: %s\nactual: %s",
				i, expectedFiles[i], scannedFiles[i])
		}
	}
}
